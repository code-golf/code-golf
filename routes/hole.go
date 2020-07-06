package routes

import (
	"net/http"

	"github.com/code-golf/code-golf/cookie"
	"github.com/code-golf/code-golf/hole"
	"github.com/code-golf/code-golf/lang"
	"github.com/code-golf/code-golf/middleware"
)

// Hole serves GET /{hole}
func Hole(w http.ResponseWriter, r *http.Request) {
	data := struct {
		HideDetails             bool
		Hole                    hole.Hole
		HoleCssPath, HoleJsPath string
		Langs                   []lang.Lang
		Solutions               map[string]string
	}{
		HoleCssPath: holeCssPath,
		HoleJsPath:  holeJsPath,
		Langs:       lang.List,
		Solutions:   map[string]string{},
	}

	var ok bool
	if data.Hole, ok = hole.ByID[param(r, "hole")]; !ok {
		render(w, r, "404", "", nil)
		return
	}

	if c, _ := r.Cookie("hide-details"); c != nil {
		data.HideDetails = true
	}

	if userID, _ := cookie.Read(r); userID != 0 {
		// Fetch all the code per lang.
		rows, err := middleware.Database(r).Query(
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

	render(w, r, "hole", data.Hole.Name, data)
}
