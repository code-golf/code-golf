package github

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/shurcooL/graphql"
)

func ideas(db *sqlx.DB) (limits []rateLimit) {
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
					Labels     struct {
						Nodes []struct {
							Name string
						}
					} `graphql:"labels(first: 100)"`
				}
				PageInfo pageInfo
			} `graphql:"issues(after: $cursor first: 100 labels: \"idea\" states: OPEN)"`
		} `graphql:"repository(name: \"code-golf\" owner: \"code-golf\")"`
	}

	variables := map[string]any{"cursor": (*graphql.String)(nil)}

	tx := db.MustBegin()
	defer tx.Rollback()

	tx.MustExec("TRUNCATE ideas")

	for {
		if err := client.Query(context.Background(), &query, variables); err != nil {
			panic(err)
		}

		limits = append(limits, query.RateLimit)

		for _, node := range query.Repository.Issues.Nodes {
			kind := "other"
			for _, label := range node.Labels.Nodes {
				if label.Name == "lang-idea" {
					kind = "lang"
				} else if label.Name == "hole-idea" {
					kind = "hole"
				} else if label.Name == "cheevo-idea" {
					kind = "cheevo"
				}
			}
			tx.MustExec(
				"INSERT INTO ideas VALUES ($1, $2, $3, $4, $5)",
				node.Number,
				node.ThumbsDown.TotalCount,
				node.ThumbsUp.TotalCount,
				node.Title,
				kind,
			)
		}

		if !query.Repository.Issues.PageInfo.HasNextPage {
			break
		}

		variables["cursor"] = &query.Repository.Issues.PageInfo.EndCursor
	}

	if err := tx.Commit(); err != nil {
		panic(err)
	}

	return
}
