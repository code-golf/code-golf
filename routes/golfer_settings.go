package routes

import (
	"net/http"
	"time"

	"github.com/code-golf/code-golf/session"
	"github.com/code-golf/code-golf/zone"
)

// GolferSettings serves GET /golfer/settings
func GolferSettings(w http.ResponseWriter, r *http.Request) {
	render(w, r, "golfer/settings", "Settings", struct {
		TimeZones []zone.Zone
	}{zone.List()})
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
