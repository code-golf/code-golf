package routes

import (
	"net/http"
	"time"

	"github.com/code-golf/code-golf/session"
)

// GET /golfers/{golfer}/cheevos
func golferCheevosGET(w http.ResponseWriter, r *http.Request) {
	golfer := session.GolferInfo(r).Golfer

	type Progress struct {
		Count, Percent, Progress int
		Earned                   *time.Time
	}

	data := map[string]Progress{}

	db := session.Database(r)
	rows, err := db.Query(
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
	defer rows.Close()

	for rows.Next() {
		var cheevoID string
		var progress Progress

		if err := rows.Scan(
			&cheevoID,
			&progress.Count,
			&progress.Earned,
			&progress.Percent,
		); err != nil {
			panic(err)
		}

		data[cheevoID] = progress
	}

	// Caclulate progress
	// TODO Bake it into the cheevos table rather than calculating on the fly.
	cheevoProgress := func(sql string, cheevoIDs ...string) {
		var count int
		if err := db.QueryRow(sql, golfer.ID).Scan(&count); err != nil {
			panic(err)
		}

		for _, cheevoID := range cheevoIDs {
			progress := data[cheevoID]
			progress.Progress = count
			data[cheevoID] = progress
		}
	}

	cheevoProgress(
		`SELECT COUNT(DISTINCT hole)
		   FROM solutions
		  WHERE NOT failing AND user_id = $1`,
		"up-to-eleven", "bakers-dozen", "the-watering-hole", "blackjack",
		"rule-34", "forty-winks", "dont-panic", "bullseye",
		"gone-in-60-holes", "cunning-linguist",
	)

	cheevoProgress(
		`WITH langs AS (
		    SELECT COUNT(DISTINCT lang)
		      FROM solutions
		     WHERE NOT failing AND user_id = $1
		  GROUP BY hole
		) SELECT COALESCE(MAX(count), 0) FROM langs`,
		"polyglot", "polyglutton", "omniglot",
	)

	cheevoProgress(
		"SELECT COALESCE(MAX(points), 0) FROM points WHERE user_id = $1",
		"its-over-9000", "twenty-kiloleagues", "marathon-runner",
	)

	render(w, r, "golfer/cheevos", data, golfer.Name)
}
