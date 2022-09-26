package routes

import (
	"net/http"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/session"
)

// GET /golfers/{golfer}/rankings/{scoring}
func golferRankingsGET(w http.ResponseWriter, r *http.Request) {
	type ranking struct {
		Diamond               bool
		Golfers, Rank, Points int
	}

	data := struct {
		Holes     []*config.Hole
		Langs     []*config.Lang
		LangsUsed map[string]bool
		Rankings  map[string]map[string]*ranking
		Display   string
		Displays  []string
		Scope     string
		Scopes    []string
		Scoring   string
		Scorings  []string
	}{
		Holes:     config.HoleList,
		Langs:     config.LangList,
		LangsUsed: map[string]bool{},
		Rankings:  map[string]map[string]*ranking{},
		Display:   param(r, "display"),
		Displays:  []string{"rankings", "points"},
		Scope:     param(r, "scope"),
		Scopes:    []string{"lang", "overall"},
		Scoring:   param(r, "scoring"),
		Scorings:  []string{"bytes", "chars"},
	}

	golfer := session.GolferInfo(r).Golfer
	var points = "points_for_lang"
	var rank = "rank"
	if data.Scope == "overall" {
		points = "points"
		rank = "rank_overall"
	}
	rows, err := session.Database(r).Query(
		"SELECT hole, lang, golfers, "+rank+", "+rank+" = 1 AND tie_count = 1, "+
			points+" FROM rankings WHERE user_id = $1 AND scoring = $2",
		golfer.ID,
		data.Scoring,
	)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var hole, lang string
		var r ranking

		if err := rows.Scan(
			&hole, &lang, &r.Golfers, &r.Rank, &r.Diamond, &r.Points,
		); err != nil {
			panic(err)
		}

		data.LangsUsed[lang] = true

		if _, ok := data.Rankings[hole]; !ok {
			data.Rankings[hole] = map[string]*ranking{}
		}

		data.Rankings[hole][lang] = &r
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	render(w, r, "golfer/rankings", data, golfer.Name)
}
