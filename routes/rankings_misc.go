package routes

import (
	"net/http"
	"time"

	"github.com/code-golf/code-golf/pager"
	"github.com/code-golf/code-golf/session"
)

// GET /rankings/misc/{type}
func rankingsMiscGET(w http.ResponseWriter, r *http.Request) {
	var desc, sql string

	switch param(r, "type") {
	case "holes-authored":
		desc = "Total holes authored."
		sql = `WITH holes AS (
			    SELECT user_id, COUNT(*) FROM authors GROUP BY user_id
			) SELECT count, country_flag, login,
			         RANK() OVER(ORDER BY count DESC),
			         COUNT(*) OVER () total
			    FROM holes
			    JOIN users ON id = user_id
			ORDER BY rank, login
			   LIMIT $1 OFFSET $2`
	case "oldest-diamonds":
		desc = "Oldest diamond medals."
		sql = `WITH diamonds AS (
			    SELECT hole, lang, scoring, submitted, user_id
			      FROM rankings
			     WHERE rank = 1 AND tie_count = 1
			) SELECT country_flag, hole, lang, login, scoring, submitted,
			         RANK() OVER(ORDER BY submitted),
			         COUNT(*) OVER () total
			    FROM diamonds
			    JOIN users ON id = user_id
			ORDER BY rank, hole, lang, scoring, login
			   LIMIT $1 OFFSET $2`
	case "referrals":
		desc = "Total referrals."
		sql = `WITH referrals AS (
			    SELECT referrer_id, COUNT(*) FROM users GROUP BY referrer_id
			) SELECT count, country_flag, login,
			         RANK() OVER(ORDER BY count DESC),
			         COUNT(*) OVER () total
			    FROM referrals
			    JOIN users ON id = referrals.referrer_id
			ORDER BY rank, login
			   LIMIT $1 OFFSET $2`
	case "solutions":
		desc = "Total solutions."
		sql = `WITH solutions AS (
			    SELECT user_id,
			           COUNT(*),
			           SUM(bytes)              bytes,
			           COALESCE(SUM(chars), 0) chars
			      FROM solutions
			     WHERE NOT failing
			  GROUP BY user_id
			) SELECT bytes, chars, count, country_flag, login,
			         RANK() OVER(ORDER BY count DESC),
			         COUNT(*) OVER () total
			    FROM solutions
			    JOIN users ON id = user_id
			ORDER BY rank, bytes, chars, login
			   LIMIT $1 OFFSET $2`
	default:
		w.WriteHeader(http.StatusNotFound)
		return
	}

	data := struct {
		Pager *pager.Pager
		Rows  []struct {
			Bytes, Chars, Count, Rank, Total        int
			CountryFlag, Hole, Lang, Login, Scoring string
			Submitted                               time.Time
		}
	}{Pager: pager.New(r)}

	if err := session.Database(r).Select(
		&data.Rows, sql, pager.PerPage, data.Pager.Offset,
	); err != nil {
		panic(err)
	}

	if len(data.Rows) > 0 {
		data.Pager.Total = data.Rows[0].Total
	}

	if data.Pager.Calculate() {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	render(w, r, "rankings/misc", data, "Rankings: Miscellaneous", desc)
}
