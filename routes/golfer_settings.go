package routes

import (
	"net/http"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/null"
	"github.com/code-golf/code-golf/session"
	"github.com/code-golf/code-golf/zone"
)

// POST /golfer/cancel-delete
func golferCancelDeletePOST(w http.ResponseWriter, r *http.Request) {
	session.Database(r).MustExec(
		"UPDATE users SET delete = NULL WHERE id = $1",
		session.Golfer(r).ID,
	)

	http.Redirect(w, r, "/golfer/settings/delete-account", http.StatusSeeOther)
}

// POST /golfer/delete
func golferDeletePOST(w http.ResponseWriter, r *http.Request) {
	session.Database(r).MustExec(
		"UPDATE users SET delete = TIMEZONE('UTC', NOW()) + INTERVAL '7 days' WHERE id = $1",
		session.Golfer(r).ID,
	)

	http.Redirect(w, r, "/golfer/settings/delete-account", http.StatusSeeOther)
}

// GET /golfer/settings
func golferSettingsGET(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Countries map[string][]*config.Country
		Pronouns  []string
		Themes    []string
		TimeZones []zone.Zone
	}{
		config.CountryTree,
		[]string{"he/him", "he/they", "she/her", "she/they", "they/them"},
		[]string{"auto", "dark", "light"},
		zone.List(),
	}

	render(w, r, "golfer/settings", data, "Settings: General")
}

// GET /golfer/settings/banners
func golferSettingsBannersGET(w http.ResponseWriter, r *http.Request) {
	render(w, r, "golfer/settings", nil, "Settings: Banners")
}

// GET /golfer/settings/export-data
func golferSettingsExportDataGET(w http.ResponseWriter, r *http.Request) {
	render(w, r, "golfer/settings", nil, "Settings: Export Data")
}

// GET /golfer/settings/delete-account
func golferSettingsDeleteAccountGET(w http.ResponseWriter, r *http.Request) {
	render(w, r, "golfer/settings", nil, "Settings: Delete Account")
}

// POST /golfer/settings/save
func golferSettingsSavePOST(w http.ResponseWriter, r *http.Request) {
	golfer := session.Golfer(r)
	page := r.FormValue("page")

	// If the posted value is valid, update the golfer's settings map.
	for _, setting := range config.Settings[page] {
		if _, ok := golfer.Settings[page]; !ok {
			golfer.Settings[page] = map[string]any{}
		}

		golfer.Settings[page][setting.ID] =
			setting.FromFormValue(r.FormValue(setting.ID))
	}

	golfer.SaveSettings(session.Database(r))

	http.Redirect(w, r, r.FormValue("path"), http.StatusFound)
}

// POST /golfer/settings/reset
func golferSettingsResetPOST(w http.ResponseWriter, r *http.Request) {
	golfer := session.Golfer(r)

	delete(golfer.Settings, r.FormValue("page"))

	golfer.SaveSettings(session.Database(r))

	http.Redirect(w, r, r.FormValue("path"), http.StatusFound)
}

// POST /golfer/settings
func golferSettingsPOST(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if len(r.Form.Get("about")) > 255 {
		http.Error(w, "Invalid about", http.StatusBadRequest)
		return
	}

	if c := r.Form.Get("country"); c != "" && config.CountryByID[c] == nil {
		http.Error(w, "Invalid country", http.StatusBadRequest)
		return
	}

	switch r.Form.Get("pronouns") {
	case "", "he/him", "he/they", "she/her", "she/they", "they/them":
	default:
		http.Error(w, "Invalid pronouns", http.StatusBadRequest)
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

	tx := session.Database(r).MustBeginTx(r.Context(), nil)
	defer tx.Rollback()

	golfer := session.Golfer(r)
	tx.MustExec(
		`UPDATE users
		    SET about = $1,
		      country = $2,
		     pronouns = $3,
		  referrer_id = (SELECT id FROM users WHERE name = $4 AND id != $8),
		 show_country = $5,
		        theme = $6,
		    time_zone = $7
		  WHERE id = $8`,
		r.Form.Get("about"),
		null.NullIfZero(r.Form.Get("country")),
		null.NullIfZero(r.Form.Get("pronouns")),
		r.Form.Get("referrer"),
		r.Form.Get("show_country") == "on",
		r.Form.Get("theme"),
		r.Form.Get("time_zone"),
		golfer.ID,
	)

	// TODO Add "flash" messages so we can show the cheevo after the redirect.
	if r.Form.Get("about") != "" {
		golfer.Earn(tx, "biohazard")
	}

	if err := tx.Commit(); err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/golfer/settings", http.StatusSeeOther)
}
