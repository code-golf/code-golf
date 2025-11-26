package routes

import (
	"net/http"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/session"
)

// GET /{hole}
func holeGET(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Authors                      []struct{ AvatarURL, Name string }
		HideDetails                  bool
		Hole, PrevHole, NextHole     *config.Hole
		HoleRedirects, LangRedirects map[string]string
		Langs                        map[string]*config.Lang
		RankingsView                 string
		Solutions                    []map[string]string
	}{
		Hole:          config.AllHoleByID[param(r, "hole")],
		HoleRedirects: config.HoleRedirects,
		LangRedirects: config.LangRedirects,
		Langs:         config.AllLangByID,
		RankingsView:  "me",
		Solutions:     []map[string]string{{}, {}},
	}

	if data.Hole == nil {
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
	if err := session.Database(r).Select(
		&data.Authors,
		`SELECT avatar_url, name
		   FROM authors
		   JOIN golfers_with_avatars ON id = user_id
		  WHERE hole = $1`,
		data.Hole.ID,
	); err != nil {
		panic(err)
	}

	golfer := session.Golfer(r)

	if golfer != nil {
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
