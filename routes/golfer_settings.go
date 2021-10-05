package routes

import (
	"net/http"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/session"
	"github.com/code-golf/code-golf/zone"
)

// GolferCancelDelete serves POST /golfer/cancel-delete
func GolferCancelDelete(w http.ResponseWriter, r *http.Request) {
	if _, err := session.Database(r).Exec(
		"UPDATE users SET delete = NULL WHERE id = $1",
		session.Golfer(r).ID,
	); err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/golfer/settings", http.StatusSeeOther)
}

// GolferDelete serves POST /golfer/delete
func GolferDelete(w http.ResponseWriter, r *http.Request) {
	if _, err := session.Database(r).Exec(
		"UPDATE users SET delete = TIMEZONE('UTC', NOW()) + INTERVAL '7 days' WHERE id = $1",
		session.Golfer(r).ID,
	); err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/golfer/settings", http.StatusSeeOther)
}

// GolferSettings serves GET /golfer/settings
func GolferSettings(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Countries map[string][]*config.Country
		Keymaps   []string
		Themes    []string
		TimeZones []zone.Zone
	}{
		config.CountryTree,
		[]string{"default", "vim"},
		[]string{"auto", "dark", "light"},
		zone.List(),
	}

	render(w, r, "golfer/settings", data, "Settings")
}

// GolferSettingsPost serves POST /golfer/settings
func GolferSettingsPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if c := r.Form.Get("country"); c != "" && config.CountryByID[c] == nil {
		http.Error(w, "Invalid country", http.StatusBadRequest)
		return
	}

	if k := r.Form.Get("keymap"); k != "default" && k != "vim" {
		http.Error(w, "Invalid keymap", http.StatusBadRequest)
		return
	}

	if t := r.Form.Get("theme"); t != "auto" && t != "dark" && t != "light" {
		http.Error(w, "Invalid theme", http.StatusBadRequest)
		return
	}

	if _, ok := zone.ByID[r.Form.Get("time_zone")]; !ok && r.Form.Get("time_zone") != "" {
		http.Error(w, "Invalid time_zone", http.StatusBadRequest)
		return
	}

	if _, err := session.Database(r).Exec(
		`UPDATE users
		    SET country = $1,
		         keymap = $2,
		    referrer_id = (SELECT id FROM users WHERE login = $3 AND id != $7),
		   show_country = $4,
		          theme = $5,
		      time_zone = $6
		  WHERE id = $7`,
		r.Form.Get("country"),
		r.Form.Get("keymap"),
		r.Form.Get("referrer"),
		r.Form.Get("show_country") == "on",
		r.Form.Get("theme"),
		r.Form.Get("time_zone"),
		session.Golfer(r).ID,
	); err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/golfer/settings", http.StatusSeeOther)
}
