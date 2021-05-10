package routes

import (
	"net/http"

	"github.com/code-golf/code-golf/lang"
	"github.com/code-golf/code-golf/session"
	"github.com/code-golf/code-golf/trophy"
)

// About serves GET /about
func About(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Langs  []lang.Lang
		Trophy *trophy.Trophy
	}{Langs: lang.List}

	if golfer := session.Golfer(r); golfer != nil {
		data.Trophy = golfer.Earn(session.Database(r), "rtfm")
	}

	render(w, r, "about", data, "About")
}
