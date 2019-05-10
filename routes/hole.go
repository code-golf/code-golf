package routes

import (
	"database/sql"
	"net/http"

	"github.com/JRaspass/code-golf/cookie"
	"github.com/julienschmidt/httprouter"
)

func hole(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	type solution struct{ Code, Lang string }

	data := struct {
		Hole                          Hole
		HoleCssPath, HoleJsPath, Lang string
		Langs                         []Lang
		Solutions                     []solution
	}{
		Hole:        holeByID[r.URL.Path[1:]],
		HoleCssPath: holeCssPath,
		HoleJsPath:  holeJsPath,
		Langs:       langs,
	}

	if userID, _ := cookie.Read(r); userID != 0 {
		// Fetch the latest successful lang.
		if err := db.QueryRow(
			"SELECT lang FROM solutions WHERE user_id = $1 AND hole = $2",
			userID, data.Hole.ID,
		).Scan(&data.Lang); err != nil && err != sql.ErrNoRows {
			panic(err)
		}

		// Fetch all the code per lang.
		rows, err := db.Query(
			`SELECT code, lang
			   FROM solutions
			  WHERE hole = $1 AND user_id = $2`,
			data.Hole.ID, userID,
		)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		for rows.Next() {
			var solution solution

			if err := rows.Scan(&solution.Code, &solution.Lang); err != nil {
				panic(err)
			}

			data.Solutions = append(data.Solutions, solution)
		}

		if err := rows.Err(); err != nil {
			panic(err)
		}
	}

	Render(w, r, http.StatusOK, "hole", data)
}
