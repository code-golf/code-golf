package routes

import (
	"net/http"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/session"
)

type Solution struct {
	Code    string
	Hole    string
	Lang    string
	Scoring string
}

// GET /golfer/search
func golferSearchGET(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Langs     map[string]string
		Holes     map[string]string
		Solutions []Solution
	}{
		Langs:     make(map[string]string),
		Holes:     make(map[string]string),
		Solutions: []Solution{},
	}

	for k, v := range config.AllLangByID {
		data.Langs[k] = v.Name
	}
	for k, v := range config.AllHoleByID {
		data.Holes[k] = v.Name
	}

	golfer := session.Golfer(r)

	rows, err := session.Database(r).Query(
		`SELECT code, hole, lang, scoring
		   FROM solutions
		  WHERE user_id = $1`,
		golfer.ID,
	)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var code, hole, lang, scoring string

		if err := rows.Scan(&code, &hole, &lang, &scoring); err != nil {
			panic(err)
		}

		data.Solutions = append(data.Solutions, Solution{code, hole, lang, scoring})
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	render(w, r, "golfer/search", data, "Solution search")
}
