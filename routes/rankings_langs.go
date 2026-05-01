package routes

import (
	"net/http"
	"time"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/golfer"
	"github.com/code-golf/code-golf/session"
)

// GET /rankings/langs/{hole}/{lang}/{scoring}
func rankingsLangsGET(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Hole, PrevHole, NextHole *config.Hole
		HoleID, LangID, Scoring  string
		Langs                    []*config.Lang
		Rows                     []struct {
			golfer.GolferLink
			Hole                    *config.Hole
			Lang                    *config.Lang
			Me                      bool
			Golds, Silvers, Bronzes int
			Points, Rank, Strokes   int
			Submitted               time.Time
			Ties                    int
		}
	}{
		HoleID:  param(r, "hole"),
		LangID:  param(r, "lang"),
		Langs:   config.LangList,
		Scoring: param(r, "scoring"),
	}

	if data.HoleID != "all" && config.HoleByID[data.HoleID] == nil ||
		data.LangID != "all" && config.LangByID[data.LangID] == nil ||
		data.Scoring != "chars" && data.Scoring != "bytes" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if data.Hole = config.HoleByID[data.HoleID]; data.Hole != nil {
		data.PrevHole, data.NextHole = getPrevNextHole(r, data.Hole)
	}

	if data.HoleID != "all" {
		data.Hole = config.HoleByID[data.HoleID]

		// Get best scores per language for this specific hole
		if err := session.Database(r).Select(
			&data.Rows,
			`WITH best_per_lang AS (
			    SELECT DISTINCT ON (r.lang)
			           r.hole, r.lang, r.strokes, r.points, r.user_id, r.submitted,
			           g.avatar_url, g.country_flag, g.name
			      FROM rankings r
			      JOIN golfers_with_avatars g ON r.user_id = g.id
			     WHERE r.hole = $1 AND r.scoring = $2 AND NOT r.experimental
			  ORDER BY r.lang, r.strokes, r.submitted
			), ranked_langs AS (
			    SELECT avatar_url,
			           country_flag,
			           hole,
			           lang,
			           name,
			           CASE WHEN $3 = 'all' THEN CAST(ROUND(MIN(strokes) OVER () * 1000.0 / strokes) AS int) ELSE points END AS points,
			           RANK() OVER (ORDER BY strokes) AS rank,
			           strokes,
			           submitted
			      FROM best_per_lang
			)
			SELECT r.*, m.golds - 1 ties, m.me
			  FROM ranked_langs r
			  JOIN (
			    SELECT lang, COUNT(*) AS golds,
			           COALESCE(bool_or(user_id = $4), false) me
			      FROM medals
			     WHERE hole = $1 AND scoring = $2 AND medal = 'gold'
			  GROUP BY lang
			 ) m ON r.lang = m.lang
			 WHERE $3 IN ('all', r.lang::text)
			ORDER BY r.rank, r.submitted, r.lang`,
			data.HoleID, data.Scoring, data.LangID, session.Golfer(r),
		); err != nil {
			panic(err)
		}
	} else if data.LangID != "all" {
		// Get best scores per hole for this specific language
		if err := session.Database(r).Select(
			&data.Rows,
			`WITH best_per_lang AS (
			    SELECT DISTINCT ON (r.hole, r.lang)
			           r.hole, r.lang, r.strokes, r.points, r.user_id, r.submitted,
			           g.avatar_url, g.country_flag, g.name
			      FROM rankings r
			      JOIN golfers_with_avatars g ON r.user_id = g.id
			     WHERE r.scoring = $1 AND NOT r.experimental
			  ORDER BY r.hole, r.lang, r.strokes, r.submitted
			), ranked_langs AS (
			    SELECT avatar_url,
			           country_flag,
			           hole,
			           lang,
			           name,
			           points,
			           RANK() OVER (PARTITION BY hole ORDER BY strokes) AS rank,
			           strokes,
			           submitted
			      FROM best_per_lang
			)
			SELECT r.*, m.golds - 1 ties, m.me
			  FROM ranked_langs r
			  JOIN (
			    SELECT hole, COUNT(*) AS golds,
			           COALESCE(bool_or(user_id = $3), false) me
			      FROM medals
			     WHERE lang = $2 AND scoring = $1 AND medal = 'gold'
			  GROUP BY hole
			 ) m ON r.hole = m.hole
			 WHERE r.lang = $2
			ORDER BY r.hole`,
			data.Scoring, data.LangID, session.Golfer(r),
		); err != nil {
			panic(err)
		}
	} else {
		// Overall default medals view (All Holes & All Langs)
		if err := session.Database(r).Select(
			&data.Rows,
			`WITH hole_stats AS (
			    SELECT hole, lang, MIN(strokes) AS strokes
			      FROM rankings
			     WHERE scoring = $1 AND NOT experimental
			     GROUP BY hole, lang
			), ratio_scores AS (
			    SELECT lang, strokes,
			           CAST(ROUND(MIN(strokes) OVER (PARTITION BY hole) * 1000.0 / strokes) AS int) AS points,
					   RANK() OVER (PARTITION BY hole ORDER BY strokes ASC) AS score_rank
			      FROM hole_stats
			), lang_scores AS (
			    SELECT lang,
			           SUM(points)  AS points,
			           SUM(strokes) AS strokes,
			           RANK() OVER(ORDER BY SUM(points) DESC, SUM(strokes) DESC) AS rank
			      FROM ratio_scores
			     GROUP BY lang
			), lang_medals AS (
			    SELECT lang,
					   COUNT(*) FILTER (WHERE score_rank = 1) AS golds,
					   COUNT(*) FILTER (WHERE score_rank = 2) AS silvers,
					   COUNT(*) FILTER (WHERE score_rank = 3) AS bronzes
				  FROM ratio_scores
				 WHERE score_rank <= 3
			     GROUP BY lang
			)
			SELECT s.lang,
			       s.points,
			       s.strokes,
			       s.rank,
			       COALESCE(m.golds, 0)   AS golds,
			       COALESCE(m.silvers, 0) AS silvers,
			       COALESCE(m.bronzes, 0) AS bronzes
			  FROM lang_scores s
			  LEFT JOIN lang_medals m ON s.lang = m.lang
			 ORDER BY s.rank, s.lang`,
			data.Scoring,
		); err != nil {
			panic(err)
		}
	}

	var description string
	if data.HoleID != "all" && data.LangID != "all" {
		description = "Best " + data.Scoring + " score in " + data.Hole.Name + " in " + config.LangByID[data.LangID].Name + "."
	} else if data.HoleID != "all" {
		description = "Best " + data.Scoring + " scores per language in " + data.Hole.Name + ". Points are the ratio to the overall minimum " + data.Scoring + "."
	} else if data.LangID != "all" {
		description = "Best " + data.Scoring + " scores per hole in " + config.LangByID[data.LangID].Name + "."
	} else {
		description = "All languages in " + data.Scoring + ". Points per hole are the ratio to the overall minimum " + data.Scoring + "."
	}

	render(w, r, "rankings/langs", data, "Rankings: Languages", description)
}
