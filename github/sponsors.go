package github

import (
	"context"
	"database/sql"
	"log"
	"time"
)

// Sponsors sets the golfer sponsor flag for sponsoring.
func Sponsors(db *sql.DB) {
	if accessToken == "" {
		return
	}

	var query struct {
		RateLimit struct {
			Cost, Limit, Remaining int
			ResetAt                time.Time
		}
		Viewer struct {
			SponsorshipsAsMaintainer struct {
				Nodes []struct {
					SponsorEntity struct {
						User struct{ DatabaseID int } `graphql:"... on User"`
					}
				}
			} `graphql:"sponsorshipsAsMaintainer(first: 100)"`
		}
	}

	if err := client.Query(context.Background(), &query, nil); err != nil {
		panic(err)
	}

	log.Printf(
		"GitHub API: Spent %d, %d/%d left, resets in %v",
		query.RateLimit.Cost,
		query.RateLimit.Remaining,
		query.RateLimit.Limit,
		time.Until(query.RateLimit.ResetAt).Round(time.Second),
	)

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	defer tx.Rollback()

	if _, err := tx.Exec("UPDATE users SET sponsor = false"); err != nil {
		panic(err)
	}

	for _, node := range query.Viewer.SponsorshipsAsMaintainer.Nodes {
		if _, err := tx.Exec(
			"UPDATE users SET sponsor = true WHERE id = $1",
			node.SponsorEntity.User.DatabaseID,
		); err != nil {
			panic(err)
		}
	}

	if err := tx.Commit(); err != nil {
		panic(err)
	}
}
