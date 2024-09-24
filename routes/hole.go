package routes

import (
	"net/http"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/session"
	"github.com/lib/pq"
)

// GET /{hole}
func holeGET(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Authors                  []string
		HideDetails              bool
		Hole, PrevHole, NextHole *config.Hole
		Langs                    map[string]*config.Lang
		RankingsView             string
		Solutions                []map[string]string
		IsSponsor                bool
		HasNotes                 bool
	}{
		Langs:        config.AllLangByID,
		RankingsView: "me",
		Solutions:    []map[string]string{{}, {}},
	}

	var ok bool
	if data.Hole, ok = config.AllHoleByID[param(r, "hole")]; !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	data.PrevHole, data.NextHole = getPrevNextHole(r, data.Hole)

	if c, _ := r.Cookie("hide-details"); c != nil {
		data.HideDetails = true
	}

	if c, _ := r.Cookie("rankings-view"); c != nil {
		if c.Value == "top" || c.Value == "following" {
			data.RankingsView = c.Value
		}
	}

	// Lookup the hole's author(s).
	if data.Hole.Experiment == 0 {
		if err := session.Database(r).QueryRow(
			`SELECT array_agg(login ORDER BY login)
			   FROM authors
			   JOIN users ON id = user_id
			  WHERE hole = $1`,
			data.Hole.ID,
		).Scan(pq.Array(&data.Authors)); err != nil {
			panic(err)
		}
	}

	golfer := session.Golfer(r)

	if golfer != nil {
		data.IsSponsor = golfer.Admin || golfer.Sponsor
		if !data.IsSponsor {
			var notesCount int
			if err := session.Database(r).QueryRow(
				`SELECT COUNT(*) FROM notes WHERE user_id = $1`,
				session.Golfer(r).ID,
			).Scan(&notesCount); err != nil {
				panic(err)
			}
			data.HasNotes = notesCount > 0
		}
	}

	if golfer != nil && data.Hole.Experiment == 0 {
		// Fetch all the code per lang.
		rows, err := session.Database(r).Query(
			`SELECT code, lang, scoring
			   FROM solutions
			  WHERE hole = $1 AND user_id = $2`,
			data.Hole.ID, golfer.ID,
		)
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		for rows.Next() {
			var code, lang, scoring string

			if err := rows.Scan(&code, &lang, &scoring); err != nil {
				panic(err)
			}

			solution := 0
			if scoring == "chars" {
				solution = 1
			}

			data.Solutions[solution][lang] = code
		}

		if err := rows.Err(); err != nil {
			panic(err)
		}
	}

	view := "hole"
	if golfer != nil && golfer.Settings["hole"]["multi-window-layout"].(bool) {
		view = "hole-tabs"
	}

	render(w, r, view, data, data.Hole.Name, data.Hole.Synopsis)
}
