package routes

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/pager"
	"github.com/code-golf/code-golf/session"
)

// RecentSolutions serves GET /recent/solutions/{hole}/{lang}/{scoring}
func RecentSolutions(w http.ResponseWriter, r *http.Request) {
	type row struct {
		Hole                    *config.Hole
		Lang                    *config.Lang
		Login                   string
		Rank, Strokes, TieCount int
		Submitted               time.Time
	}

	data := struct {
		HoleID, LangID, Scoring string
		Holes                   []*config.Hole
		Langs                   []*config.Lang
		Pager                   *pager.Pager
		Rows                    []row
	}{
		HoleID:  param(r, "hole"),
		Holes:   config.HoleList,
		LangID:  param(r, "lang"),
		Langs:   config.LangList,
		Pager:   pager.New(r),
		Rows:    make([]row, 0, pager.PerPage),
		Scoring: param(r, "scoring"),
	}

	if data.HoleID != "all" && config.HoleByID[data.HoleID] == nil ||
		data.LangID != "all" && config.LangByID[data.LangID] == nil ||
		data.Scoring != "chars" && data.Scoring != "bytes" {
		NotFound(w, r)
		return
	}

	rows, err := session.Database(r).Query(
		` SELECT hole, lang, login, strokes, rank, submitted, tie_count
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
	)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		var holeID, langID string
		var r row

		if err := rows.Scan(
			&holeID,
			&langID,
			&r.Login,
			&r.Strokes,
			&r.Rank,
			&r.Submitted,
			&r.TieCount,
		); err != nil {
			panic(err)
		}

		r.Hole = config.HoleByID[holeID]
		r.Lang = config.LangByID[langID]

		data.Rows = append(data.Rows, r)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	render(w, r, "recent/solutions", data, "Recent Solutions")
}
