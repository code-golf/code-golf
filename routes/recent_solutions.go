package routes

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/pager"
	"github.com/code-golf/code-golf/session"
)

// GET /recent/solutions/{hole}/{lang}/{scoring}
func recentSolutionsGET(w http.ResponseWriter, r *http.Request) {
	type row struct {
		Country                 config.NullCountry
		Name                    string
		Hole                    *config.Hole
		Lang                    *config.Lang
		Rank, Strokes, TieCount int
		Submitted               time.Time
	}

	data := struct {
		HoleID, LangID, Scoring string
		Holes                   []*config.Hole
		Langs                   []*config.Lang
		LangsShown              map[string]bool
		Pager                   *pager.Pager
		Rows                    []row
	}{
		HoleID:     param(r, "hole"),
		Holes:      config.HoleList,
		LangID:     param(r, "lang"),
		Langs:      config.LangList,
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

	if err := session.Database(r).Select(
		&data.Rows,
		` SELECT hole, lang, login name, strokes, rank, submitted, tie_count
		    FROM rankings
		    JOIN users ON user_id = id
		   WHERE (hole = $1 OR $1 IS NULL)
		     AND (lang = $2 OR $2 IS NULL)
		     AND scoring = $3
		ORDER BY submitted DESC LIMIT $4`,
		sql.NullString{String: data.HoleID, Valid: data.HoleID != "all"},
		sql.NullString{String: data.LangID, Valid: data.LangID != "all"},
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
