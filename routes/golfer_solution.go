package routes

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/discord"
	"github.com/code-golf/code-golf/session"
)

// GET /golfers/{golfer}/{hole}/{lang}/{scoring}
func golferSolutionGET(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Failing bool
		Hole    *config.Hole
		Lang    *config.Lang
		Log     []struct {
			Bytes, Chars int
			Submitted    time.Time
		}
		Scoring string
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
	if err := session.Database(r).QueryRow(
		`SELECT failing
		   FROM solutions
		  WHERE hole = $1 AND lang = $2 AND scoring = $3 AND user_id = $4`,
		data.Hole.ID,
		data.Lang.ID,
		data.Scoring,
		golfer.ID,
	).Scan(&data.Failing); errors.Is(err, sql.ErrNoRows) {
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
// nolint deadcode
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

	// Best of three runs. Given they're fickle, timeouts count as passes.
	var passes, fails int
	/* FIXME for i := 0; passes < 2 && fails < 2; i++ {
		score := h.Play(ctx, hole.ID, lang.ID, code)
		if score.Pass || score.Timeout || score.ExitCode != 0 {
			passes++
		} else {
			fails++
		}
	} */

	if fails > passes {
		res := db.MustExec(
			`UPDATE solutions
			    SET failing = true
			  WHERE NOT failing AND code = $1 AND hole = $2 AND lang = $3`,
			code, hole.ID, lang.ID,
		)

		// FIXME Technically we can end up failing multiple golfers.
		if rows, _ := res.RowsAffected(); rows > 0 {
			go discord.LogFailedRejudge(&golfer, hole, lang, scoring)
		}
	}

	http.Redirect(w, r, param(r, "scoring"), http.StatusSeeOther)
}
