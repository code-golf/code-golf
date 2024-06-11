package routes

import (
	"net/http"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/session"
)

// GET /golfers/{golfer}/holes/{display}/{scope}/{scoring}
func golferHolesGET(w http.ResponseWriter, r *http.Request) {
	type ranking struct{ Golfers, Points, Rank, TieCount int }

	data := struct {
		Holes                      []*config.Hole
		Langs                      []*config.Lang
		LangsUsed                  map[string]bool
		Rankings                   map[string]map[string]*ranking
		Display, Scope, Scoring    string
		Displays, Scopes, Scorings []string
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
	rows, err := session.Database(r).Query(
		`SELECT hole, lang, tie_count,
		        CASE WHEN $1 THEN golfers_overall ELSE golfers         END,
		        CASE WHEN $1 THEN points          ELSE points_for_lang END,
		        CASE WHEN $1 THEN rank_overall    ELSE rank            END
		   FROM rankings
		  WHERE user_id = $2 AND scoring = $3`,
		data.Scope == "overall",
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
			&hole, &lang, &r.TieCount, &r.Golfers, &r.Points, &r.Rank,
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

	render(w, r, "golfer/holes", data, golfer.Name)
}
