package github

import (
	"context"
	"database/sql"
	"log"
)

// Updates the usernames of every GitHub connection.
func updateUsernames(db *sql.DB) (limits []rateLimit) {
	const batchSize = 100 // GitHub limits us to 100 nodes per query.

	var ids []string

	rows, err := db.Query(
		`SELECT encode(('04:User' || id)::bytea, 'base64')
		   FROM connections WHERE connection = 'github'`,
	)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			panic(err)
		}

		ids = append(ids, id)

		if len(ids) == batchSize {
			if limit, err := updateUsernamesBatch(db, ids); err != nil {
				panic(err)
			} else {
				limits = append(limits, *limit)
			}

			ids = []string{}
		}
	}

	if len(ids) > 0 {
		if limit, err := updateUsernamesBatch(db, ids); err != nil {
			panic(err)
		} else {
			limits = append(limits, *limit)
		}
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	return
}

func updateUsernamesBatch(db *sql.DB, ids []string) (*rateLimit, error) {
	var query struct {
		RateLimit rateLimit
		Nodes     []struct {
			User struct {
				DatabaseID int
				Login      string
			} `graphql:"... on User"`
		} `graphql:"nodes(ids: $ids)"`
	}

	if err := client.Query(
		context.Background(),
		&query,
		map[string]any{"ids": ids},
	); err != nil {
		// Log GraphQL errors, they're not fatal.
		log.Print(err)
	}

	updateConn, err := db.Prepare(
		`WITH old AS (
		    SELECT username FROM connections
		     WHERE connection = 'github' AND id = $1
		)  UPDATE connections SET username = $2
		    WHERE connection = 'github' AND id = $1
		RETURNING (SELECT username FROM old)`,
	)
	if err != nil {
		return nil, err
	}
	defer updateConn.Close()

	updateUser, err := db.Prepare("UPDATE users SET login = $2 WHERE id = $1")
	if err != nil {
		return nil, err
	}
	defer updateUser.Close()

	for _, node := range query.Nodes {
		// Skip unresolved nodes.
		if node.User.DatabaseID == 0 {
			continue
		}

		var username string
		if err := updateConn.QueryRow(
			node.User.DatabaseID, node.User.Login,
		).Scan(&username); err != nil {
			return nil, err
		}

		if username != node.User.Login {
			log.Printf("%8d: %s → %s\n",
				node.User.DatabaseID, username, node.User.Login)

			// Also update users, ATM they match the GitHub connection.
			if _, err := updateUser.Exec(
				node.User.DatabaseID, node.User.Login,
			); err != nil {
				return nil, err
			}
		}
	}

	return &query.RateLimit, nil
}
