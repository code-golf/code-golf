package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

var accessToken = os.Getenv("GITHUB_ACCESS_TOKEN")

// FIXME Not a route, but only routes have access to the DB.
func Ideas() {
	if accessToken == "" {
		return
	}

	query := struct {
		Query     string `json:"query"`
		Variables struct {
			Cursor string `json:"cursor,omitempty"`
		} `json:"variables"`
	}{
		Query: `query($cursor: String) {
			rateLimit { cost limit remaining resetAt }
			repository(name: "code-golf" owner: "JRaspass") {
				issues(after: $cursor first: 100 labels: "idea" states: OPEN) {
					edges {
						node {
							number
							thumbsDown: reactions(content: THUMBS_DOWN) { totalCount }
							thumbsUp:   reactions(content: THUMBS_UP  ) { totalCount }
							title
						}
					}
					pageInfo { endCursor hasNextPage }
				}
			}
		}`,
	}

	type thumbs struct {
		TotalCount int `json:"totalCount"`
	}

	var data struct {
		Data struct {
			RateLimit struct {
				Cost      int       `json:"cost"`
				Limit     int       `json:"limit"`
				Remaining int       `json:"remaining"`
				ResetAt   time.Time `json:"resetAt"`
			}
			Repository struct {
				Issues struct {
					Edges []struct {
						Node struct {
							Number     int    `json:"number"`
							ThumbsDown thumbs `json:"thumbsDown"`
							ThumbsUp   thumbs `json:"thumbsUp"`
							Title      string `json:"title"`
						} `json:"node"`
					} `json:"edges"`
					PageInfo struct {
						EndCursor   string `json:"endCursor"`
						HasNextPage bool   `json:"hasNextPage"`
					} `json:"pageInfo"`
				} `json:"issues"`
			} `json:"repository"`
		} `json:"data"`
		Errors []interface{} `json:"errors"`
	}

	if err := graphQL(query, &data); err != nil {
		panic(err)
	}

	if len(data.Errors) != 0 {
		panic(fmt.Sprint(data.Errors))
	}

	fmt.Printf(
		"GitHub API: Spent %d, %d/%d left, resets in %v\n",
		data.Data.RateLimit.Cost,
		data.Data.RateLimit.Remaining,
		data.Data.RateLimit.Limit,
		time.Until(data.Data.RateLimit.ResetAt).Round(time.Second),
	)

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	defer tx.Rollback()

	if _, err := tx.Exec("TRUNCATE ideas"); err != nil {
		panic(err)
	}

	for _, edge := range data.Data.Repository.Issues.Edges {
		if _, err := tx.Exec(
			"INSERT INTO ideas VALUES ($1, $2, $3, $4)",
			edge.Node.Number,
			edge.Node.ThumbsDown.TotalCount,
			edge.Node.ThumbsUp.TotalCount,
			edge.Node.Title,
		); err != nil {
			panic(err)
		}
	}

	if err := tx.Commit(); err != nil {
		panic(err)
	}
}

// FIXME Not a route, but only routes have access to the DB.
func Stars() {
	if accessToken == "" {
		return
	}

	stargazers := map[int]time.Time{}

	query := struct {
		Query     string `json:"query"`
		Variables struct {
			Cursor string `json:"cursor,omitempty"`
		} `json:"variables"`
	}{
		Query: `query($cursor: String) {
			rateLimit { cost limit remaining resetAt }
			repository(name: "code-golf" owner: "JRaspass") {
				stargazers(after: $cursor first: 100) {
					edges { node { databaseId } starredAt }
					pageInfo { endCursor hasNextPage }
				}
			}
		}`,
	}

	for {
		var data struct {
			Data struct {
				RateLimit struct {
					Cost      int       `json:"cost"`
					Limit     int       `json:"limit"`
					Remaining int       `json:"remaining"`
					ResetAt   time.Time `json:"resetAt"`
				}
				Repository struct {
					Stargazers struct {
						Edges []struct {
							Node struct {
								DatabaseID int `json:"databaseId"`
							} `json:"node"`
							StarredAt time.Time `json:"starredAt"`
						} `json:"edges"`
						PageInfo struct {
							EndCursor   string `json:"endCursor"`
							HasNextPage bool   `json:"hasNextPage"`
						} `json:"pageInfo"`
					} `json:"stargazers"`
				} `json:"repository"`
			} `json:"data"`
			Errors []interface{} `json:"errors"`
		}

		if err := graphQL(query, &data); err != nil {
			panic(err)
		}

		if len(data.Errors) != 0 {
			panic(fmt.Sprint(data.Errors))
		}

		fmt.Printf(
			"GitHub API: Spent %d, %d/%d left, resets in %v\n",
			data.Data.RateLimit.Cost,
			data.Data.RateLimit.Remaining,
			data.Data.RateLimit.Limit,
			time.Until(data.Data.RateLimit.ResetAt).Round(time.Second),
		)

		for _, edge := range data.Data.Repository.Stargazers.Edges {
			stargazers[edge.Node.DatabaseID] = edge.StarredAt
		}

		if data.Data.Repository.Stargazers.PageInfo.HasNextPage {
			query.Variables.Cursor =
				data.Data.Repository.Stargazers.PageInfo.EndCursor
		} else {
			break
		}
	}

	rows, err := db.Query(`
		SELECT earned, user_id
		  FROM trophies
		 WHERE trophy = 'my-god-its-full-of-stars'
	`)
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

		if newEarned, ok := stargazers[userID]; ok {
			delete(stargazers, userID)

			if earned != newEarned {
				if _, err := db.Exec(
					`UPDATE trophies
					    SET earned = $1
					  WHERE trophy = 'my-god-its-full-of-stars'
					    AND user_id = $2`,
					newEarned,
					userID,
				); err != nil {
					panic(err)
				}
			}
		} else if _, err := db.Exec(
			`DELETE FROM trophies
			  WHERE trophy = 'my-god-its-full-of-stars'
			    AND user_id = $1`,
			userID,
		); err != nil {
			panic(err)
		}
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	for userID, earned := range stargazers {
		if _, err := db.Exec(
			`INSERT INTO trophies
				  SELECT $1, $2, 'my-god-its-full-of-stars'
			WHERE EXISTS (SELECT * FROM users WHERE id = $2)`,
			earned,
			userID,
		); err != nil {
			panic(err)
		}
	}
}

func graphQL(query, data interface{}) error {
	var body bytes.Buffer

	if err := json.NewEncoder(&body).Encode(query); err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://api.github.com/graphql", &body)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "bearer "+accessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	return json.NewDecoder(res.Body).Decode(&data)
}
