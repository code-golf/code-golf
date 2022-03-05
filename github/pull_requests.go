package github

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/shurcooL/githubv4"
)

func pullRequests(db *pgxpool.Pool) (limits []rateLimit) {
	pullRequests := map[int]time.Time{}

	var query struct {
		RateLimit  rateLimit
		Repository struct {
			PullRequests struct {
				Nodes []struct {
					Author struct {
						User struct{ DatabaseID int } `graphql:"... on User"`
					}
					MergedAt time.Time
				}
				PageInfo pageInfo
			} `graphql:"pullRequests(after: $cursor first: 100 states: MERGED)"`
		} `graphql:"repository(name: \"code-golf\" owner: \"code-golf\")"`
	}

	variables := map[string]interface{}{"cursor": (*githubv4.String)(nil)}

	for {
		if err := client.Query(context.Background(), &query, variables); err != nil {
			panic(err)
		}

		limits = append(limits, query.RateLimit)

		for _, node := range query.Repository.PullRequests.Nodes {
			mergedAt, ok := pullRequests[node.Author.User.DatabaseID]
			if !ok || node.MergedAt.Before(mergedAt) {
				pullRequests[node.Author.User.DatabaseID] = node.MergedAt
			}
		}

		if !query.Repository.PullRequests.PageInfo.HasNextPage {
			break
		}

		variables["cursor"] = &query.Repository.PullRequests.PageInfo.EndCursor
	}

	awardCheevos(db, pullRequests, "patches-welcome")

	return
}
