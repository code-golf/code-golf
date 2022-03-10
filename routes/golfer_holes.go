package routes

import (
	"net/http"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/session"
)

// GolferHoles serves GET /golfers/{golfer}/holes
func GolferHoles(w http.ResponseWriter, r *http.Request) {
	golfer := session.GolferInfo(r).Golfer
	data := struct {
		Holes    []*config.Hole
		Langs    []*config.Lang
		Ranks    map[string]int
		Scoring  string
		Scorings []string
	}{
		Holes:    config.HoleList,
		Langs:    config.LangList,
		Ranks:    map[string]int{},
		Scoring:  param(r, "scoring"),
		Scorings: []string{"bytes", "chars"},
	}

	if data.Scoring != "bytes" && data.Scoring != "chars" {
		if data.Scoring == "" {
			http.Redirect(w, r, "/golfers/"+golfer.Name+"/holes/chars", http.StatusPermanentRedirect)
			return
		}

		NotFound(w, r)
		return
	}

	rows, err := session.Database(r).Query(
		`SELECT hole::text || lang::text || scoring::text, rank
		   FROM rankings WHERE user_id = $1`,
		golfer.ID,
	)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var key string
		var rank int

		if err := rows.Scan(&key, &rank); err != nil {
			panic(err)
		}

		data.Ranks[key] = rank
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	render(w, r, "golfer/holes", data, golfer.Name)
}
