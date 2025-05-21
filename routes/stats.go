package routes

import (
	"cmp"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/session"
)

// GET /stats
func statsGET(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Bytes, Cheevos, CheevosEarned, Countries, Golfers,
		Holes, HolesExp, Langs, LangsExp, Solutions int
	}{
		Cheevos:  len(config.CheevoList),
		Holes:    len(config.AllHoleList),
		HolesExp: len(config.ExpHoleList),
		Langs:    len(config.AllLangList),
		LangsExp: len(config.ExpLangList),
	}

	db := session.Database(r)

	if err := db.QueryRow(
		"SELECT COUNT(*), COUNT(DISTINCT user_id) FROM trophies",
	).Scan(&data.CheevosEarned, &data.Golfers); err != nil {
		panic(err)
	}

	if err := db.QueryRow(
		"SELECT COUNT(DISTINCT country) FROM users WHERE country IS NOT NULL",
	).Scan(&data.Countries); err != nil {
		panic(err)
	}

	if err := db.QueryRow(
		`SELECT COALESCE(SUM(bytes), 0), COUNT(*)
		   FROM solutions
		  WHERE NOT failing`,
	).Scan(&data.Bytes, &data.Solutions); err != nil {
		panic(err)
	}

	render(w, r, "stats", data, "Statistics")
}

// GET /stats/{page:cheevos}
func statsCheevosGET(w http.ResponseWriter, r *http.Request) {
	type row struct {
		Cheevo        *config.Cheevo
		Golfers, Rank int
		Percent       string
	}

	var data []row
	if err := session.Database(r).Select(
		&data,
		` SELECT RANK() OVER (ORDER BY COUNT(*) DESC)             rank,
		         trophy                                           cheevo,
		         COUNT(*)                                         golfers,
		         ROUND(COUNT(*) / SUM(COUNT(*)) OVER () * 100, 2) percent
		    FROM trophies
		GROUP BY trophy`,
	); err != nil {
		panic(err)
	}

	render(w, r, "stats", data, "Statistics: Achievements")
}

// GET /stats/{page:countries}
func statsCountriesGET(w http.ResponseWriter, r *http.Request) {
	type row struct {
		Country       *config.Country
		Golfers, Rank int
		Percent       string
	}

	var data []row
	if err := session.Database(r).Select(
		&data,
		` SELECT RANK() OVER (ORDER BY COUNT(*) DESC)             rank,
		         country                                          country,
		         COUNT(*)                                         golfers,
		         ROUND(COUNT(*) / SUM(COUNT(*)) OVER () * 100, 2) percent
		    FROM users
		GROUP BY country`,
	); err != nil {
		panic(err)
	}

	// Sort in Go because SQL doesn't have access to country name, just ID.
	slices.SortFunc(data, func(a, b row) int {
		if c := cmp.Compare(a.Rank, b.Rank); c != 0 {
			return c
		}

		var aName, bName string
		if a.Country != nil {
			aName = a.Country.Name
		}
		if b.Country != nil {
			bName = b.Country.Name
		}
		return cmp.Compare(aName, bName)
	})

	render(w, r, "stats", data, "Statistics: Countries")
}

// GET /stats/{page:golfers}
func statsGolfersGET(w http.ResponseWriter, r *http.Request) {
	var data []struct {
		Count, Sum int
		Date       time.Time
		PerDay     string
	}

	if err := session.Database(r).Select(
		&data,
		`-- Only consider golfers that have a cheevo (i.e here to stay).
		WITH earnt_golfers AS (
		    SELECT DISTINCT id, started FROM users JOIN trophies ON id = user_id
		), first_golfer AS (
		    SELECT started date, 1 count, 1 sum
		      FROM earnt_golfers
		  ORDER BY started
		     LIMIT 1
		), counts AS (
		    SELECT LEAST(
		               DATE_TRUNC('year', started)
		                   + INTERVAL '1 year' - INTERVAL '1 microsecond',
		               TIMEZONE('UTC', NOW())
		           ) date,
		           COUNT(*)
		      FROM earnt_golfers
		  GROUP BY date
		) SELECT *, 0 per_day FROM first_golfer
		   UNION ALL
		  SELECT *, SUM(count) OVER (ORDER BY date),
		         ROUND(count / EXTRACT(doy FROM date), 2) per_day
		    FROM counts`,
	); err != nil {
		panic(err)
	}

	render(w, r, "stats", data, "Statistics: Golfers")
}

// GET /stats/{page:holes|langs}
func statsTableGET(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Fact string
		Rows []struct {
			Count, Golfers, Rank int
			PerGolfer            string
			Hole                 *config.Hole
			Lang                 *config.Lang
		}
	}

	column := ""
	switch param(r, "page") {
	case "holes":
		column = "hole"
		data.Fact = "Hole"
	case "langs":
		column = "lang"
		data.Fact = "Language"
	}

	if err := session.Database(r).Select(
		&data.Rows,
		` SELECT RANK() OVER (ORDER BY COUNT(*) DESC, `+column+`),
		         `+column+`,
		         COUNT(*),
		         COUNT(DISTINCT user_id) golfers,
		         ROUND(COUNT(*)::decimal / COUNT(DISTINCT user_id), 2) per_golfer
		    FROM solutions
		   WHERE NOT failing
		GROUP BY `+column,
	); err != nil {
		panic(err)
	}

	render(w, r, "stats", data, "Statistics: "+data.Fact+"s")
}

// GET /stats/{page:unsolved-holes}
func statsUnsolvedHolesGET(w http.ResponseWriter, r *http.Request) {
	type holeLang struct {
		Hole *config.Hole
		Lang *config.Lang
	}

	var data []holeLang
	if err := session.Database(r).Select(
		&data,
		`WITH solves AS (
		    SELECT DISTINCT hole, lang FROM stable_passing_solutions
		),
		holes AS (SELECT id hole FROM holes WHERE experiment = 0),
		langs AS (SELECT id lang FROM langs WHERE experiment = 0),
		combo AS (SELECT hole, lang FROM holes CROSS JOIN langs)
		   SELECT hole, lang
		     FROM combo
		LEFT JOIN solves USING (hole, lang)
		    WHERE solves.hole IS NULL`,
	); err != nil {
		panic(err)
	}

	slices.SortFunc(data, func(a, b holeLang) int {
		return cmp.Or(
			cmp.Compare(strings.ToLower(a.Hole.Name), strings.ToLower(b.Hole.Name)),
			cmp.Compare(strings.ToLower(a.Lang.Name), strings.ToLower(b.Lang.Name)),
		)
	})

	render(w, r, "stats", data, "Statistics: Unsolved Holes")
}
