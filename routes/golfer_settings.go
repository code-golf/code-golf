package routes

import (
	"net/http"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/null"
	"github.com/code-golf/code-golf/oauth"
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
		Connections    []oauth.Connection
		Countries      map[string][]*config.Country
		Keymaps        []string
		OAuthProviders map[string]*oauth.Config
		OAuthState     string
		Pronouns       []string
		Themes         []string
		TimeZones      []zone.Zone
	}{
		oauth.GetConnections(session.Database(r), session.Golfer(r).ID, false),
		config.CountryTree,
		[]string{"default", "vim"},
		oauth.Providers,
		nonce(),
		[]string{"he/him", "she/her", "they/them"},
		[]string{"auto", "dark", "light"},
		zone.List(),
	}

	http.SetCookie(w, &http.Cookie{
		HttpOnly: true,
		Name:     "__Host-oauth-state",
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		Secure:   true,
		Value:    data.OAuthState,
	})

	render(w, r, "golfer/settings", data, "Settings: General")
}

// GET /golfer/settings/export-data
func golferSettingsExportDataGET(w http.ResponseWriter, r *http.Request) {
	render(w, r, "golfer/settings", nil, "Settings: Export Data")
}

// GET /golfer/settings/export-data
func golferSettingsDeleteAccountGET(w http.ResponseWriter, r *http.Request) {
	render(w, r, "golfer/settings", nil, "Settings: Delete Account")
}

// POST /golfer/settings/save
func golferSettingsSavePOST(w http.ResponseWriter, r *http.Request) {
	golfer := session.Golfer(r)
	page := r.FormValue("page")

	// If the posted value is valid, update the golfer's settings map.
	for _, setting := range config.Settings[page] {
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

	if k := r.Form.Get("keymap"); k != "default" && k != "vim" {
		http.Error(w, "Invalid keymap", http.StatusBadRequest)
		return
	}

	switch r.Form.Get("pronouns") {
	case "", "he/him", "she/her", "they/them":
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
		       keymap = $3,
		     pronouns = $4,
		  referrer_id = (SELECT id FROM users WHERE login = $5 AND id != $9),
		 show_country = $6,
		        theme = $7,
		    time_zone = $8
		  WHERE id = $9`,
		r.Form.Get("about"),
		r.Form.Get("country"),
		r.Form.Get("keymap"),
		null.NullIfZero(r.Form.Get("pronouns")),
		r.Form.Get("referrer"),
		r.Form.Get("show_country") == "on",
		r.Form.Get("theme"),
		r.Form.Get("time_zone"),
		golfer.ID,
	)

	// Update connections' publicness if they differ from DB.
	for _, c := range oauth.GetConnections(tx, golfer.ID, false) {
		if show := r.Form.Get("show_"+c.Connection) == "on"; show != c.Public {
			tx.MustExec(
				`UPDATE connections
				    SET public = $1
				  WHERE connection = $2 AND user_id = $3`,
				show,
				c.Connection,
				golfer.ID,
			)
		}
	}

	// TODO Add "flash" messages so we can show the cheevo after the redirect.
	if r.Form.Get("about") != "" {
		golfer.Earn(tx, "biohazard")
	}

	if err := tx.Commit(); err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/golfer/settings", http.StatusSeeOther)
}
