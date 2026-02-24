package github

import (
	"context"

	"github.com/jmoiron/sqlx"
)

func sponsors(db *sqlx.DB) (limits []rateLimit, err error) {
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
		return nil, err
	}

	limits = append(limits, query.RateLimit)

	tx := db.MustBegin()
	defer tx.Rollback()

	tx.MustExec("UPDATE users SET sponsor = false")

	for _, node := range query.Viewer.SponsorshipsAsMaintainer.Nodes {
		tx.MustExec(
			"UPDATE users SET sponsor = true WHERE id = $1",
			node.SponsorEntity.User.DatabaseID,
		)
	}

	return limits, tx.Commit()
}
