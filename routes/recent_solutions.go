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
		Country                          *config.Country
		Experimental                     bool
		Name                             string
		Hole                             *config.Hole
		Lang                             *config.Lang
		Golfers, Rank, Strokes, TieCount int
		Submitted                        time.Time
		Time                             *time.Duration
	}

	data := struct {
		Hole, PrevHole, NextHole *config.Hole
		HoleID, LangID, Scoring  string
		Pager                    *pager.Pager
		Rows                     []row
	}{
		HoleID:  param(r, "hole"),
		LangID:  param(r, "lang"),
		Pager:   pager.New(r),
		Rows:    make([]row, 0, pager.PerPage),
		Scoring: param(r, "scoring"),
	}

	if data.HoleID != "all" && config.AllHoleByID[data.HoleID] == nil ||
		data.LangID != "all" && config.AllLangByID[data.LangID] == nil ||
		data.Scoring != "chars" && data.Scoring != "bytes" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if data.Hole = config.AllHoleByID[data.HoleID]; data.Hole != nil {
		data.PrevHole, data.NextHole = getPrevNextHole(r, data.Hole)
	}

	if err := session.Database(r).Select(
		&data.Rows,
		` SELECT experimental, golfers, hole, lang, login name, strokes, rank,
		         submitted, tie_count, time_ms * 1e6 time
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

	render(w, r, "recent/solutions", data, "Recent Solutions")
}
