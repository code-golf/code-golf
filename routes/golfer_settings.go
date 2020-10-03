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
		TimeZones []zone.Zone
	}{country.Tree, zone.List()})
}

// GolferSettingsPost serves POST /golfer/settings
func GolferSettingsPost(w http.ResponseWriter, r *http.Request) {
	// TODO Send a 400 if the time zone isn't valid?
	timeZone, _ := time.LoadLocation(r.FormValue("time_zone"))
	if timeZone != nil && timeZone != time.Local {
		if _, err := session.Database(r).Exec(
			"UPDATE users SET time_zone = $1 WHERE id = $2",
			timeZone.String(), session.Golfer(r).ID,
		); err != nil {
			panic(err)
		}
	}

	http.Redirect(w, r, "/golfer/settings", http.StatusSeeOther)
}
