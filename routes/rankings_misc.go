package routes

import (
	"html/template"
	"net/http"
	"time"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/golfer"
	"github.com/code-golf/code-golf/pager"
	"github.com/code-golf/code-golf/session"
)

// GET /rankings/misc/{type}
func rankingsMiscGET(w http.ResponseWriter, r *http.Request) {
	var desc template.HTML
	var sql string

	data := struct {
		Pager *pager.Pager
		Rows  []struct {
			golfer.GolferLink

			Bytes, Chars, Count, Rank, Total int
			Hole                             *config.Hole
			Lang                             *config.Lang
			Me                               bool
			Scoring                          string
			Submitted                        time.Time
		}
	}{Pager: pager.New(r)}

	args := []any{pager.PerPage, data.Pager.Offset}

	switch t := param(r, "type"); t {
	case "followers":
		desc = "Total followers."
		sql = `WITH followers AS (
			    SELECT followee_id, COUNT(*) FROM follows GROUP BY followee_id
			) SELECT avatar_url, count, country_flag, name,
			         RANK() OVER(ORDER BY count DESC),
			         COUNT(*) OVER () total
			    FROM followers
			    JOIN golfers_with_avatars ON id = followee_id
			ORDER BY rank, name
			   LIMIT $1 OFFSET $2`
	case "holes-authored":
		desc = "Total holes authored."
		sql = `WITH holes AS (
			    SELECT user_id, COUNT(*) FROM authors GROUP BY user_id
			) SELECT avatar_url, count, country_flag, name,
			         RANK() OVER(ORDER BY count DESC),
			         COUNT(*) OVER () total
			    FROM holes
			    JOIN golfers_with_avatars ON id = user_id
			ORDER BY rank, name
			   LIMIT $1 OFFSET $2`
	case "holes-of-the-week":
		desc, _ = config.HoleOfTheWeek()
		sql = `WITH holes AS (
			    SELECT user_id, COUNT(*), MAX(completed) submitted
			      FROM weekly_solves
			  GROUP BY user_id
			) SELECT avatar_url, count, country_flag, name, submitted,
			         RANK() OVER(ORDER BY count DESC, submitted),
			         COUNT(*) OVER () total
			    FROM holes
			    JOIN golfers_with_avatars ON id = user_id
			ORDER BY rank, name
			   LIMIT $1 OFFSET $2`
	case "referrals":
		desc = "Total referrals."
		sql = `WITH referrals AS (
			    SELECT referrer_id, COUNT(*) FROM users GROUP BY referrer_id
			) SELECT avatar_url, count, country_flag, name,
			         RANK() OVER(ORDER BY count DESC),
			         COUNT(*) OVER () total
			    FROM referrals
			    JOIN golfers_with_avatars ON id = referrals.referrer_id
			ORDER BY rank, name
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
			) SELECT avatar_url, bytes, chars, count, country_flag, name,
			         RANK() OVER(ORDER BY count DESC),
			         COUNT(*) OVER () total
			    FROM solutions
			    JOIN golfers_with_avatars ON id = user_id
			ORDER BY rank, bytes, chars, name
			   LIMIT $1 OFFSET $2`
	default:
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err := session.Database(r).Select(&data.Rows, sql, args...); err != nil {
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
