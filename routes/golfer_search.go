package routes

import (
	"net/http"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/session"
)

// GET /golfer/search
func golferSearchGET(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Langs     map[string]string
		Holes     map[string]string
		Solutions []struct {
			Code    string `json:"code"`
			Hole    string `json:"hole"`
			Lang    string `json:"lang"`
			Scoring string `json:"scoring"`
		}
	}{
		Langs: make(map[string]string),
		Holes: make(map[string]string),
	}

	for k, v := range config.AllLangByID {
		data.Langs[k] = v.Name
	}
	for k, v := range config.AllHoleByID {
		data.Holes[k] = v.Name
	}

	if err := session.Database(r).Select(
		&data.Solutions,
		`SELECT code, hole, lang, scoring
		   FROM solutions
		  WHERE user_id = $1`,
		session.Golfer(r).ID,
	); err != nil {
		panic(err)
	}

	render(w, r, "golfer/search", data, "Solution search")
}
