package github

import (
	"database/sql"
	"encoding/json"
)

// Ideas imports open issues with the idea label into the ideas table.
func Ideas(db *sql.DB) {
	if accessToken == "" {
		return
	}

	edges, err := graphQL("issues", `
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
	`)
	if err != nil {
		panic(err)
	}

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	defer tx.Rollback()

	if _, err := tx.Exec("TRUNCATE ideas"); err != nil {
		panic(err)
	}

	for _, e := range edges {
		type thumbs struct{ TotalCount int }

		var edge struct {
			Node struct {
				Number     int
				ThumbsDown thumbs
				ThumbsUp   thumbs
				Title      string
			}
		}

		if err := json.Unmarshal(e, &edge); err != nil {
			panic(err)
		}

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
