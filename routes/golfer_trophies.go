package routes

import (
	"database/sql"
	"net/http"

	"github.com/code-golf/code-golf/session"
	"github.com/code-golf/code-golf/trophy"
)

// GolferTrophies serves GET /golfers/{golfer}
func GolferTrophies(w http.ResponseWriter, r *http.Request) {
	golfer := session.GolferInfo(r).Golfer

	type EarnedTrophy struct {
		Count, Percent int
		Earned         sql.NullTime
	}

	data := struct {
		Earned   map[string]EarnedTrophy
		Trophies map[string][]*trophy.Trophy
	}{
		map[string]EarnedTrophy{},
		trophy.Tree,
	}

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
		var earned EarnedTrophy
		var trophyID string

		if err := rows.Scan(
			&trophyID,
			&earned.Count,
			&earned.Earned,
			&earned.Percent,
		); err != nil {
			panic(err)
		}

		data.Earned[trophyID] = earned
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	render(w, r, "golfer/trophies", data, golfer.Name)
}
