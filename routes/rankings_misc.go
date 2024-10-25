package routes

import (
	"net/http"
	"time"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/pager"
	"github.com/code-golf/code-golf/session"
)

// GET /rankings/misc/{type}
func rankingsMiscGET(w http.ResponseWriter, r *http.Request) {
	var desc, sql string

	data := struct {
		Pager *pager.Pager
		Rows  []struct {
			Bytes, Chars, Count, Rank, Total int
			Country                          config.NullCountry `db:"country_flag"`
			Hole, Lang, Scoring              string
			Me                               bool
			Name                             string `db:"login"`
			Submitted                        time.Time
		}
	}{Pager: pager.New(r)}

	args := []any{pager.PerPage, data.Pager.Offset}

	switch t := param(r, "type"); t {
	case "diamond-deltas":
		desc = "Deltas between diamonds and silvers."
		sql = `WITH diamonds AS (
			    SELECT * FROM rankings WHERE rank = 1 AND tie_count = 1
			), silvers AS (
			    SELECT DISTINCT hole, lang, scoring, strokes
			      FROM rankings
			     WHERE rank = 2
			) SELECT country_flag, hole, lang, login, scoring,
			         silvers.strokes - diamonds.strokes count,
			         RANK() OVER(ORDER BY silvers.strokes - diamonds.strokes DESC),
			         COUNT(*) OVER () total
			    FROM diamonds
			    JOIN silvers USING (hole, lang, scoring)
			    JOIN users ON id = user_id
			ORDER BY rank, scoring
			   LIMIT $1 OFFSET $2`
	case "followers":
		desc = "Total followers."
		sql = `WITH followers AS (
			    SELECT followee_id, COUNT(*) FROM follows GROUP BY followee_id
			) SELECT count, country_flag, login,
			         RANK() OVER(ORDER BY count DESC),
			         COUNT(*) OVER () total
			    FROM followers
			    JOIN users ON id = followee_id
			ORDER BY rank, login
			   LIMIT $1 OFFSET $2`
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
	case "most-tied-golds":
		args = append(args, session.Golfer(r))
		desc = "Most tied gold medals"
		sql = `SELECT hole, lang, scoring, COUNT(*),
			          RANK() OVER(ORDER BY COUNT(*) DESC),
			          COUNT(*) FILTER (WHERE user_id = $3) > 0 me,
			          COUNT(*) OVER () total
			     FROM medals
			    WHERE medal = 'gold'
			 GROUP BY hole, lang, scoring
			 ORDER BY rank, hole, lang, scoring
			    LIMIT $1 OFFSET $2`
	case "oldest-diamonds", "oldest-unicorns":
		if t == "oldest-diamonds" {
			args = append(args, "diamond")
			desc = "💎 Oldest diamonds (uncontested gold medals)."
		} else {
			args = append(args, "unicorn")
			desc = "🦄 Oldest unicorns (uncontested solves)."
		}

		sql = `SELECT country_flag, hole, lang, login, scoring, submitted,
			         RANK() OVER(ORDER BY submitted),
			         COUNT(*) OVER () total
			    FROM medals
			    JOIN users ON id = user_id
			   WHERE medal = $3
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
