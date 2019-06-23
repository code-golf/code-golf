package routes

import (
	"net/http"

	"github.com/JRaspass/code-golf/cookie"
	"github.com/julienschmidt/httprouter"
)

func hole(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := struct {
		HideDetails             bool
		Hole                    Hole
		HoleCssPath, HoleJsPath string
		Langs                   []Lang
		Solutions               map[string]string
	}{
		Hole:        holeByID[r.URL.Path[1:]],
		HoleCssPath: holeCssPath,
		HoleJsPath:  holeJsPath,
		Langs:       langs,
		Solutions:   map[string]string{},
	}

	if c, _ := r.Cookie("hide-details"); c != nil {
		data.HideDetails = true
	}

	if userID, _ := cookie.Read(r); userID != 0 {
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
			var code, lang string

			if err := rows.Scan(&code, &lang); err != nil {
				panic(err)
			}

			data.Solutions[lang] = code
		}

		if err := rows.Err(); err != nil {
			panic(err)
		}
	}

	Render(w, r, http.StatusOK, "hole", data.Hole.Name, data)
}
