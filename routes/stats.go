package routes

import (
	"fmt"
	"net/http"

	"github.com/code-golf/code-golf/hole"
	"github.com/code-golf/code-golf/lang"
	"github.com/code-golf/code-golf/pie"
	"github.com/code-golf/code-golf/session"
)

// Stats serves GET /stats
func Stats(w http.ResponseWriter, r *http.Request) {
	type row struct {
		Count, Golfers, Rank   int
		Fact, PerGolfer, Route string
	}

	type table struct {
		Fact string
		Rows []row
	}

	data := struct {
		Bytes, Golfers, Holes, Langs, Solutions int
		SolutionsByHole, SolutionsByLang        pie.Pie
		Tables                                  []table
	}{
		Holes:  len(hole.List),
		Langs:  len(lang.List),
		Tables: []table{{Fact: "Hole"}, {Fact: "Language"}},
	}

	db := session.Database(r)

	if err := db.QueryRow(
		"SELECT COUNT(DISTINCT user_id) FROM trophies",
	).Scan(&data.Golfers); err != nil {
		panic(err)
	}

	if err := db.QueryRow(
		`WITH distinct_solutions AS (
		  SELECT DISTINCT ON (hole, lang, user_id, code_id) bytes
		    FROM solutions
		    JOIN code ON code_id = code.id
		   WHERE NOT failing
		) SELECT COALESCE(SUM(bytes), 0), COUNT(*) FROM distinct_solutions`,
	).Scan(&data.Bytes, &data.Solutions); err != nil {
		panic(err)
	}

	for i, fact := range [...]string{"hole", "lang"} {
		rows, err := db.Query(
			`WITH distinct_solutions AS (
			  SELECT DISTINCT ON (hole, lang, user_id, code_id)
			         hole,
			         lang,
			         user_id,
			         code_id
			    FROM solutions
			   WHERE NOT failing
			) SELECT RANK() OVER (ORDER BY COUNT(*) DESC, ` + fact + `),
			         ` + fact + `,
			         COUNT(*),
			         COUNT(DISTINCT user_id)
			    FROM distinct_solutions
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
				row.Route = "/" + row.Fact
				row.Fact = hole.ByID[row.Fact].Name
			} else {
				row.Route = "/recent/" + row.Fact
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

	render(w, r, "stats", data, "Stats")
}
