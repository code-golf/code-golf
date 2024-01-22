package routes

import (
	"cmp"
	"net/http"
	"slices"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/session"
)

// GET /stats
func statsGET(w http.ResponseWriter, r *http.Request) {
	var data struct{ Bytes, Golfers, Solutions int }

	db := session.Database(r)

	if err := db.QueryRow(
		"SELECT COUNT(DISTINCT user_id) FROM trophies",
	).Scan(&data.Golfers); err != nil {
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

// GET /stats/{page:countries}
func statsCountriesGET(w http.ResponseWriter, r *http.Request) {
	type row struct {
		Country       config.NullCountry
		Golfers, Rank int
		Percent       string
	}

	var data []row
	if err := session.Database(r).Select(
		&data,
		` SELECT RANK() OVER (ORDER BY COUNT(*) DESC)             rank,
		         COALESCE(country, '')                            country,
		         COUNT(*)                                         golfers,
		         ROUND(COUNT(*) / SUM(COUNT(*)) OVER () * 100, 2) percent
		    FROM users
		GROUP BY COALESCE(country, '')`,
	); err != nil {
		panic(err)
	}

	// Sort in Go because SQL doesn't have access to country name, just ID.
	slices.SortFunc(data, func(a, b row) int {
		if c := cmp.Compare(a.Rank, b.Rank); c != 0 {
			return c
		}

		var aName, bName string
		if a.Country.Valid {
			aName = a.Country.Country.Name
		}
		if b.Country.Valid {
			bName = b.Country.Country.Name
		}
		return cmp.Compare(aName, bName)
	})

	render(w, r, "stats", data, "Statistics: Countries")
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
