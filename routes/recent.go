package routes

import (
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

type Recent struct {
	Hole      Hole
	Lang      Lang
	Login     string
	Strokes   int
	Submitted time.Time
}

func recent(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	rows, err := db.Query(
		` SELECT hole, lang, login, LENGTH(code), submitted
		    FROM solutions JOIN users ON id = user_id
		ORDER BY submitted DESC LIMIT 100`,
	)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var recents []Recent

	for rows.Next() {
		var holeID, langID string
		var recent Recent

		if err := rows.Scan(
			&holeID,
			&langID,
			&recent.Login,
			&recent.Strokes,
			&recent.Submitted,
		); err != nil {
			panic(err)
		}

		recent.Hole = holeByID[holeID]
		recent.Lang = langByID[langID]

		recents = append(recents, recent)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	Render(w, r, http.StatusOK, "recent", "Recent Solutions", recents)
}
