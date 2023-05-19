package routes

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/session"
	"golang.org/x/exp/slices"
)

// GET /admin
func adminGET(w http.ResponseWriter, r *http.Request) {
	type country struct {
		Flag, ID, Name string
		Golfers        int
		Percent        float64
	}

	data := struct {
		Countries []*country
		Sessions  []struct {
			Country  config.NullCountry
			LastUsed time.Time
			Name     string
		}
		Tables []struct {
			Name       sql.NullString
			Rows, Size int
		}
	}{}

	db := session.Database(r)

	if err := db.Select(
		&data.Sessions,
		`WITH grouped_sessions AS (
		    SELECT user_id, MAX(last_used) last_used
		      FROM sessions
		     WHERE user_id != $1
		       AND last_used > TIMEZONE('UTC', NOW()) - INTERVAL '1 day'
		  GROUP BY user_id
		) SELECT country_flag country, last_used, login name
		    FROM grouped_sessions
		    JOIN users ON id = user_id
		ORDER BY last_used DESC`,
		session.Golfer(r).ID,
	); err != nil {
		panic(err)
	}

	if err := db.Select(
		&data.Tables,
		` SELECT relname                                           name,
		         CASE WHEN relkind = 'i' THEN 0 ELSE reltuples END rows,
		         PG_TOTAL_RELATION_SIZE(c.oid)                     size
		    FROM pg_class     c
		    JOIN pg_namespace n
		      ON n.oid = relnamespace
		     AND nspname = 'public'
		   WHERE reltuples != 0
		   UNION
		  SELECT NULL, 0, PG_DATABASE_SIZE('code-golf')
		ORDER BY name`,
	); err != nil {
		panic(err)
	}

	if err := db.Select(
		&data.Countries,
		` SELECT COALESCE(country, '')                  id,
		         COUNT(*)                               golfers,
		         COUNT(*) / SUM(COUNT(*)) OVER () * 100 percent
		    FROM users
		GROUP BY COALESCE(country, '')`,
	); err != nil {
		panic(err)
	}

	for _, c := range data.Countries {
		if country, ok := config.CountryByID[c.ID]; ok {
			c.Flag = country.Flag
			c.Name = country.Name
		}
	}

	slices.SortStableFunc(data.Countries, func(a, b *country) bool {
		if a.Golfers != b.Golfers {
			return a.Golfers > b.Golfers
		}

		return a.Name < b.Name
	})

	render(w, r, "admin/info", data, "Admin Info")
}
