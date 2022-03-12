package routes

import (
	"net/http"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/session"
)

// GET /about
func aboutGET(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Cheevo *config.Cheevo
		Langs  []*config.Lang
	}{Langs: config.LangList}

	if golfer := session.Golfer(r); golfer != nil {
		data.Cheevo = golfer.Earn(session.Database(r), "rtfm")
	}

	render(w, r, "about", data, "About")
}
