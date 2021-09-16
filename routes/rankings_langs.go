package routes

import (
	"net/http"

	"github.com/code-golf/code-golf/hole"
	"github.com/code-golf/code-golf/lang"
	"github.com/code-golf/code-golf/session"
)

// RankingsLangs serves GET /rankings/langs/{scoring}
func RankingsLangs(w http.ResponseWriter, r *http.Request) {
	type row struct {
		Hole                    hole.Hole
		Lang                    lang.Lang
		Golds, Silvers, Bronzes int
		Points, Rank, Strokes   int
	}

	data := struct {
		LangID, Scoring string
		Langs           []lang.Lang
		Rows            []row
	}{
		LangID:  param(r, "lang"),
		Langs:   lang.List,
		Scoring: param(r, "scoring"),
	}

	if data.LangID != "all" && lang.ByID[data.LangID].ID == "" ||
		data.Scoring != "chars" && data.Scoring != "bytes" {
		NotFound(w, r)
		return
	}

	if data.LangID != "all" {
		rows, err := session.Database(r).Query(
			`WITH ranks AS (
			    SELECT hole, lang,
			           RANK() OVER (PARTITION BY hole ORDER BY points DESC)
			      FROM rankings
			     WHERE scoring = $1
			) SELECT DISTINCT hole, rank
			    FROM ranks
			   WHERE lang = $2 AND rank < 4
			ORDER BY hole, rank`,
			data.Scoring,
			data.LangID,
		)
		if err != nil {
			panic(err)
		}

		for rows.Next() {
			var r row
			var holeID string

			if err := rows.Scan(&holeID, &r.Rank); err != nil {
				panic(err)
			}

			r.Hole = hole.ByID[holeID]

			data.Rows = append(data.Rows, r)
		}

		if err := rows.Err(); err != nil {
			panic(err)
		}
	} else {
		rows, err := session.Database(r).Query(
			`WITH ranks AS (
			    SELECT hole, lang, points, strokes,
			           RANK()       OVER (PARTITION BY hole       ORDER BY points DESC),
			           ROW_NUMBER() OVER (PARTITION BY hole, lang ORDER BY points DESC)
			      FROM rankings
			     WHERE scoring = $1
			), medals AS (
			    SELECT DISTINCT hole, lang, rank FROM ranks WHERE rank < 4
			) SELECT lang,
			         SUM(points),
			         SUM(strokes),
			         RANK() OVER(ORDER BY SUM(points) DESC, SUM(strokes) DESC),
			         (SELECT COUNT(*) FROM medals WHERE lang = ranks.lang AND rank = 1),
			         (SELECT COUNT(*) FROM medals WHERE lang = ranks.lang AND rank = 2),
			         (SELECT COUNT(*) FROM medals WHERE lang = ranks.lang AND rank = 3)
			    FROM ranks
			   WHERE row_number = 1
			GROUP BY lang
			ORDER BY rank, lang`,
			data.Scoring,
		)
		if err != nil {
			panic(err)
		}

		for rows.Next() {
			var r row
			var langID string

			if err := rows.Scan(
				&langID,
				&r.Points,
				&r.Strokes,
				&r.Rank,
				&r.Golds,
				&r.Silvers,
				&r.Bronzes,
			); err != nil {
				panic(err)
			}

			r.Lang = lang.ByID[langID]

			data.Rows = append(data.Rows, r)
		}

		if err := rows.Err(); err != nil {
			panic(err)
		}
	}

	var description string
	if lang := lang.ByID[data.LangID]; lang.ID != "" {
		description = lang.Name + " in "
	} else {
		description = "All languages in "
	}
	description += data.Scoring + "."

	render(w, r, "rankings/langs", data, "Rankings: Languages", description)
}
