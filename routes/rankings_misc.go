package routes

import (
	"net/http"
	"time"

	"github.com/code-golf/code-golf/pager"
	"github.com/code-golf/code-golf/session"
)

// GET /rankings/misc/{type}
func rankingsMiscGET(w http.ResponseWriter, r *http.Request) {
	type row struct {
		Bytes, Chars, Rank, Total           int
		Country, Hole, Lang, Login, Scoring string
		Submitted                           *time.Time
	}

	data := struct {
		Pager *pager.Pager
		Rows  []row
	}{pager.New(r), make([]row, 0, pager.PerPage)}

	var desc, sql string

	switch param(r, "type") {
	case "holes-authored":
		desc = "Total holes authored."
		sql =
			`WITH holes AS (
			    SELECT user_id, COUNT(*) FROM authors GROUP BY user_id
			) SELECT 0, 0,
			         country_flag,
			         '', '',
			         login,
			         RANK() OVER(ORDER BY count DESC),
			         '', NULL,
			         count,
			         COUNT(*) OVER ()
			    FROM holes
			    JOIN users ON id = user_id
			ORDER BY rank, login
			   LIMIT $1 OFFSET $2`
	case "oldest-diamonds":
		desc = "Oldest diamond medals."
		sql =
			`WITH diamonds AS (
			    SELECT hole, lang, scoring, submitted, user_id
			      FROM rankings
			     WHERE rank = 1 AND tie_count = 1
			) SELECT 0, 0,
			         country_flag,
			         hole,
			         lang,
			         login,
			         RANK() OVER(ORDER BY submitted),
			         scoring,
			         submitted,
			         0,
			         COUNT(*) OVER ()
			    FROM diamonds
			    JOIN users ON id = user_id
			ORDER BY rank, hole, lang, scoring, login
			   LIMIT $1 OFFSET $2`
	case "referrals":
		desc = "Total referrals."
		sql =
			`WITH referrals AS (
			    SELECT referrer_id, COUNT(*) FROM users GROUP BY referrer_id
			) SELECT 0, 0,
			         country_flag,
			         '', '',
			         login,
			         RANK() OVER(ORDER BY count DESC),
			         '', NULL,
			         count,
			         COUNT(*) OVER ()
			    FROM referrals
			    JOIN users ON id = referrals.referrer_id
			ORDER BY rank, login
			   LIMIT $1 OFFSET $2`
	case "solutions":
		desc = "Total solutions."
		sql =
			`WITH solutions AS (
			    SELECT user_id,
			           COUNT(*),
			           SUM(bytes)              bytes,
			           COALESCE(SUM(chars), 0) chars
			      FROM solutions
			     WHERE NOT failing
			  GROUP BY user_id
			) SELECT bytes,
			         chars,
			         country_flag,
			         '', '',
			         login,
			         RANK() OVER(ORDER BY count DESC),
			         '', NULL,
			         count,
			         COUNT(*) OVER ()
			    FROM solutions
			    JOIN users ON id = user_id
			ORDER BY rank, bytes, chars, login
			   LIMIT $1 OFFSET $2`
	default:
		w.WriteHeader(http.StatusNotFound)
		return
	}

	rows, err := session.Database(r).Query(sql, pager.PerPage, data.Pager.Offset)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var r row

		if err := rows.Scan(
			&r.Bytes,
			&r.Chars,
			&r.Country,
			&r.Hole,
			&r.Lang,
			&r.Login,
			&r.Rank,
			&r.Scoring,
			&r.Submitted,
			&r.Total,
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
		w.WriteHeader(http.StatusNotFound)
		return
	}

	render(w, r, "rankings/misc", data, "Rankings: Miscellaneous", desc)
}
