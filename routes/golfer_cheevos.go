package routes

import (
	"database/sql"
	"net/http"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/session"
)

// GolferCheevos serves GET /golfers/{golfer}
func GolferCheevos(w http.ResponseWriter, r *http.Request) {
	golfer := session.GolferInfo(r).Golfer

	type EarnedCheevo struct {
		Count, Percent int
		Earned         sql.NullTime
	}

	data := struct {
		Cheevos map[string][]*config.Cheevo
		Earned  map[string]EarnedCheevo
	}{config.CheevoTree, map[string]EarnedCheevo{}}

	rows, err := session.Database(r).Query(
		`WITH count AS (
		    SELECT trophy, COUNT(*) FROM trophies GROUP BY trophy
		), earned AS (
		    SELECT trophy, earned FROM trophies WHERE user_id = $1
		) SELECT *, count * 100 / (SELECT COUNT(*) FROM users)
		    FROM count LEFT JOIN earned USING(trophy)`,
		golfer.ID,
	)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var cheevoID string
		var earned EarnedCheevo

		if err := rows.Scan(
			&cheevoID,
			&earned.Count,
			&earned.Earned,
			&earned.Percent,
		); err != nil {
			panic(err)
		}

		data.Earned[cheevoID] = earned
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	render(w, r, "golfer/cheevos", data, golfer.Name)
}
