package routes

import (
	"database/sql"
	"net/http"

	"github.com/code-golf/code-golf/session"
)

// GET /golfer/export
func golferExportGET(w http.ResponseWriter, r *http.Request) {
	golfer := session.Golfer(r)

	rows, err := session.Database(r).Query(
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
	)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var json sql.RawBytes

		if err := rows.Scan(&json); err != nil {
			panic(err)
		}

		w.Header().Set("Content-Disposition",
			`attachment; filename="code-golf-export.json"`)
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	golfer.Earn(session.Database(r), "takeout")
}
