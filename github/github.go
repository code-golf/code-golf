package github

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

var accessToken = os.Getenv("GITHUB_ACCESS_TOKEN")

func awardTrophies(db *sql.DB, earnedUsers map[int]time.Time, trophy string) {
	rows, err := db.Query(
		"SELECT earned, user_id FROM trophies WHERE trophy = $1",
		trophy,
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
					`UPDATE trophies
					    SET earned  = $1
					  WHERE trophy  = $2
					    AND user_id = $3`,
					newEarned,
					trophy,
					userID,
				); err != nil {
					panic(err)
				}
			}
		} else if _, err := db.Exec(
			"DELETE FROM trophies WHERE trophy = $1 AND user_id = $2",
			trophy,
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
			`INSERT INTO trophies SELECT $1, $2, $3
			WHERE EXISTS (SELECT * FROM users WHERE id = $2)`,
			earned,
			userID,
			trophy,
		); err != nil {
			panic(err)
		}
	}
}

func graphQL(key, query string) ([]json.RawMessage, error) {
	payload := struct {
		Query     string `json:"query"`
		Variables struct {
			Cursor string `json:"cursor,omitempty"`
		} `json:"variables"`
	}{
		Query: `query($cursor: String) {
			rateLimit { cost limit remaining resetAt }
			repository(name: "code-golf" owner: "code-golf") {` + query + `}
		}`,
	}

	var edges []json.RawMessage

	for {
		var body bytes.Buffer

		if err := json.NewEncoder(&body).Encode(payload); err != nil {
			return nil, err
		}

		req, err := http.NewRequest("POST", "https://api.github.com/graphql", &body)
		if err != nil {
			return nil, err
		}

		req.Header.Add("Authorization", "bearer "+accessToken)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}

		defer res.Body.Close()

		var data struct {
			Data struct {
				RateLimit struct {
					Cost      int
					Limit     int
					Remaining int
					ResetAt   time.Time
				}
				Repository map[string]struct {
					Edges    []json.RawMessage
					PageInfo struct {
						EndCursor   string
						HasNextPage bool
					}
				}
			}
			Errors []interface{}
		}

		if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
			return nil, err
		}

		if len(data.Errors) != 0 {
			//nolint:goerr113
			return nil, fmt.Errorf("%v", data.Errors)
		}

		fmt.Printf(
			"GitHub API: Spent %d, %d/%d left, resets in %v\n",
			data.Data.RateLimit.Cost,
			data.Data.RateLimit.Remaining,
			data.Data.RateLimit.Limit,
			time.Until(data.Data.RateLimit.ResetAt).Round(time.Second),
		)

		edges = append(edges, data.Data.Repository[key].Edges...)

		if data.Data.Repository[key].PageInfo.HasNextPage {
			payload.Variables.Cursor = data.Data.Repository[key].PageInfo.EndCursor
		} else {
			break
		}
	}

	return edges, nil
}
