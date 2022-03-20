package routes

import (
	"net/http"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/session"
)

// POST /golfers/delete/solution
func golferDeleteSolutionPOST(w http.ResponseWriter, r *http.Request) {
	hole := r.PostFormValue("hole")
	lang := r.PostFormValue("lang")

	if config.HoleByID[hole] == nil || config.LangByID[lang] == nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	if _, err := session.Database(r).Exec(
		"DELETE FROM solutions WHERE hole = $1 AND lang = $2 AND user_id = $3",
		hole, lang, session.Golfer(r).ID,
	); err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/"+hole+"#"+lang, http.StatusSeeOther)
}
