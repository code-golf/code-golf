package github

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/shurcooL/githubv4"
)

func ideas(db *pgxpool.Pool) (limits []rateLimit) {
	type thumbs struct{ TotalCount int }

	var query struct {
		RateLimit  rateLimit
		Repository struct {
			Issues struct {
				Nodes []struct {
					Number     int
					ThumbsDown thumbs `graphql:"thumbsDown: reactions(content: THUMBS_DOWN)"`
					ThumbsUp   thumbs `graphql:"thumbsUp:   reactions(content: THUMBS_UP)"`
					Title      string
				}
				PageInfo pageInfo
			} `graphql:"issues(after: $cursor first: 100 labels: \"idea\" states: OPEN)"`
		} `graphql:"repository(name: \"code-golf\" owner: \"code-golf\")"`
	}

	variables := map[string]interface{}{"cursor": (*githubv4.String)(nil)}

	tx, err := db.Begin(context.Background())
	if err != nil {
		panic(err)
	}

	defer tx.Rollback(context.Background())

	if _, err := tx.Exec(context.Background(), "TRUNCATE ideas"); err != nil {
		panic(err)
	}

	for {
		if err := client.Query(context.Background(), &query, variables); err != nil {
			panic(err)
		}

		limits = append(limits, query.RateLimit)

		for _, node := range query.Repository.Issues.Nodes {
			if _, err := tx.Exec(
				context.Background(),
				"INSERT INTO ideas VALUES ($1, $2, $3, $4)",
				node.Number,
				node.ThumbsDown.TotalCount,
				node.ThumbsUp.TotalCount,
				node.Title,
			); err != nil {
				panic(err)
			}
		}

		if !query.Repository.Issues.PageInfo.HasNextPage {
			break
		}

		variables["cursor"] = &query.Repository.Issues.PageInfo.EndCursor
	}

	if err := tx.Commit(context.Background()); err != nil {
		panic(err)
	}

	return
}
