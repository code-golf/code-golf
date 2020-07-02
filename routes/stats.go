package routes

import (
	"fmt"
	"net/http"

	"github.com/code-golf/code-golf/hole"
	"github.com/code-golf/code-golf/lang"
	"github.com/code-golf/code-golf/pie"
)

// Stats serves GET /stats
func Stats(w http.ResponseWriter, r *http.Request) {
	type row struct {
		Count, Golfers, Rank int
		Fact, PerGolfer      string
	}

	type table struct {
		Fact string
		Rows []row
	}

	data := struct {
		Golfers, Holes, Langs, Solutions int
		SolutionsByHole, SolutionsByLang pie.Pie
		Tables                           []table
	}{
		Holes:  len(hole.List),
		Langs:  len(lang.List),
		Tables: []table{{Fact: "Hole"}, {Fact: "Language"}},
	}

	db := db(r)

	if err := db.QueryRow(
		`SELECT (SELECT COUNT(DISTINCT user_id) FROM trophies),
		        (SELECT COUNT(*) FROM solutions WHERE NOT failing)`,
	).Scan(&data.Golfers, &data.Solutions); err != nil {
		panic(err)
	}

	for i, fact := range [...]string{"hole", "lang"} {
		rows, err := db.Query(
			`WITH top AS (
			    SELECT ` + fact + `, ` + fact + `::text txt, COUNT(*)
			      FROM solutions
			     WHERE NOT failing
			  GROUP BY ` + fact + `
			  ORDER BY count DESC, ` + fact + `
			     LIMIT 5
			) (SELECT txt, count FROM top)
			    UNION ALL
			SELECT 'Other' txt,
			       COUNT(*)
			  FROM solutions
			 WHERE NOT failing
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
					slice.Label = hole.ByID[slice.Label].Name
				} else {
					slice.Label = lang.ByID[slice.Label].Name
				}
			}

			slices = append(slices, slice)
		}

		if err := rows.Err(); err != nil {
			panic(err)
		}

		if fact == "hole" {
			data.SolutionsByHole = pie.New(slices)
		} else {
			data.SolutionsByLang = pie.New(slices)
		}

		rows, err = db.Query(
			` SELECT RANK() OVER (ORDER BY COUNT(*) DESC, ` + fact + `),
			         ` + fact + `, COUNT(*), COUNT(DISTINCT user_id)
			    FROM solutions
			   WHERE NOT failing
			GROUP BY ` + fact,
		)
		if err != nil {
			panic(err)
		}

		for rows.Next() {
			var row row

			if err := rows.Scan(
				&row.Rank,
				&row.Fact,
				&row.Count,
				&row.Golfers,
			); err != nil {
				panic(err)
			}

			if fact == "hole" {
				row.Fact = hole.ByID[row.Fact].Name
			} else {
				row.Fact = lang.ByID[row.Fact].Name
			}

			row.PerGolfer = fmt.Sprintf("%.2f", float64(row.Count)/float64(row.Golfers))

			data.Tables[i].Rows = append(data.Tables[i].Rows, row)
		}

		if err := rows.Err(); err != nil {
			panic(err)
		}
	}

	render(w, r, http.StatusOK, "stats", "Stats", data)
}
