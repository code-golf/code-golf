package routes

import (
	"net/http"

	"github.com/code-golf/code-golf/lang"
	"github.com/code-golf/code-golf/session"
)

// RankingsLangs serves GET /rankings/langs/{scoring}
func RankingsLangs(w http.ResponseWriter, r *http.Request) {
	type row struct {
		Lang                    lang.Lang
		Golds, Silvers, Bronzes int
		Points, Rank, Strokes   int
	}

	data := struct {
		Rows    []row
		Scoring string
	}{Scoring: param(r, "scoring")}

	if data.Scoring != "bytes" && data.Scoring != "chars" {
		NotFound(w, r)
		return
	}

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

	render(w, r, "rankings/langs", data, "Rankings: Languages",
		"All languages in "+data.Scoring+".")
}
