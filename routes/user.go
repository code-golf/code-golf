package routes

import (
	"database/sql"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func user(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var one int

	switch err := db.QueryRow(
		"SELECT 1 FROM users WHERE login = $1", ps[0].Value,
	).Scan(&one); err {
	case sql.ErrNoRows:
		print404(w, r)
	case nil:
		printHeader(w, r, 200)
		w.Write([]byte("<article><h1>" + ps[0].Value + "</h1>"))
	default:
		panic(err)
	}
}
