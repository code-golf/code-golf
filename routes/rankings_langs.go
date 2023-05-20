package routes

import (
	"net/http"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/session"
)

// GET /rankings/langs/{scoring}
func rankingsLangsGET(w http.ResponseWriter, r *http.Request) {
	type row struct {
		Hole                    *config.Hole
		Lang                    *config.Lang
		Golds, Silvers, Bronzes int
		Points, Rank, Strokes   int
	}

	data := struct {
		LangID, Scoring string
		Langs           []*config.Lang
		Rows            []row
	}{
		LangID:  param(r, "lang"),
		Langs:   config.LangList,
		Scoring: param(r, "scoring"),
	}

	if data.LangID != "all" && config.LangByID[data.LangID] == nil ||
		data.Scoring != "chars" && data.Scoring != "bytes" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if data.LangID != "all" {
		if err := session.Database(r).Select(
			&data.Rows,
			`WITH ranks AS (
			    SELECT hole, lang,
			           RANK() OVER (PARTITION BY hole
			                            ORDER BY points DESC, strokes)
			      FROM rankings
			     WHERE scoring = $1
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
			     WHERE scoring = $1
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
	if lang, ok := config.LangByID[data.LangID]; ok {
		description = lang.Name + " in "
	} else {
		description = "All languages in "
	}
	description += data.Scoring + "."

	render(w, r, "rankings/langs", data, "Rankings: Languages", description)
}
