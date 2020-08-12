package github

import (
	"context"
	"database/sql"
	"time"

	"github.com/shurcooL/githubv4"
)

func stars(db *sql.DB) (limits []rateLimit) {
	stargazers := map[int]time.Time{}

	var query struct {
		RateLimit  rateLimit
		Repository struct {
			Stargazers struct {
				Edges []struct {
					Node      struct{ DatabaseID int }
					StarredAt time.Time
				}
				PageInfo pageInfo
			} `graphql:"stargazers(after: $cursor first: 100)"`
		} `graphql:"repository(name: \"code-golf\" owner: \"code-golf\")"`
	}

	variables := map[string]interface{}{"cursor": (*githubv4.String)(nil)}

	for {
		if err := client.Query(context.Background(), &query, variables); err != nil {
			panic(err)
		}

		limits = append(limits, query.RateLimit)

		for _, edge := range query.Repository.Stargazers.Edges {
			stargazers[edge.Node.DatabaseID] = edge.StarredAt
		}

		if !query.Repository.Stargazers.PageInfo.HasNextPage {
			break
		}

		variables["cursor"] = &query.Repository.Stargazers.PageInfo.EndCursor
	}

	awardTrophies(db, stargazers, "my-god-its-full-of-stars")

	return
}
