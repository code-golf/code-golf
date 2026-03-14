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
		HoleID, LangID, Scoring string
		Hole                    *config.Hole
		Langs                   []*config.Lang
		Rows                    []struct {
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

	if data.HoleID != "all" {
		data.Hole = config.HoleByID[data.HoleID]

		var langFilter string
		var args []any
		if data.LangID != "all" {
			langFilter = "AND lang = $3"
			args = []any{data.HoleID, data.Scoring, data.LangID}
		} else {
			args = []any{data.HoleID, data.Scoring}
		}

		if err := session.Database(r).Select(
			&data.Rows,
			`WITH best AS (
			    SELECT lang, points, strokes, user_id, submitted,
			           ROW_NUMBER() OVER (PARTITION BY lang
			                                  ORDER BY strokes, submitted) AS rn
			      FROM rankings
			     WHERE hole = $1 AND scoring = $2 `+langFilter+`
			) SELECT avatar_url,
			         country_flag,
			         lang,
			         name,
			         points,
			         RANK() OVER (ORDER BY strokes) AS rank,
			         strokes,
			         submitted
			    FROM best
			    JOIN golfers_with_avatars ON user_id = id
			   WHERE rn = 1
			ORDER BY rank, lang`,
			args...,
		); err != nil {
			panic(err)
		}
	} else if data.LangID != "all" {
		if err := session.Database(r).Select(
			&data.Rows,
			`WITH ranks AS (
			    SELECT hole, lang,
			           RANK() OVER (PARTITION BY hole
			                            ORDER BY points DESC, strokes)
			      FROM rankings
			     WHERE scoring = $1 AND NOT experimental
			) SELECT DISTINCT hole, rank
			    FROM ranks
			   WHERE lang = $2 AND rank < 4
			ORDER BY hole, rank`,
			data.Scoring,
			data.LangID,
		); err != nil {
			panic(err)
		}
	} else {
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
	if data.HoleID != "all" {
		description = "Best " + data.Scoring + " scores per language in " + data.Hole.Name + "."
	} else {
		if lang, ok := config.LangByID[data.LangID]; ok {
			description = lang.Name + " in "
		} else {
			description = "All languages in "
		}
		description += data.Scoring + "."
	}

	render(w, r, "rankings/langs", data, "Rankings: Languages", description)
}
