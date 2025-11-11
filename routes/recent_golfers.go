package routes

import (
	"net/http"
	"time"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/pager"
	"github.com/code-golf/code-golf/session"
)

// GET /recent/golfers
func recentGolfersGET(w http.ResponseWriter, r *http.Request) {
	var data []struct {
		Cheevos config.Cheevos
		Country *config.Country
		Date    time.Time
		Langs   config.Langs
		Name    string
	}

	if err := session.Database(r).Select(
		&data,
		`WITH langs AS (
		    SELECT user_id, array_agg(DISTINCT lang) langs
		      FROM solutions
		     WHERE NOT failing
		  GROUP BY user_id
		), cheevos AS (
		    SELECT user_id, array_agg(cheevo ORDER BY cheevo) cheevos, MIN(earned) date
		      FROM cheevos
		  GROUP BY user_id
		  ORDER BY date DESC
		     LIMIT $1
		)  SELECT country_flag country, cheevos, date, login name, langs
		     FROM cheevos
		     JOIN users ON id = user_id
		LEFT JOIN langs USING (user_id)
		 ORDER BY date DESC`,
		pager.PerPage,
	); err != nil {
		panic(err)
	}

	render(w, r, "recent/golfers", data, "Recent Golfers")
}
