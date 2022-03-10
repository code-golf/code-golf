package routes

import (
	"database/sql"
	"net/http"
	"sort"
	"time"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/session"
)

// Admin serves GET /admin
func Admin(w http.ResponseWriter, r *http.Request) {
	type Country struct {
		Flag, Name string
		Golfers    int
		Percent    float64
	}

	type Session struct {
		Golfer   string
		LastUsed time.Time
	}

	type Table struct {
		Name       sql.NullString
		Rows, Size int
	}

	data := struct {
		Countries []Country
		Sessions  []Session
		Tables    []Table
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
		`SELECT COALESCE(country, ''), COUNT(*), COUNT(*) / SUM(COUNT(*)) OVER () * 100
		   FROM users GROUP BY COALESCE(country, '')`,
	)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var c Country
		var id string

		if err := rows.Scan(&id, &c.Golfers, &c.Percent); err != nil {
			panic(err)
		}

		if country, ok := config.CountryByID[id]; ok {
			c.Flag = country.Flag
			c.Name = country.Name
		}

		data.Countries = append(data.Countries, c)
	}

	sort.Slice(data.Countries, func(i, j int) bool {
		if data.Countries[i].Golfers != data.Countries[j].Golfers {
			return data.Countries[i].Golfers > data.Countries[j].Golfers
		}

		return data.Countries[i].Name < data.Countries[j].Name
	})

	if err := rows.Err(); err != nil {
		panic(err)
	}

	render(w, r, "admin/info", data, "Admin Info")
}
