package routes

import (
	"net/http"
	"time"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/golfer"
	"github.com/code-golf/code-golf/pager"
	"github.com/code-golf/code-golf/session"
)

// GET /rankings/history/{hole}/{lang}/{scoring}/{type}
func rankingsHistoryGET(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Hole, PrevHole, NextHole *config.Hole
		HoleID, LangID, Scoring  string
		Type                     string
		Pager                    *pager.Pager
		Rows                     []struct {
			golfer.GolferLink
			Hole      *config.Hole
			Lang      *config.Lang
			Scoring   string
			Count     int
			Rank      int
			Total     int
			Me        bool
			Submitted time.Time
		}
	}{
		HoleID:  param(r, "hole"),
		LangID:  param(r, "lang"),
		Pager:   pager.New(r),
		Scoring: param(r, "scoring"),
		Type:    param(r, "type"),
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

	var description string
	var sql string

	args := []any{data.HoleID, data.LangID, data.Scoring, pager.PerPage, data.Pager.Offset}

	switch data.Type {
	case "diamond-deltas":
		description = "Deltas between diamonds and silvers"
		sql = `WITH diamonds AS (
			    SELECT *
			      FROM rankings
			     WHERE rank = 1 AND tie_count = 1 AND NOT experimental
			), silvers AS (
			    SELECT DISTINCT hole, lang, scoring, strokes
			      FROM rankings
			     WHERE rank = 2 AND NOT experimental
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
		var userID int
		if golfer := session.Golfer(r); golfer != nil {
			userID = golfer.ID
		}
		args = append(args, userID) // $6

		description = "Most tied gold medals"
		sql = `SELECT hole, lang, scoring, COUNT(*) count,
			          RANK() OVER(ORDER BY COUNT(*) DESC),
			          COUNT(*) FILTER (WHERE user_id = $6) > 0 me,
			          COUNT(*) OVER () total
			     FROM medals
			    WHERE medal = 'gold'
			      AND $1 IN ('all', hole::text)
			      AND $2 IN ('all', lang::text)
			      AND $3 IN ('all', scoring::text)
			 GROUP BY hole, lang, scoring
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
		args = append(args, medal) // $6

		sql = `SELECT avatar_url, country_flag, name,
			         hole, lang, scoring, submitted,
			         RANK() OVER(ORDER BY submitted),
			         COUNT(*) OVER () total
			    FROM medals
			    JOIN golfers_with_avatars ON id = user_id
			   WHERE medal = $6
			     AND $1 IN ('all', hole::text)
			     AND $2 IN ('all', lang::text)
			     AND $3 IN ('all', scoring::text)
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

	render(w, r, "rankings/history", data, "Rankings: History", description)
}