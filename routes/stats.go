package routes

import (
	"fmt"
	"net/http"

	"github.com/code-golf/code-golf/hole"
	"github.com/code-golf/code-golf/lang"
	"github.com/code-golf/code-golf/middleware"
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

	db := middleware.Database(r)

	if err := db.QueryRow(
		`SELECT (SELECT COUNT(DISTINCT user_id) FROM trophies),
		        (SELECT COUNT(*) FROM solutions WHERE NOT failing)`,
	).Scan(&data.Golfers, &data.Solutions); err != nil {
		panic(err)
	}

	for i, fact := range [...]string{"hole", "lang"} {
		rows, err := db.Query(
			` SELECT RANK() OVER (ORDER BY COUNT(*) DESC, ` + fact + `),
			         ` + fact + `, COUNT(*), COUNT(DISTINCT user_id)
			    FROM solutions
			   WHERE NOT failing
			GROUP BY ` + fact,
		)
		if err != nil {
			panic(err)
		}

		var slices []pie.Slice

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

			// Pie chart
			switch len(slices) {
			case 6:
				slices[5].Quantity += row.Count
			case 5:
				slices = append(slices, pie.Slice{Label: "Other", Quantity: row.Count})
			default:
				slices = append(slices, pie.Slice{Label: row.Fact, Quantity: row.Count})
			}
		}

		if err := rows.Err(); err != nil {
			panic(err)
		}

		if fact == "hole" {
			data.SolutionsByHole = pie.New(slices)
		} else {
			data.SolutionsByLang = pie.New(slices)
		}
	}

	render(w, r, "stats", "Stats", data)
}
