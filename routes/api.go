package routes

import (
	"net/http"

	"github.com/code-golf/code-golf/session"
)

// APISuggestionsGolfers serves GET /api/v1/suggestions/golfers
func APISuggestionsGolfers(w http.ResponseWriter, r *http.Request) {
	var json []byte

	if err := session.Database(r).QueryRow(
		`SELECT COALESCE(json_agg(login ORDER BY login), '[]')
		   FROM users
		  WHERE strpos(login, $1) > 0 AND login != $2`,
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
