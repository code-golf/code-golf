package main

import (
	"database/sql"
	"net/http"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
	var err error

	if db, err = sql.Open("postgres", os.Getenv("DB_DSN")); err != nil {
		panic(err)
	}
}

func addSolution(userID int, lang, code string) {
	if _, err := db.Exec(
		`INSERT INTO solutions(user_id, lang, code) VALUES($1, $2, $3)
		 ON CONFLICT ON CONSTRAINT solutions_pkey DO UPDATE SET code = $3`,
		userID, lang, code,
	); err != nil {
		panic(err)
	}
}

func addUser(id int, login string) {
	if _, err := db.Exec(
		`INSERT INTO users VALUES($1, $2)
		 ON CONFLICT(id) DO UPDATE SET login = $2`,
		id, login,
	); err != nil {
		panic(err)
	}
}

func printLeaderboards(w http.ResponseWriter) {
	rows, err := db.Query(
		`SELECT login, lang, length(code)
		   FROM solutions
		   JOIN users on user_id = id`,
	)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	w.Write([]byte("<table>"))

	i := 0

	for rows.Next() {
		var login, lang string
		var length int

		if err := rows.Scan(&login, &lang, &length); err != nil {
			panic(err)
		}

		i++

		w.Write([]byte(
			"<tr><td>" + strconv.Itoa(i)
			"<td>" + strconv.Itoa(length) +
			"<td>" + lang +
			"<td>" + login,
		))
	}

	w.Write([]byte("</table>"))

	if err := rows.Err(); err != nil {
		panic(err)
	}
}
