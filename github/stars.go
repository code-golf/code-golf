package github

import (
	"database/sql"
	"encoding/json"
	"time"
)

// Stars awards a trophy for starring the repo.
func Stars(db *sql.DB) {
	if accessToken == "" {
		return
	}

	stargazers := map[int]time.Time{}

	edges, err := graphQL("stargazers", `
		stargazers(after: $cursor first: 100) {
			edges { node { databaseId } starredAt }
			pageInfo { endCursor hasNextPage }
		}
	`)
	if err != nil {
		panic(err)
	}

	for _, e := range edges {
		var edge struct {
			Node      struct{ DatabaseID int }
			StarredAt time.Time
		}

		if err := json.Unmarshal(e, &edge); err != nil {
			panic(err)
		}

		stargazers[edge.Node.DatabaseID] = edge.StarredAt
	}

	awardTrophies(db, stargazers, "my-god-its-full-of-stars")
}
