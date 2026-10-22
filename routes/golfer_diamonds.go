package routes

import (
	"net/http"
	"time"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/pager"
	"github.com/code-golf/code-golf/session"
)

// GET /golfers/{golfer}/diamonds/{hole}/{lang}/{scoring}
func golferDiamondsGET(w http.ResponseWriter, r *http.Request) {
	data := struct {
		HoleID, LangID, Scoring string
		Pager                   *pager.Pager
		Rows                    []struct {
			Hole        *config.Hole
			Lang        *config.Lang
			Rank, Total int
			Scoring     string
			Submitted   time.Time
		}
	}{
		HoleID:  param(r, "hole"),
		LangID:  param(r, "lang"),
		Pager:   pager.New(r),
		Scoring: param(r, "scoring"),
	}

	if data.HoleID != "all" && config.HoleByID[data.HoleID] == nil ||
		data.LangID != "all" && config.LangByID[data.LangID] == nil ||
		data.Scoring != "all" && data.Scoring != "chars" && data.Scoring != "bytes" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	golfer := session.GolferInfo(r).Golfer
	if err := session.Database(r).Select(
		&data.Rows,
		`  SELECT hole, lang, scoring, submitted,
		          RANK() OVER(ORDER BY submitted),
		          COUNT(*) OVER () total
		     FROM medals
		    WHERE medal = 'diamond'
		      AND user_id = $1
		      AND $2 IN ('all', hole::text)
		      AND $3 IN ('all', lang::text)
		      AND $4 IN ('all', scoring::text)
		 ORDER BY rank, hole, lang, scoring
		    LIMIT $5 OFFSET $6`,
		golfer.ID,
		data.HoleID,
		data.LangID,
		data.Scoring,
		pager.PerPage,
		data.Pager.Offset,
	); err != nil {
		panic(err)
	}

	if len(data.Rows) > 0 {
		data.Pager.Total = data.Rows[0].Total
	}

	if data.Pager.Calculate() {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	render(w, r, "golfer/diamonds", data, golfer.Name)
}
