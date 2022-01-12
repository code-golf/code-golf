package routes

import (
	"net/http"
	"time"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/session"
)

// GolferWall serves GET /golfers/{golfer}
func GolferWall(w http.ResponseWriter, r *http.Request) {
	const limit = 100

	type row struct {
		Cheevo *config.Cheevo
		Date   time.Time
		Golfer string
		Hole   *config.Hole
		Lang   *config.Lang
	}

	data := make([]row, 0, limit)
	golfer := session.GolferInfo(r).Golfer

	// TODO Support friends/follow.
	rows, err := session.Database(r).Query(
		`WITH data AS (
		 -- Cheevos
		    SELECT earned       date,
		           trophy::text cheevo,
		           ''           hole,
		           ''           lang,
		           user_id
		      FROM trophies
		     WHERE user_id = $1
		 UNION ALL
		 -- Holes
		    SELECT MAX(submitted) date,
		           ''             cheevo,
		           hole::text     hole,
		           lang::text     lang,
		           user_id
		      FROM solutions
		     WHERE user_id = $1
		  GROUP BY user_id, hole, lang
		) SELECT cheevo, date, login, hole, lang
		    FROM data JOIN users ON id = user_id
		ORDER BY date DESC LIMIT $2`,
		golfer.ID,
		limit,
	)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		var cheevo, golfer, hole, lang string
		var date time.Time

		if err := rows.Scan(&cheevo, &date, &golfer, &hole, &lang); err != nil {
			panic(err)
		}

		// TODO Parse date into viewers location.
		data = append(data, row{
			Cheevo: config.CheevoByID[cheevo],
			Date:   date,
			Golfer: golfer,
			Hole:   config.HoleByID[hole],
			Lang:   config.LangByID[lang],
		})
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	render(w, r, "golfer/wall", data, golfer.Name)
}
