package routes

import (
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func stats(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := struct {
		DistributionOfHoles, DistributionOfLangs, HolesByDifficulty, SolutionsByLang template.JS
		Golfers, Holes, Langs, Solutions                                             int
		StatsCssPath, StatsJsPath                                                    string
	}{
		Holes:             len(holes),
		HolesByDifficulty: template.JS(HolesByDifficulty),
		Langs:             len(langs),
		StatsCssPath:      statsCssPath,
		StatsJsPath:       statsJsPath,
	}

	if err := db.QueryRow(
		"SELECT COUNT(DISTINCT user_id), COUNT(*) FROM solutions WHERE NOT failing",
	).Scan(&data.Golfers, &data.Solutions); err != nil {
		panic(err)
	}

	if err := db.QueryRow(
		`WITH top AS (
		    SELECT lang, lang::text lang_text, COUNT(*)
		      FROM solutions
		     WHERE NOT failing
		  GROUP BY lang
		  ORDER BY count DESC
		     LIMIT 7
		), data AS (
		    (SELECT lang_text, count FROM top ORDER BY lang)
		    UNION ALL
		    SELECT 'other' lang_text,
		           COUNT(*)
		      FROM solutions
		     WHERE NOT FAILING
		       AND lang NOT IN (SELECT lang FROM top)
		)
		SELECT '[' ||
		       ARRAY_TO_JSON(ARRAY_AGG(lang_text))
		       || ',' ||
		       ARRAY_TO_JSON(ARRAY_AGG(count))
		       || ']'
		  FROM data`,
	).Scan(&data.SolutionsByLang); err != nil {
		panic(err)
	}

	if err := db.QueryRow(
		`SELECT ARRAY_TO_JSON(ARRAY_AGG(ROW_TO_JSON(t)))
		   FROM (SELECT x, COUNT(*) y
		   FROM (SELECT COUNT(DISTINCT hole) x FROM solutions WHERE NOT failing GROUP BY user_id) z
		GROUP BY x ORDER BY x) t`,
	).Scan(&data.DistributionOfHoles); err != nil {
		panic(err)
	}

	if err := db.QueryRow(
		`SELECT ARRAY_TO_JSON(ARRAY_AGG(ROW_TO_JSON(t)))
		   FROM (SELECT x, COUNT(*) y
		   FROM (SELECT COUNT(DISTINCT lang) x FROM solutions WHERE NOT failing GROUP BY user_id) z
		GROUP BY x ORDER BY x) t`,
	).Scan(&data.DistributionOfLangs); err != nil {
		panic(err)
	}

	Render(w, r, http.StatusOK, "stats", "Stats", data)
}
