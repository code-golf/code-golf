package routes

import (
	"net/http"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/golfer"
	"github.com/code-golf/code-golf/session"
)

// GET /about
func aboutGET(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Cheevo                          *config.Cheevo
		FollowLimit, FollowLimitSponsor int
		Langs                           []*config.Lang
	}{
		FollowLimit:        golfer.FollowLimit,
		FollowLimitSponsor: golfer.FollowLimitSponsor,
		Langs:              config.LangList,
	}

	if golfer := session.Golfer(r); golfer != nil {
		data.Cheevo = golfer.Earn(session.Database(r), "rtfm")
	}

	render(w, r, "about", data, "About")
}
