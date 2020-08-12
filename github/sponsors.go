package github

import (
	"context"
	"database/sql"
)

func sponsors(db *sql.DB) (limits []rateLimit) {
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

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	defer tx.Rollback()

	if _, err := tx.Exec("UPDATE users SET sponsor = false"); err != nil {
		panic(err)
	}

	for _, node := range query.Viewer.SponsorshipsAsMaintainer.Nodes {
		if _, err := tx.Exec(
			"UPDATE users SET sponsor = true WHERE id = $1",
			node.SponsorEntity.User.DatabaseID,
		); err != nil {
			panic(err)
		}
	}

	if err := tx.Commit(); err != nil {
		panic(err)
	}

	return
}
