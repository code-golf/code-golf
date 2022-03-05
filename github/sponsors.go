package github

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

func sponsors(db *pgxpool.Pool) (limits []rateLimit) {
	var query struct {
		RateLimit rateLimit
		Viewer    struct {
			SponsorshipsAsMaintainer struct {
				Nodes []struct {
					SponsorEntity struct {
						User struct{ DatabaseID int } `graphql:"... on User"`
					}
				}
			} `graphql:"sponsorshipsAsMaintainer(first: 100)"`
		}
	}

	if err := client.Query(context.Background(), &query, nil); err != nil {
		panic(err)
	}

	limits = append(limits, query.RateLimit)

	tx, err := db.Begin(context.Background())
	if err != nil {
		panic(err)
	}
	defer tx.Rollback(context.Background())

	if _, err := tx.Exec(
		context.Background(),
		"UPDATE users SET sponsor = false",
	); err != nil {
		panic(err)
	}

	for _, node := range query.Viewer.SponsorshipsAsMaintainer.Nodes {
		if _, err := tx.Exec(
			context.Background(),
			"UPDATE users SET sponsor = true WHERE id = $1",
			node.SponsorEntity.User.DatabaseID,
		); err != nil {
			panic(err)
		}
	}

	if err := tx.Commit(context.Background()); err != nil {
		panic(err)
	}

	return
}
