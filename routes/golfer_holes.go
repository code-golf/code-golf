package routes

import (
	"net/http"

	"github.com/code-golf/code-golf/golfer"
	"github.com/code-golf/code-golf/lang"
)

// GolferHoles serves GET /golfers/{golfer}/holes
func GolferHoles(w http.ResponseWriter, r *http.Request) {
	golfer := r.Context().Value("golferInfo").(*golfer.GolferInfo).Golfer

	data := struct {
		Holes []Hole
		Langs []lang.Lang
		Ranks map[string]map[string]int
	}{holes, lang.List, map[string]map[string]int{}}

	rows, err := db(r).Query(
		`WITH matrix AS (
		  SELECT user_id,
		         hole,
		         lang,
		         RANK() OVER (PARTITION BY hole, lang ORDER BY LENGTH(code))
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

	render(w, r, http.StatusOK, "golfer-holes", golfer.Name, data)
}
