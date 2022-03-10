package routes

import (
	"net/http"
	"time"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/pager"
	"github.com/code-golf/code-golf/session"
	"github.com/lib/pq"
)

// RecentGolfers serves GET /recent/golfers
func RecentGolfers(w http.ResponseWriter, r *http.Request) {
	type golfer struct {
		Cheevos       int
		Country, Name string
		Date          time.Time
		Langs         []string
	}

	data := struct {
		Cheevos int
		Golfers []golfer
	}{len(config.CheevoList), make([]golfer, 0, pager.PerPage)}

	rows, err := session.Database(r).Query(
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
		)  SELECT country_flag, cheevos, date, login, langs
		     FROM recent
		     JOIN users ON id = user_id
		LEFT JOIN langs USING (user_id)
		 ORDER BY date DESC`,
		pager.PerPage,
	)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var g golfer

		if err := rows.Scan(
			&g.Country,
			&g.Cheevos,
			&g.Date,
			&g.Name,
			pq.Array(&g.Langs),
		); err != nil {
			panic(err)
		}

		data.Golfers = append(data.Golfers, g)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	render(w, r, "recent/golfers", data, "Recent Golfers")
}
