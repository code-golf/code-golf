package routes

import (
	"crypto/rand"
	"net/http"

	"github.com/code-golf/code-golf/oauth"
	"github.com/code-golf/code-golf/session"
)

// GET /golfer/settings/connections
func golferSettingsConnectionsGET(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Connections    []oauth.Connection
		OAuthProviders map[string]*oauth.Config
		OAuthState     string
	}{
		oauth.GetConnections(session.Database(r), session.Golfer(r).ID, false),
		oauth.Providers,
		rand.Text(),
	}

	http.SetCookie(w, &http.Cookie{
		HttpOnly: true,
		Name:     "__Host-oauth-state",
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		Secure:   true,
		Value:    data.OAuthState,
	})

	render(w, r, "golfer/settings", data, "Settings: Connections")
}

// POST /golfer/settings/connections
func golferSettingsConnectionsPOST(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	golfer := session.Golfer(r)

	tx := session.Database(r).MustBeginTx(r.Context(), nil)
	defer tx.Rollback()

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

	if err := tx.Commit(); err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/golfer/settings/connections", http.StatusSeeOther)
}
