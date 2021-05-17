package routes

import (
	"net/http"

	"github.com/code-golf/code-golf/pager"
	"github.com/code-golf/code-golf/session"
)

// RankingsSolutions serves GET /rankings/solutions
func RankingsSolutions(w http.ResponseWriter, r *http.Request) {
	type row struct {
		BytesPer, CharsPer, Country, Login string
		Bytes, Chars, Rank, Count, Langs   int
	}

	data := struct {
		Pager *pager.Pager
		Rows  []row
	}{pager.New(r), make([]row, 0, pager.PerPage)}

	rows, err := session.Database(r).Query(
		`WITH solutions AS (
		    SELECT user_id,
		           COUNT(*),
		           COUNT(DISTINCT lang) langs,
		           SUM(bytes) bytes,
		           SUM(chars) chars
		      FROM solutions
		     WHERE NOT failing
		  GROUP BY user_id
		) SELECT bytes,
		         TO_CHAR(bytes::decimal / count, 'FM999,999.0'),
		         chars,
		         TO_CHAR(chars::decimal / count, 'FM999,999.0'),
		         count,
		         COALESCE(CASE WHEN show_country THEN country END, ''),
		         langs,
		         login,
		         RANK() OVER(ORDER BY count DESC),
		         COUNT(*) OVER ()
		    FROM solutions
		    JOIN users on id = user_id
		ORDER BY rank, bytes, chars, login
		   LIMIT $1 OFFSET $2`,
		pager.PerPage,
		data.Pager.Offset,
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
			&r.Langs,
			&r.Login,
			&r.Rank,
			&data.Pager.Total,
		); err != nil {
			panic(err)
		}

		data.Rows = append(data.Rows, r)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	if data.Pager.Calculate() {
		NotFound(w, r)
		return
	}

	render(w, r, "rankings/solutions", data, "Rankings: Solutions", "All solutions")
}
