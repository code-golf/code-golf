package routes

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/code-golf/code-golf/config"
	h "github.com/code-golf/code-golf/hole"
	"github.com/code-golf/code-golf/session"
)

// GET /golfers/{golfer}/{hole}/{lang}/{scoring}
func golferSolutionGET(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Failing bool
		Hole    *config.Hole
		Lang    *config.Lang
		Log     []struct {
			Bytes     int
			Chars     *int
			Submitted time.Time
		}
		Rank, RankOverall, Row, RowOverall *int
		Scoring                            string
		Tested                             time.Time
	}{
		Hole:    config.HoleByID[param(r, "hole")],
		Lang:    config.LangByID[param(r, "lang")],
		Scoring: param(r, "scoring"),
	}

	if data.Hole == nil || data.Lang == nil ||
		(data.Scoring != "bytes" && data.Scoring != "chars") {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	golfer := session.GolferInfo(r).Golfer
	if err := session.Database(r).Get(
		&data,
		`  SELECT failing, rank, rank_overall, row, row_overall, tested
		     FROM solutions
		LEFT JOIN rankings USING (hole, lang, scoring, user_id)
		    WHERE hole = $1 AND lang = $2 AND scoring = $3 AND user_id = $4`,
		data.Hole.ID,
		data.Lang.ID,
		data.Scoring,
		golfer.ID,
	); errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		panic(err)
	}

	if err := session.Database(r).Select(
		&data.Log,
		` SELECT bytes, chars, submitted
		    FROM solutions_log
		   WHERE hole = $1 AND lang = $2 AND scoring = $3 AND user_id = $4
		ORDER BY submitted DESC`,
		data.Hole.ID,
		data.Lang.ID,
		data.Scoring,
		golfer.ID,
	); err != nil {
		panic(err)
	}

	render(w, r, "golfer/solution", data, golfer.Name)
}

// POST /golfers/{golfer}/{hole}/{lang}/{scoring}
func golferSolutionPOST(w http.ResponseWriter, r *http.Request) {
	hole := config.HoleByID[param(r, "hole")]
	lang := config.LangByID[param(r, "lang")]
	scoring := param(r, "scoring")

	if hole == nil || lang == nil || (scoring != "bytes" && scoring != "chars") {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	code := ""
	ctx := r.Context()
	db := session.Database(r)
	golfer := session.GolferInfo(r).Golfer

	if err := db.QueryRowContext(
		ctx,
		`SELECT code
		   FROM solutions
		  WHERE NOT failing
		    AND hole    = $1
		    AND lang    = $2
		    AND scoring = $3
		    AND user_id = $4`,
		hole.ID,
		lang.ID,
		scoring,
		golfer.ID,
	).Scan(&code); errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		panic(err)
	}

	// Subset runs so we don't leak too much info.
	type SubsetRun struct {
		Pass    bool          `json:"pass"`
		Time    time.Duration `json:"time"`
		Timeout bool          `json:"timeout"`
	}

	runs := h.Play(ctx, hole, lang, code)
	subsetRuns := make([]SubsetRun, len(runs))

	for i, run := range runs {
		subsetRuns[i] = SubsetRun{run.Pass, run.Time, run.Timeout}
	}

	if err := json.NewEncoder(w).Encode(subsetRuns); err != nil {
		panic(err)
	}
}
