package routes

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/discord"
	h "github.com/code-golf/code-golf/hole"
	"github.com/code-golf/code-golf/session"
)

// GolferSolution serves GET /golfers/{golfer}/{hole}/{lang}/{scoring}
func GolferSolution(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Failing bool
		Hole    *config.Hole
		Lang    *config.Lang
		Scoring string
	}{
		Hole:    config.HoleByID[param(r, "hole")],
		Lang:    config.LangByID[param(r, "lang")],
		Scoring: param(r, "scoring"),
	}

	if data.Hole == nil || data.Lang == nil ||
		(data.Scoring != "bytes" && data.Scoring != "chars") {
		NotFound(w, r)
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
		NotFound(w, r)
		return
	} else if err != nil {
		panic(err)
	}

	render(w, r, "golfer/solution", data, golfer.Name)
}

// GolferSolutionPost serves POST /golfers/{golfer}/{hole}/{lang}/{scoring}
func GolferSolutionPost(w http.ResponseWriter, r *http.Request) {
	hole := config.HoleByID[param(r, "hole")]
	lang := config.LangByID[param(r, "lang")]
	scoring := param(r, "scoring")

	if hole == nil || lang == nil || (scoring != "bytes" && scoring != "chars") {
		NotFound(w, r)
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
		NotFound(w, r)
		return
	} else if err != nil {
		panic(err)
	}

	// Best of three runs. Given they're fickle, timeouts count as passes.
	var passes, fails int
	for i := 0; passes < 2 && fails < 2; i++ {
		score := h.Play(ctx, hole.ID, lang.ID, code)
		if score.Pass || score.Timeout {
			passes++
		} else {
			fails++
		}
	}

	if fails > passes {
		if _, err := db.Exec(
			`UPDATE solutions
			    SET failing = true
			  WHERE code = $1 AND hole = $2 AND lang = $3`,
			code, hole.ID, lang.ID,
		); err != nil {
			panic(err)
		}

		go discord.LogFailedRejudge(&golfer, hole, lang, scoring)
	}

	http.Redirect(w, r, param(r, "scoring"), http.StatusSeeOther)
}
