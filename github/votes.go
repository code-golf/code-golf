package github

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/shurcooL/graphql"
)

/* https://docs.github.com/en/graphql/overview/explorer

fragment thumb on ReactionConnection {
  nodes {
    createdAt
    user { databaseId }
  }
} {
  repository(name: "code-golf" owner: "code-golf") {
    issues(first: 100 labels: "idea") {
      nodes {
        thumbsUp:   reactions(content: THUMBS_UP   first: 100) { ...thumb }
        thumbsDown: reactions(content: THUMBS_DOWN first: 100) { ...thumb }
      }
    }
  }
}
*/

func votes(db *sqlx.DB) (limits []rateLimit) {
	voters := map[int]time.Time{}

	type Thumbs struct {
		Nodes []struct {
			CreatedAt time.Time
			User      struct{ DatabaseID int }
		}
	}

	var query struct {
		RateLimit  rateLimit
		Repository struct {
			Issues struct {
				Nodes []struct {
					ThumbsUp   Thumbs `graphql:"thumbsUp:   reactions(content: THUMBS_UP   first: 100)"`
					ThumbsDown Thumbs `graphql:"thumbsDown: reactions(content: THUMBS_DOWN first: 100)"`
				}
				PageInfo pageInfo
			} `graphql:"issues(after: $cursor first: 100 labels: \"idea\")"`
		} `graphql:"repository(name: \"code-golf\" owner: \"code-golf\")"`
	}

	variables := map[string]any{"cursor": (*graphql.String)(nil)}

	for {
		if err := client.Query(context.Background(), &query, variables); err != nil {
			panic(err)
		}

		limits = append(limits, query.RateLimit)

		for _, issue := range query.Repository.Issues.Nodes {
			for _, thumb := range append(issue.ThumbsUp.Nodes, issue.ThumbsDown.Nodes...) {
				if t, ok := voters[thumb.User.DatabaseID]; !ok || thumb.CreatedAt.Before(t) {
					voters[thumb.User.DatabaseID] = thumb.CreatedAt
				}
			}
		}

		if !query.Repository.Issues.PageInfo.HasNextPage {
			break
		}

		variables["cursor"] = &query.Repository.Issues.PageInfo.EndCursor
	}

	awardCheevos(db, voters, "like-comment-subscribe")

	return
}
