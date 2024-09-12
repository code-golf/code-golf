package routes

import (
	"net/http"
	"time"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/null"
	"github.com/code-golf/code-golf/session"
)

// GET /admin
func adminGET(w http.ResponseWriter, r *http.Request) {
	var data struct {
		LastTested []struct {
			Solutions int
			TestedDay time.Time
		}
		Sessions []struct {
			Country  config.NullCountry
			LastUsed time.Time
			Name     string
		}
		Tables []struct {
			Name       null.String
			Rows, Size int
		}
	}

	db := session.Database(r)
	golfer := session.Golfer(r)

	if err := db.Select(
		&data.LastTested,
		` SELECT COUNT(*) solutions, TIMEZONE($1, DATE(tested)) tested_day
		    FROM solutions
		GROUP BY tested_day
		ORDER BY tested_day DESC`,
		golfer.TimeZone,
	); err != nil {
		panic(err)
	}

	if err := db.Select(
		&data.Sessions,
		`WITH grouped_sessions AS (
		    SELECT user_id, MAX(last_used) last_used
		      FROM sessions
		     WHERE user_id != $1
		       AND last_used > TIMEZONE('UTC', NOW()) - INTERVAL '1 hour'
		  GROUP BY user_id
		) SELECT country_flag country, last_used, login name
		    FROM grouped_sessions
		    JOIN users ON id = user_id
		ORDER BY last_used DESC`,
		golfer.ID,
	); err != nil {
		panic(err)
	}

	if err := db.Select(
		&data.Tables,
		` SELECT relname                                                name,
		         CASE WHEN relkind = 'i' THEN 0 ELSE reltuples::int END rows,
		         PG_TOTAL_RELATION_SIZE(c.oid)                          size
		    FROM pg_class     c
		    JOIN pg_namespace n
		      ON n.oid = relnamespace
		     AND nspname = 'public'
		   WHERE reltuples > 0
		   UNION
		  SELECT NULL, 0, PG_DATABASE_SIZE(CURRENT_DATABASE())
		ORDER BY name`,
	); err != nil {
		panic(err)
	}

	render(w, r, "admin/info", data, "Admin Info")
}
