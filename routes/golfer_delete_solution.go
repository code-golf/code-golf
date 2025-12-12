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

	if config.AllHoleByID[hole] == nil || config.AllLangByID[lang] == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	golfer := session.Golfer(r)

	tx := session.Database(r).MustBeginTx(r.Context(), nil)
	defer tx.Rollback()

	tx.MustExec(
		"DELETE FROM solutions WHERE hole = $1 AND lang = $2 AND user_id = $3",
		hole, lang, golfer.ID,
	)

	// TODO Add "flash" messages so we can show the cheevo after the redirect.
	golfer.Earn(tx, "rm-rf")

	if err := tx.Commit(); err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/"+hole+"#"+lang, http.StatusSeeOther)
}
