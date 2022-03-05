package routes

import (
	"net/http"

	"github.com/code-golf/code-golf/session"
)

// GolferExport serves GET /golfer/export
func GolferExport(w http.ResponseWriter, r *http.Request) {
	db := session.Database(r)
	golfer := session.Golfer(r)

	var json []byte
	if err := db.QueryRow(
		r.Context(),
		`WITH cheevos AS (
		    SELECT trophy cheevo, to_json(earned)#>>'{}' || 'Z' earned
		      FROM trophies
		     WHERE user_id = $1
		  ORDER BY trophy
		), solutions AS (
		    SELECT hole, lang, scoring, bytes, chars, failing,
		           to_json(submitted)#>>'{}' || 'Z' submitted, code
		      FROM solutions
		     WHERE user_id = $1
		  ORDER BY hole, lang, scoring
		) SELECT json_build_object(
		    'name',      login,
		    'country',   country,
		    'time_zone', time_zone,
		    'cheevos',   (SELECT COALESCE(json_agg(cheevos  ), '[]'::json) FROM cheevos  ),
		    'solutions', (SELECT COALESCE(json_agg(solutions), '[]'::json) FROM solutions)
		) FROM users WHERE id = $1`,
		golfer.ID,
	).Scan(&json); err != nil {
		panic(err)
	}

	w.Header().Set("Content-Disposition",
		`attachment; filename="code-golf-export.json"`)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)

	golfer.Earn(db, "takeout")
}
