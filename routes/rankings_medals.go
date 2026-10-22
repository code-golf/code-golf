package routes

import (
	"net/http"
	"time"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/golfer"
	"github.com/code-golf/code-golf/pager"
	"github.com/code-golf/code-golf/session"
)

// GET /rankings/medals/{type}/{hole}/{lang}/{scoring}
func rankingsMedalsGET(w http.ResponseWriter, r *http.Request) {
	data := struct {
		AI                       string
		Hole, PrevHole, NextHole *config.Hole
		HoleID, LangID, Scoring  string
		Type                     string
		Pager                    *pager.Pager
		Rows                     []struct {
			golfer.GolferLink
			Hole *config.Hole
			Lang *config.Lang

			Unicorn, Diamond, Gold, Silver, Bronze int
			Rank, Total                            int

			Count     int
			Me        bool
			Scoring   string
			Submitted time.Time
		}
	}{
		AI:      r.URL.Query().Get("ai"),
		HoleID:  param(r, "hole"),
		LangID:  param(r, "lang"),
		Pager:   pager.New(r),
		Scoring: param(r, "scoring"),
		Type:    param(r, "type"),
	}

	// aiFilter is nil for "Both", true for "AI Only", false for "Human Only".
	var aiFilter *bool

	switch data.AI {
	case "":
		data.AI = "all"
	case "human":
		aiFilter = new(false)
	case "ai":
		aiFilter = new(true)
	case "all":
	default:
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if data.HoleID != "all" && config.HoleByID[data.HoleID] == nil ||
		data.LangID != "all" && config.LangByID[data.LangID] == nil ||
		data.Scoring != "all" && data.Scoring != "chars" && data.Scoring != "bytes" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if data.Hole = config.HoleByID[data.HoleID]; data.Hole != nil {
		data.PrevHole, data.NextHole = getPrevNextHole(r, data.Hole)
	}

	var description, sql string
	args := []any{data.HoleID, data.LangID, data.Scoring, pager.PerPage, data.Pager.Offset, aiFilter}

	// 1. Determine SQL and base Description based on Type
	switch data.Type {
	case "overall":
		description = "Overall medal counts"
		sql = `WITH counts AS (
			    SELECT user_id,
			           COUNT(*) FILTER (WHERE medal = 'unicorn') unicorn,
			           COUNT(*) FILTER (WHERE medal = 'diamond') diamond,
			           COUNT(*) FILTER (WHERE medal = 'gold'   ) gold,
			           COUNT(*) FILTER (WHERE medal = 'silver' ) silver,
			           COUNT(*) FILTER (WHERE medal = 'bronze' ) bronze
			      FROM medals
			     WHERE $1 IN ('all', hole::text)
			       AND $2 IN ('all', lang::text)
			       AND $3 IN ('all', scoring::text)
			       AND uses_ai IS NOT DISTINCT FROM $6
			  GROUP BY user_id
			) SELECT RANK() OVER(
			             ORDER BY gold DESC, diamond DESC, silver DESC, bronze DESC
			         ),
			         unicorn, diamond, gold, silver, bronze,
			         avatar_url, country_flag, name, COUNT(*) OVER() total
			    FROM counts
			    JOIN golfers_with_avatars ON id = user_id
			ORDER BY rank, name
			   LIMIT $4 OFFSET $5`

	case "diamond-deltas":
		description = "Deltas between diamonds and silvers"
		sql = `WITH diamonds AS (
			    SELECT *
			      FROM medals
			     WHERE medal = 'diamond'
			       AND uses_ai IS NOT DISTINCT FROM $6
			), silvers AS (
			    SELECT DISTINCT hole, lang, scoring, strokes
			      FROM medals
			     WHERE medal = 'silver'
			       AND uses_ai IS NOT DISTINCT FROM $6
			) SELECT avatar_url, country_flag, name,
			         hole, lang, scoring,
			         silvers.strokes - diamonds.strokes count,
			         RANK() OVER(ORDER BY silvers.strokes - diamonds.strokes DESC),
			         COUNT(*) OVER () total
			    FROM diamonds
			    JOIN silvers USING (hole, lang, scoring)
			    JOIN golfers_with_avatars ON id = user_id
			   WHERE $1 IN ('all', hole::text)
			     AND $2 IN ('all', lang::text)
			     AND $3 IN ('all', scoring::text)
			ORDER BY rank, scoring
			   LIMIT $4 OFFSET $5`

	case "most-tied-golds":
		description = "Most tied gold medals"
		args = append(args, session.Golfer(r)) // $7

		sql = `SELECT hole, lang, scoring, COUNT(*) count,
			          RANK() OVER(ORDER BY COUNT(*) DESC),
			          COUNT(*) FILTER (WHERE user_id = $7) > 0 me,
			          COUNT(*) OVER () total
			     FROM medals
			    WHERE medal = 'gold'
			      AND $1 IN ('all', hole::text)
			      AND $2 IN ('all', lang::text)
			      AND $3 IN ('all', scoring::text)
			      AND uses_ai IS NOT DISTINCT FROM $6
			 GROUP BY hole, lang, scoring
				HAVING COUNT(*) > 1
			 ORDER BY rank, hole, lang, scoring
			    LIMIT $4 OFFSET $5`

	case "oldest-diamonds", "oldest-unicorns":
		medal := "unicorn"
		if data.Type == "oldest-diamonds" {
			medal = "diamond"
			description = "💎 Oldest diamonds (uncontested gold medals)"
		} else {
			description = "🦄 Oldest unicorns (uncontested solves)"
		}
		args = append(args, medal) // $7

		sql = `SELECT avatar_url, country_flag, name,
			          hole, lang, scoring, submitted,
			          RANK() OVER(ORDER BY submitted),
			          COUNT(*) OVER () total
			     FROM medals
			     JOIN golfers_with_avatars ON id = user_id
			    WHERE medal = $7
			      AND $1 IN ('all', hole::text)
			      AND $2 IN ('all', lang::text)
			      AND $3 IN ('all', scoring::text)
			      AND uses_ai IS NOT DISTINCT FROM $6
			 ORDER BY rank, hole, lang, scoring, name
			    LIMIT $4 OFFSET $5`

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

	if data.Hole != nil {
		description += " in " + data.Hole.Name
	}

	if lang, ok := config.LangByID[data.LangID]; ok {
		description += " in " + lang.Name
	}

	if data.Scoring != "all" {
		description += " in " + data.Scoring
	}

	if data.Type == "overall" {
		description += ". Ranked by golds, then diamonds, then silvers, then bronzes"
	}

	render(w, r, "rankings/medals", data, "Rankings: Medals", description+".")
}
