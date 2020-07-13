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

	data := struct {
		Sessions []Session
		Tables   []Table
	}{}

	db := session.Database(r)

	rows, err := db.Query(
		` SELECT login, MAX(last_used)
		    FROM sessions JOIN users ON user_id = users.id
		GROUP BY login
		ORDER BY max DESC`,
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

	render(w, r, "admin/info", "Admin Info", data)
}
