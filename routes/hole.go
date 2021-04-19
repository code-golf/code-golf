package routes

import (
	"net/http"

	"github.com/code-golf/code-golf/hole"
	"github.com/code-golf/code-golf/lang"
	"github.com/code-golf/code-golf/session"
)

// Hole serves GET /{hole}
func Hole(w http.ResponseWriter, r *http.Request) {
	data := struct {
		HideDetails bool
		Hole        hole.Hole
		Langs       []lang.Lang
		Solutions   []map[string]string
	}{
		Langs:     lang.List,
		Solutions: []map[string]string{{}, {}},
	}

	var ok bool
	if data.Hole, ok = hole.ByID[param(r, "hole")]; !ok {
		if data.Hole, ok = hole.ExperimentalByID[param(r, "hole")]; !ok {
			NotFound(w, r)
			return
		}
	}

	if c, _ := r.Cookie("hide-details"); c != nil {
		data.HideDetails = true
	}

	if golfer := session.Golfer(r); golfer != nil && data.Hole.Experiment == 0 {
		// Fetch all the code per lang.
		rows, err := session.Database(r).Query(
			`SELECT code, lang, scoring
				    FROM solutions
				    JOIN code ON code_id = id
				  WHERE hole = $1 AND user_id = $2`,
			data.Hole.ID, golfer.ID,
		)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		for rows.Next() {
			var code, lang, scoring string

			if err := rows.Scan(&code, &lang, &scoring); err != nil {
				panic(err)
			}

			solution := 0
			if scoring == "chars" {
				solution = 1
			}

			data.Solutions[solution][lang] = code
		}

		if err := rows.Err(); err != nil {
			panic(err)
		}
	}

	render(w, r, "hole", data, data.Hole.Name)
}

// HoleNG serves GET /ng/{hole}
func HoleNG(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Hole      hole.Hole
		Langs     map[string]lang.Lang
		Solutions map[string]map[string]string
	}{
		Langs:     lang.ByID,
		Solutions: map[string]map[string]string{},
	}

	var ok bool
	if data.Hole, ok = hole.ByID[param(r, "hole")]; !ok {
		if data.Hole, ok = hole.ExperimentalByID[param(r, "hole")]; !ok {
			NotFound(w, r)
			return
		}
	}

	if golfer := session.Golfer(r); golfer != nil && data.Hole.Experiment == 0 {
		// Fetch all the code per lang.
		rows, err := session.Database(r).Query(
			`SELECT code, lang, scoring
			   FROM solutions
			   JOIN code ON code_id = id
			  WHERE hole = $1 AND user_id = $2`,
			data.Hole.ID, golfer.ID,
		)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		for rows.Next() {
			var code, lang, scoring string

			if err := rows.Scan(&code, &lang, &scoring); err != nil {
				panic(err)
			}

			if data.Solutions[lang] == nil {
				data.Solutions[lang] = map[string]string{scoring: code}
			} else {
				data.Solutions[lang][scoring] = code
			}
		}

		if err := rows.Err(); err != nil {
			panic(err)
		}
	}

	render(w, r, "hole-ng", data, data.Hole.Name)
}
