package routes

import (
	"net/http"
	"time"

	"github.com/code-golf/code-golf/session"
	"github.com/code-golf/code-golf/trophy"
)

// RankingsCheevos serves GET /rankings/cheevos/{cheevo}
func RankingsCheevos(w http.ResponseWriter, r *http.Request) {
	cheevoID := param(r, "cheevo")

	type row struct {
		Country, Login string
		Earned         time.Time
		Rank, Count    int
	}

	data := struct {
		Cheevo  *trophy.Trophy
		Rows    []row
		Cheevos map[string][]*trophy.Trophy
		Total   int
	}{
		Cheevo:  trophy.ByID[cheevoID],
		Cheevos: trophy.Tree,
		Total:   len(trophy.List),
	}

	if cheevoID != "" && data.Cheevo == nil {
		NotFound(w, r)
		return
	}

	rows, err := session.Database(r).Query(
		`WITH count AS (
		    SELECT user_id, COUNT(*), MAX(earned) earned
		      FROM trophies
		     WHERE $1 IN ('', trophy::text)
		  GROUP BY user_id
		) SELECT count,
		         COALESCE(CASE WHEN show_country THEN country END, ''),
		         earned,
		         login,
		         CASE WHEN $1 = ''
		            THEN RANK() OVER(ORDER BY count DESC)
		            ELSE RANK() OVER(ORDER BY earned)
		         END
		    FROM count JOIN users ON id = user_id
		ORDER BY rank, earned
		   LIMIT 30`,
		cheevoID,
	)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		var r row

		if err := rows.Scan(
			&r.Count,
			&r.Country,
			&r.Earned,
			&r.Login,
			&r.Rank,
		); err != nil {
			panic(err)
		}

		data.Rows = append(data.Rows, r)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	render(w, r, "rankings/cheevos", "Rankings: Achievements", data)
}
