package routes

import (
	"net/http"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/session"
)

// GET /golfers/{golfer}/holes
func golferHolesGET(w http.ResponseWriter, r *http.Request) {
	golfer := session.GolferInfo(r).Golfer
	type rank struct {
		rank      int
		isDiamond bool
	}
	data := struct {
		Holes    []*config.Hole
		Langs    []*config.Lang
		Ranks    map[string]rank
		Scoring  string
		Scorings []string
	}{
		Holes:    config.HoleList,
		Langs:    config.LangList,
		Ranks:    map[string]rank{},
		Scoring:  param(r, "scoring"),
		Scorings: []string{"bytes", "chars"},
	}

	if data.Scoring != "bytes" && data.Scoring != "chars" {
		if data.Scoring == "" {
			http.Redirect(w, r, "/golfers/"+golfer.Name+"/holes/chars", http.StatusPermanentRedirect)
			return
		}

		w.WriteHeader(http.StatusNotFound)
		return
	}

	rows, err := session.Database(r).Query(
		`SELECT hole::text || lang::text || scoring::text, rank, rank = 1 AND tie_count = 1
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
		var isDiamond bool

		if err := rows.Scan(&key, &rank, &isDiamond); err != nil {
			panic(err)
		}

		data.Ranks[key] = rank{rank, isDiamond}
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	render(w, r, "golfer/holes", data, golfer.Name)
}
