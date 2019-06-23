package routes

import (
	"net/http"

	"github.com/JRaspass/code-golf/pie"
	"github.com/julienschmidt/httprouter"
)

func stats(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := struct {
		Golfers, Holes, Langs, Solutions int
		SolutionsByHole, SolutionsByLang pie.Pie
	}{
		Holes: len(holes),
		Langs: len(langs),
	}

	if err := db.QueryRow(
		"SELECT COUNT(DISTINCT user_id), COUNT(*) FROM solutions WHERE NOT failing",
	).Scan(&data.Golfers, &data.Solutions); err != nil {
		panic(err)
	}

	for _, fact := range [...]string{"hole", "lang"} {
		rows, err := db.Query(
			`WITH top AS (
			    SELECT ` + fact + `, ` + fact + `::text txt, COUNT(*)
			      FROM solutions
			     WHERE NOT failing
			  GROUP BY ` + fact + `
			  ORDER BY count DESC
			     LIMIT 5
			) (SELECT txt, count FROM top)
			    UNION ALL
			SELECT 'Other' txt,
			       COUNT(*)
			  FROM solutions
			 WHERE NOT FAILING
			   AND ` + fact + ` NOT IN (SELECT ` + fact + ` FROM top)`,
		)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var slices []pie.Slice

		for rows.Next() {
			var slice pie.Slice

			if err := rows.Scan(&slice.Label, &slice.Quantity); err != nil {
				panic(err)
			}

			if slice.Label != "Other" {
				if fact == "hole" {
					slice.Label = holeByID[slice.Label].Name
				} else {
					slice.Label = langByID[slice.Label].Name
				}
			}

			slices = append(slices, slice)
		}

		if fact == "hole" {
			data.SolutionsByHole = pie.New(slices)
		} else {
			data.SolutionsByLang = pie.New(slices)
		}
	}

	Render(w, r, http.StatusOK, "stats", "Stats", data)
}
