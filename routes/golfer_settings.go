package routes

import (
	"net/http"
	"time"

	"github.com/code-golf/code-golf/country"
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
	render(w, r, "golfer/settings", "Settings", struct {
		Countries map[string][]*country.Country
		Keymaps   []string
		TimeZones []zone.Zone
	}{country.Tree, []string{"default", "vim"}, zone.List()})
}

// GolferSettingsPost serves POST /golfer/settings
func GolferSettingsPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if _, ok := country.ByID[r.Form.Get("country")]; !ok && r.Form.Get("country") != "" {
		http.Error(w, "Invalid country", http.StatusBadRequest)
		return
	}

	if r.Form.Get("keymap") != "default" && r.Form.Get("keymap") != "vim" {
		http.Error(w, "Invalid keymap", http.StatusBadRequest)
		return
	}

	timeZone, _ := time.LoadLocation(r.Form.Get("time_zone"))
	if timeZone == nil || timeZone == time.Local {
		http.Error(w, "Invalid time_zone", http.StatusBadRequest)
		return
	}

	if _, err := session.Database(r).Exec(
		"UPDATE users SET keymap = $1, time_zone = $2 WHERE id = $3",
		r.Form.Get("keymap"),
		timeZone.String(),
		session.Golfer(r).ID,
	); err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/golfer/settings", http.StatusSeeOther)
}
