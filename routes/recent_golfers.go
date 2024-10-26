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
	type golfer struct {
		Cheevos int
		Country config.NullCountry
		Date    time.Time
		Langs   config.Langs
		Name    string
	}

	data := struct {
		Cheevos int
		Golfers []golfer
	}{len(config.CheevoList), make([]golfer, 0, pager.PerPage)}

	if err := session.Database(r).Select(
		&data.Golfers,
		`WITH langs AS (
		    SELECT user_id, array_agg(DISTINCT lang) langs
		      FROM solutions
		     WHERE NOT failing
		  GROUP BY user_id
		), recent AS (
		    SELECT user_id, COUNT(*) cheevos, MIN(earned) date
		      FROM trophies
		  GROUP BY user_id
		  ORDER BY date DESC
		     LIMIT $1
		)  SELECT country_flag country, cheevos, date, login name, langs
		     FROM recent
		     JOIN users ON id = user_id
		LEFT JOIN langs USING (user_id)
		 ORDER BY date DESC`,
		pager.PerPage,
	); err != nil {
		panic(err)
	}

	render(w, r, "recent/golfers", data, "Recent Golfers")
}
