package routes

import (
	"net/http"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/session"
)

// Hole serves GET /{hole}
func Hole(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	data := struct {
		Authors     []string
		HideDetails bool
		Hole        *config.Hole
		Solutions   []map[string]string
	}{
		Solutions: []map[string]string{{}, {}},
	}

	var ok bool
	if data.Hole, ok = config.HoleByID[param(r, "hole")]; !ok {
		if data.Hole, ok = config.ExpHoleByID[param(r, "hole")]; !ok {
			NotFound(w, r)
			return
		}
	}

	if c, _ := r.Cookie("hide-details"); c != nil {
		data.HideDetails = true
	}

	// Lookup the hole's author(s).
	if data.Hole.Experiment == 0 {
		if err := session.Database(r).QueryRow(
			ctx,
			`SELECT array_agg(login ORDER BY login)::text[]
			   FROM authors
			   JOIN users ON id = user_id
			  WHERE hole = $1`,
			data.Hole.ID,
		).Scan(&data.Authors); err != nil {
			panic(err)
		}
	}

	if golfer := session.Golfer(r); golfer != nil && data.Hole.Experiment == 0 {
		// Fetch all the code per lang.
		rows, err := session.Database(r).Query(
			ctx,
			`SELECT code, lang, scoring
			   FROM solutions
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
	ctx := r.Context()
	data := struct {
		Authors     []string
		HideDetails bool
		Hole        *config.Hole
		Langs       []*config.Lang
		Solutions   map[string]map[string]string
	}{
		Langs:     config.LangList,
		Solutions: map[string]map[string]string{},
	}

	var ok bool
	if data.Hole, ok = config.HoleByID[param(r, "hole")]; !ok {
		if data.Hole, ok = config.ExpHoleByID[param(r, "hole")]; !ok {
			NotFound(w, r)
			return
		}
	}

	if c, _ := r.Cookie("hide-details"); c != nil {
		data.HideDetails = true
	}

	// Lookup the hole's author(s).
	if data.Hole.Experiment == 0 {
		if err := session.Database(r).QueryRow(
			ctx,
			`SELECT array_agg(login ORDER BY login)::text[]
			   FROM authors
			   JOIN users ON id = user_id
			  WHERE hole = $1`,
			data.Hole.ID,
		).Scan(&data.Authors); err != nil {
			panic(err)
		}
	}

	if golfer := session.Golfer(r); golfer != nil && data.Hole.Experiment == 0 {
		// Fetch all the code per lang.
		rows, err := session.Database(r).Query(
			ctx,
			`SELECT code, lang, scoring
			   FROM solutions
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
