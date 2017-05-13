package main

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

var dbDSN = os.Getenv("DB_DSN")

func addSolution(userID int, lang, code string ) {
	db, err := sql.Open("postgres", dbDSN)
	if err != nil {
		panic(err)
	}

	if _, err := db.Exec(
		`INSERT INTO solutions(user_id, lang, code) VALUES($1, $2, $3)
		 ON CONFLICT ON CONSTRAINT solutions_pkey DO UPDATE SET code = $3`,
		userID, lang, code,
	); err != nil {
		panic(err)
	}
}

func addUser(id int, login string) {
	if db, err := sql.Open("postgres", dbDSN); err != nil {
		panic(err)
	} else if _, err := db.Exec(
		`INSERT INTO users VALUES($1, $2)
		 ON CONFLICT(id) DO UPDATE SET login = $2`,
		id, login,
	); err != nil {
		panic(err)
	}
}
