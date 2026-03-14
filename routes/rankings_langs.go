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
			Golds, Silvers, Bronzes int
			Points, Rank, Strokes   int
			Submitted               time.Time
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
			           r.hole, r.lang, r.points, r.strokes, r.user_id, r.submitted,
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
			           points,
			           RANK() OVER (ORDER BY strokes) AS rank,
			           strokes,
			           submitted
			      FROM best_per_lang
			)
			SELECT *
			  FROM ranked_langs
			 WHERE $3 IN ('all', lang::text)
			ORDER BY rank, submitted, lang`,
			data.HoleID, data.Scoring, data.LangID,
		); err != nil {
			panic(err)
		}
	} else if data.LangID != "all" {
		// Get best scores per hole for this specific language
		if err := session.Database(r).Select(
			&data.Rows,
			`WITH best_per_lang AS (
			    SELECT DISTINCT ON (r.hole, r.lang)
			           r.hole, r.lang, r.points, r.strokes, r.user_id, r.submitted,
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
			SELECT *
			  FROM ranked_langs
			 WHERE lang = $2
			ORDER BY hole`,
			data.Scoring, data.LangID,
		); err != nil {
			panic(err)
		}
	} else {
		// Overall default medals view (All Holes & All Langs)
		if err := session.Database(r).Select(
			&data.Rows,
			`WITH ranks AS (
			    SELECT hole, lang, points, strokes,
			           RANK()       OVER (PARTITION BY hole
			                                  ORDER BY points DESC, strokes),
			           ROW_NUMBER() OVER (PARTITION BY hole, lang
			                                  ORDER BY points DESC, strokes)
			      FROM rankings
			     WHERE scoring = $1 AND NOT experimental
			), medals AS (
			    SELECT DISTINCT hole, lang, rank FROM ranks WHERE rank < 4
			) SELECT lang,
			         SUM(points)  points,
			         SUM(strokes) strokes,
			         RANK() OVER(ORDER BY SUM(points) DESC, SUM(strokes) DESC),
			         (SELECT COUNT(*) FROM medals WHERE lang = ranks.lang AND rank = 1) golds,
			         (SELECT COUNT(*) FROM medals WHERE lang = ranks.lang AND rank = 2) silvers,
			         (SELECT COUNT(*) FROM medals WHERE lang = ranks.lang AND rank = 3) bronzes
			    FROM ranks
			   WHERE row_number = 1
			GROUP BY lang
			ORDER BY rank, lang`,
			data.Scoring,
		); err != nil {
			panic(err)
		}
	}

	var description string
	if data.HoleID != "all" && data.LangID != "all" {
		description = "Best " + data.Scoring + " score in " + data.Hole.Name + " in " + config.LangByID[data.LangID].Name + "."
	} else if data.HoleID != "all" {
		description = "Best " + data.Scoring + " scores per language in " + data.Hole.Name + "."
	} else if data.LangID != "all" {
		description = "Best " + data.Scoring + " scores per hole in " + config.LangByID[data.LangID].Name + "."
	} else {
		description = "All languages in " + data.Scoring + "."
	}

	render(w, r, "rankings/langs", data, "Rankings: Languages", description)
}
