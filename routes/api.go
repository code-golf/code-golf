package routes

import (
	_ "embed"
	"net/http"

	"github.com/code-golf/code-golf/session"
)

//go:embed api.yml
var yml []byte

// API serves GET /api
func API(w http.ResponseWriter, r *http.Request) { w.Write(yml) }

// APISuggestionsGolfers serves GET /api/suggestions/golfers
func APISuggestionsGolfers(w http.ResponseWriter, r *http.Request) {
	var json []byte

	if err := session.Database(r).QueryRow(
		`WITH golfers AS (
		    SELECT login
		      FROM users
		     WHERE strpos(login, $1) > 0 AND login != $2
		  ORDER BY login
		     LIMIT 10
		) SELECT COALESCE(json_agg(login), '[]') FROM golfers`,
		r.FormValue("q"),
		r.FormValue("ignore"),
	).Scan(&json); err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(json); err != nil {
		panic(err)
	}
}
