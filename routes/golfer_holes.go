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
	scoring := param(r, "scoring")

	if scoring != "bytes" && scoring != "chars" {
		if scoring == "" {
			http.Redirect(w, r, "/golfers/"+golfer.Name+"/holes/chars", http.StatusPermanentRedirect)
			return
		}

		NotFound(w, r)
		return
	}

	data := struct {
		Holes        []hole.Hole
		Langs        []lang.Lang
		Ranks        []map[string]map[string]int
		ScoringIndex int
		Scorings     []string
	}{
		Holes:    hole.List,
		Langs:    lang.List,
		Ranks:    []map[string]map[string]int{{}, {}},
		Scorings: []string{"bytes", "chars"},
	}

	if scoring == data.Scorings[1] {
		data.ScoringIndex = 1
	}

	rows, err := session.Database(r).Query(
		`WITH matrix AS (
		  SELECT user_id,
		         hole,
		         lang,
		         RANK() OVER (PARTITION BY hole, lang ORDER BY `+data.Scorings[0]+`)
		    FROM solutions
		   WHERE NOT failing
		     AND scoring = $2
		), other_matrix AS (
		  SELECT user_id,
		         hole,
		         lang,
		         RANK() OVER (PARTITION BY hole, lang ORDER BY `+data.Scorings[1]+`) other_rank
		    FROM solutions
		   WHERE NOT failing
		     AND scoring = $3
		) SELECT hole, lang, COALESCE(rank, 0), COALESCE(other_rank, 0)
		    FROM matrix
	   FULL JOIN other_matrix USING (user_id, hole, lang)
		   WHERE user_id = $1`,
		golfer.ID,
		data.Scorings[0],
		data.Scorings[1],
	)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		var hole, lang string
		var rank0, rank1 int

		if err := rows.Scan(&hole, &lang, &rank0, &rank1); err != nil {
			panic(err)
		}

		if data.Ranks[0][hole] == nil {
			data.Ranks[0][hole] = map[string]int{}
			data.Ranks[1][hole] = map[string]int{}
		}

		data.Ranks[0][hole][lang] = rank0
		data.Ranks[1][hole][lang] = rank1
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	render(w, r, "golfer/holes", data, golfer.Name)
}
