package routes

import (
	"net/http"

	"github.com/code-golf/code-golf/cheevo"
	"github.com/code-golf/code-golf/lang"
	"github.com/code-golf/code-golf/session"
)

// About serves GET /about
func About(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Cheevo *cheevo.Cheevo
		Langs  []lang.Lang
	}{Langs: lang.List}

	if golfer := session.Golfer(r); golfer != nil {
		data.Cheevo = golfer.Earn(session.Database(r), "rtfm")
	}

	render(w, r, "about", data, "About")
}
