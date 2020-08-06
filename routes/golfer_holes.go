package routes

import (
	"net/http"

	"github.com/code-golf/code-golf/hole"
	"github.com/code-golf/code-golf/lang"
	"github.com/code-golf/code-golf/session"
)

// GolferHoles serves GET /golfers/{golfer}/holes
func GolferHoles(w http.ResponseWriter, r *http.Request) {
	golfer := session.GolferInfo(r).Golfer

	data := struct {
		Holes []hole.Hole
		Langs []lang.Lang
		Ranks map[string]map[string]int
	}{hole.List, lang.List, map[string]map[string]int{}}

	rows, err := session.Database(r).Query(
		`WITH matrix AS (
		  SELECT user_id,
		         hole,
		         lang,
		         RANK() OVER (PARTITION BY hole, lang ORDER BY chars)
		    FROM solutions
		   WHERE NOT failing
		) SELECT hole, lang, rank
		    FROM matrix
		   WHERE user_id = $1`,
		golfer.ID,
	)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		var hole, lang string
		var rank int

		if err := rows.Scan(&hole, &lang, &rank); err != nil {
			panic(err)
		}

		if data.Ranks[hole] == nil {
			data.Ranks[hole] = map[string]int{}
		}

		data.Ranks[hole][lang] = rank
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	render(w, r, "golfer/holes", golfer.Name, data)
}
