package routes

import (
	"net/http"
	"time"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/null"
	"github.com/code-golf/code-golf/pager"
	"github.com/code-golf/code-golf/session"
)

// GET /recent/solutions/{hole}/{lang}/{scoring}
func recentSolutionsGET(w http.ResponseWriter, r *http.Request) {
	type row struct {
		Country                          config.NullCountry
		Name                             string
		Hole                             *config.Hole
		Lang                             *config.Lang
		Golfers, Rank, Strokes, TieCount int
		Submitted                        time.Time
	}

	data := struct {
		Hole, PrevHole, NextHole *config.Hole
		HoleID, LangID, Scoring  string
		LangsShown               map[string]bool
		Pager                    *pager.Pager
		Rows                     []row
	}{
		HoleID:     param(r, "hole"),
		LangID:     param(r, "lang"),
		LangsShown: map[string]bool{},
		Pager:      pager.New(r),
		Rows:       make([]row, 0, pager.PerPage),
		Scoring:    param(r, "scoring"),
	}

	if data.HoleID != "all" && config.HoleByID[data.HoleID] == nil ||
		data.LangID != "all" && config.LangByID[data.LangID] == nil ||
		data.Scoring != "chars" && data.Scoring != "bytes" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if data.Hole = config.HoleByID[data.HoleID]; data.Hole != nil {
		data.PrevHole, data.NextHole = getPrevNextHole(r, data.Hole)
	}

	if err := session.Database(r).Select(
		&data.Rows,
		` SELECT golfers, hole, lang, login name, strokes, rank, submitted, tie_count
		    FROM rankings
		    JOIN users ON user_id = id
		   WHERE (hole = $1 OR $1 IS NULL)
		     AND (lang = $2 OR $2 IS NULL)
		     AND scoring = $3
		ORDER BY submitted DESC LIMIT $4`,
		null.New(data.HoleID, data.HoleID != "all"),
		null.New(data.LangID, data.LangID != "all"),
		data.Scoring,
		pager.PerPage,
	); err != nil {
		panic(err)
	}

	for _, row := range data.Rows {
		data.LangsShown[row.Lang.ID] = true
	}

	render(w, r, "recent/solutions", data, "Recent Solutions")
}
