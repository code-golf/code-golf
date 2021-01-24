package routes

import (
	"net/http"

	"github.com/code-golf/code-golf/session"
)

// RankingsSolutions serves GET /rankings/solutions/{page}
func RankingsSolutions(w http.ResponseWriter, r *http.Request) {
	type row struct {
		BytesPer, CharsPer, Country, Login string
		Bytes, Chars, Rank, Count          int
	}

	var data []row

	rows, err := session.Database(r).Query(
		`WITH solutions AS (
		    SELECT user_id,
		           COUNT(*),
		           SUM(bytes) bytes,
		           SUM(chars) chars
		      FROM solutions
		      JOIN code ON id = code_id
		     WHERE NOT failing
		  GROUP BY user_id
		) SELECT bytes,
		         TO_CHAR(bytes::decimal / count, 'FM999,999.0'),
		         chars,
		         TO_CHAR(chars::decimal / count, 'FM999,999.0'),
		         count,
		         COALESCE(CASE WHEN show_country THEN country END, ''),
		         login,
		         RANK() OVER(ORDER BY count DESC)
		    FROM solutions
		    JOIN users on id = user_id
		ORDER BY rank, bytes, chars, login
		   LIMIT 25`,
	)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var r row

		if err := rows.Scan(
			&r.Bytes,
			&r.BytesPer,
			&r.Chars,
			&r.CharsPer,
			&r.Count,
			&r.Country,
			&r.Login,
			&r.Rank,
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
