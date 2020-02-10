package github

import (
	"database/sql"
	"encoding/json"
	"time"
)

// PullRequests awards a trophy for merging a PR.
func PullRequests(db *sql.DB) {
	if accessToken == "" {
		return
	}

	pullRequests := map[int]time.Time{}

	edges, err := graphQL("pullRequests", `
		pullRequests(after: $cursor first: 100 states: MERGED) {
			edges { node { author { ...on User { databaseId } } mergedAt } }
			pageInfo { endCursor hasNextPage }
		}
	`)
	if err != nil {
		panic(err)
	}

	for _, e := range edges {
		var edge struct {
			Node struct {
				Author   struct{ DatabaseID int }
				MergedAt time.Time
			}
		}

		if err := json.Unmarshal(e, &edge); err != nil {
			panic(err)
		}

		mergedAt, ok := pullRequests[edge.Node.Author.DatabaseID]
		if !ok || edge.Node.MergedAt.Before(mergedAt) {
			pullRequests[edge.Node.Author.DatabaseID] = edge.Node.MergedAt
		}
	}

	awardTrophies(db, pullRequests, "patches-welcome")
}
