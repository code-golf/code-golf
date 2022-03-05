package github

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

var accessToken = os.Getenv("GITHUB_ACCESS_TOKEN")

var client = githubv4.NewClient(
	oauth2.NewClient(
		context.Background(),
		oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken}),
	),
)

type pageInfo struct {
	EndCursor   githubv4.String
	HasNextPage bool
}

type rateLimit struct {
	Cost, Limit, Remaining int
	ResetAt                time.Time
}

func Run(db *pgxpool.Pool) {
	if accessToken == "" {
		return
	}

	var (
		cost  int
		limit rateLimit
	)

	for _, job := range []func(*pgxpool.Pool) []rateLimit{
		ideas, pullRequests, sponsors, stars,
	} {
		for _, limit = range job(db) {
			cost += limit.Cost
		}
	}

	log.Printf(
		"GitHub API: Spent %d, %d/%d left, resets in %v",
		cost, limit.Remaining, limit.Limit, time.Until(limit.ResetAt).Round(time.Second),
	)
}

func awardCheevos(db *pgxpool.Pool, earnedUsers map[int]time.Time, cheevoID string) {
	rows, err := db.Query(
		context.Background(),
		"SELECT earned, user_id FROM trophies WHERE trophy = $1",
		cheevoID,
	)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var earned time.Time
		var userID int

		if err := rows.Scan(&earned, &userID); err != nil {
			panic(err)
		}

		if newEarned, ok := earnedUsers[userID]; ok {
			delete(earnedUsers, userID)

			if earned != newEarned {
				if _, err := db.Exec(
					context.Background(),
					`UPDATE trophies
					    SET earned  = $1
					  WHERE trophy  = $2
					    AND user_id = $3`,
					newEarned,
					cheevoID,
					userID,
				); err != nil {
					panic(err)
				}
			}
		} else if _, err := db.Exec(
			context.Background(),
			"DELETE FROM trophies WHERE trophy = $1 AND user_id = $2",
			cheevoID,
			userID,
		); err != nil {
			panic(err)
		}
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	for userID, earned := range earnedUsers {
		if _, err := db.Exec(
			context.Background(),
			`INSERT INTO trophies SELECT $1, $2, $3
			WHERE EXISTS (SELECT * FROM users WHERE id = $2)`,
			earned,
			userID,
			cheevoID,
		); err != nil {
			panic(err)
		}
	}
}
