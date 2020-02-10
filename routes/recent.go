package routes

import (
	"net/http"
	"time"
)

// Recent serves GET /recent
func Recent(w http.ResponseWriter, r *http.Request) {
	rows, err := db(r).Query(
		` SELECT hole, lang, login, LENGTH(code), submitted
		    FROM solutions JOIN users ON id = user_id
		ORDER BY submitted DESC LIMIT 100`,
	)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	type recent struct {
		Hole      Hole
		Lang      Lang
		Login     string
		Strokes   int
		Submitted time.Time
	}

	var recents []recent

	for rows.Next() {
		var holeID, langID string
		var r recent

		if err := rows.Scan(
			&holeID,
			&langID,
			&r.Login,
			&r.Strokes,
			&r.Submitted,
		); err != nil {
			panic(err)
		}

		r.Hole = holeByID[holeID]
		r.Lang = langByID[langID]

		recents = append(recents, r)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	render(w, r, http.StatusOK, "recent", "Recent Solutions", recents)
}
