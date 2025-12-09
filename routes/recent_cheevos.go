package routes

import (
	"net/http"
	"time"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/golfer"
	"github.com/code-golf/code-golf/pager"
	"github.com/code-golf/code-golf/session"
)

// GET /recent/cheevos
func recentCheevosGET(w http.ResponseWriter, r *http.Request) {
	var data []struct {
		golfer.GolferLink

		Cheevo *config.Cheevo
		Earned time.Time
	}

	var cheevo *config.Cheevo
	if id := param(r, "cheevo"); id != "all" {
		if cheevo = config.CheevoByID[id]; cheevo == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}

	if err := session.Database(r).Select(
		&data,
		` SELECT avatar_url, country_flag, cheevo, earned, name
		    FROM cheevos
		    JOIN golfers_with_avatars ON id = user_id
		   WHERE cheevo = $1 OR $1 IS NULL
		ORDER BY earned DESC
		   LIMIT $2`,
		cheevo,
		pager.PerPage,
	); err != nil {
		panic(err)
	}

	render(w, r, "recent/cheevos", data, "Recent Achievements")
}
