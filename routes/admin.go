package routes

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/code-golf/code-golf/session"
)

// Admin serves GET /admin
func Admin(w http.ResponseWriter, r *http.Request) {
	type Session struct {
		Golfer   string
		LastUsed time.Time
	}

	type Table struct {
		Name       sql.NullString
		Rows, Size int
	}

	type TimeZone struct {
		Name    string
		Golfers int
		Percent float64
	}

	data := struct {
		Sessions  []Session
		Tables    []Table
		TimeZones []TimeZone
	}{}

	db := session.Database(r)

	rows, err := db.Query(
		` SELECT login, MAX(last_used)
		    FROM sessions JOIN users ON user_id = users.id
		   WHERE user_id != $1
		     AND last_used > TIMEZONE('UTC', NOW()) - INTERVAL '1 day'
		GROUP BY login
		ORDER BY max DESC`,
		session.Golfer(r).ID,
	)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		var session Session

		if err := rows.Scan(&session.Golfer, &session.LastUsed); err != nil {
			panic(err)
		}

		data.Sessions = append(data.Sessions, session)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	rows, err = db.Query(
		` SELECT relname,
		         CASE WHEN relkind = 'i' THEN 0 ELSE reltuples END,
		         PG_TOTAL_RELATION_SIZE(c.oid)
		    FROM pg_class     c
		    JOIN pg_namespace n
		      ON n.oid = relnamespace
		     AND nspname = 'public'
		   WHERE reltuples != 0
		   UNION
		  SELECT NULL, 0, PG_DATABASE_SIZE('code-golf')
		ORDER BY relname`,
	)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		var table Table

		if err := rows.Scan(&table.Name, &table.Rows, &table.Size); err != nil {
			panic(err)
		}

		data.Tables = append(data.Tables, table)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	rows, err = db.Query(
		`SELECT COALESCE(time_zone, ''), COUNT(*), COUNT(*) / SUM(COUNT(*)) OVER () * 100
		   FROM users GROUP BY time_zone ORDER BY COUNT(*) DESC, COALESCE(time_zone, '')`,
	)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		var timeZone TimeZone

		if err := rows.Scan(&timeZone.Name, &timeZone.Golfers, &timeZone.Percent); err != nil {
			panic(err)
		}

		data.TimeZones = append(data.TimeZones, timeZone)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	render(w, r, "admin/info", "Admin Info", data)
}
