package routes

import (
	"net/http"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/oauth"
	"github.com/code-golf/code-golf/session"
	"github.com/code-golf/code-golf/zone"
)

// POST /golfer/cancel-delete
func golferCancelDeletePOST(w http.ResponseWriter, r *http.Request) {
	if _, err := session.Database(r).Exec(
		"UPDATE users SET delete = NULL WHERE id = $1",
		session.Golfer(r).ID,
	); err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/golfer/settings", http.StatusSeeOther)
}

// POST /golfer/delete
func golferDeletePOST(w http.ResponseWriter, r *http.Request) {
	if _, err := session.Database(r).Exec(
		"UPDATE users SET delete = TIMEZONE('UTC', NOW()) + INTERVAL '7 days' WHERE id = $1",
		session.Golfer(r).ID,
	); err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/golfer/settings", http.StatusSeeOther)
}

// GET /golfer/settings
func golferSettingsGET(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Connections    []oauth.Connection
		Countries      map[string][]*config.Country
		Layouts        []string
		OAuthProviders map[string]*oauth.Config
		OAuthState     string
		Themes         []string
		TimeZones      []zone.Zone
	}{
		oauth.GetConnections(session.Database(r), session.Golfer(r).ID, false),
		config.CountryTree,
		[]string{"default", "tabs"},
		oauth.Providers,
		nonce(),
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

	render(w, r, "golfer/settings", data, "Settings")
}

// POST /golfer/settings
func golferSettingsPOST(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if c := r.Form.Get("country"); c != "" && config.CountryByID[c] == nil {
		http.Error(w, "Invalid country", http.StatusBadRequest)
		return
	}

	if k := r.Form.Get("layout"); k != "default" && k != "tabs" {
		http.Error(w, "Invalid layout", http.StatusBadRequest)
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

	tx, err := session.Database(r).BeginTx(r.Context(), nil)
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()

	userID := session.Golfer(r).ID
	if _, err := tx.Exec(
		`UPDATE users
		    SET country = $1,
				     layout = $2,
		    referrer_id = (SELECT id FROM users WHERE login = $3 AND id != $7),
		   show_country = $4,
		          theme = $5,
		      time_zone = $6
		  WHERE id = $7`,
		r.Form.Get("country"),
		r.Form.Get("layout"),
		r.Form.Get("referrer"),
		r.Form.Get("show_country") == "on",
		r.Form.Get("theme"),
		r.Form.Get("time_zone"),
		userID,
	); err != nil {
		panic(err)
	}

	// Update connections' publicness if they differ from DB.
	for _, c := range oauth.GetConnections(tx, userID, false) {
		if show := r.Form.Get("show_"+c.Connection) == "on"; show != c.Public {
			if _, err := tx.Exec(
				`UPDATE connections
				    SET public = $1
				  WHERE connection = $2 AND user_id = $3`,
				show,
				c.Connection,
				userID,
			); err != nil {
				panic(err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/golfer/settings", http.StatusSeeOther)
}
