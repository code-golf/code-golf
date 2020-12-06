package routes

import (
	"net/http"
	"time"

	"github.com/code-golf/code-golf/hole"
	"github.com/code-golf/code-golf/lang"
	"github.com/code-golf/code-golf/session"
)

// RankingsHoles serves GET /rankings/holes/{hole}/{lang}/{scoring}/{page}
func RankingsHoles(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Holes []hole.Hole
		Langs []lang.Lang
	}{
		Holes: hole.List,
		Langs: lang.List,
	}

	render(w, r, "rankings/holes", "Rankings: Holes", data)
}

func RankingsSolutions(w http.ResponseWriter, r *http.Request) {
	type row struct {
		Country, Login            string
		Bytes, Chars, Rank, Count int
		Submitted                 time.Time
	}

	var data []row

	rows, err := session.Database(r).Query(
		`WITH solutions AS (
		    SELECT user_id,
		           COUNT(*),
		           SUM(bytes)     bytes,
		           SUM(chars)     chars,
		           MIN(submitted) submitted
		      FROM solutions
		      JOIN code ON id = code_id
		     WHERE NOT failing
		  GROUP BY user_id
		) SELECT bytes,
		         chars,
		         count,
		         COALESCE(CASE WHEN show_country THEN country END, ''),
		         login,
		         RANK() OVER(ORDER BY count DESC),
		         submitted
		    FROM solutions
		    JOIN users on id = user_id
		ORDER BY rank, bytes, chars, submitted
		   LIMIT 25`,
	)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var r row

		if err := rows.Scan(
			&r.Bytes,
			&r.Chars,
			&r.Count,
			&r.Country,
			&r.Login,
			&r.Rank,
			&r.Submitted,
		); err != nil {
			panic(err)
		}

		data = append(data, r)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	render(w, r, "rankings/solutions", "Rankings: Solutions", data)
}
