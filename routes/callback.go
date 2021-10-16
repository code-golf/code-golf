package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/code-golf/code-golf/session"
	"github.com/code-golf/code-golf/zone"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var oauthConfig = oauth2.Config{
	ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
	ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
	Endpoint:     github.Endpoint,
}

// /callback/dev exists because GitHub doesn't support multiple URLs.

// CallbackDev serves GET /callback/dev
func CallbackDev(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://localhost/callback?"+r.URL.RawQuery, http.StatusSeeOther)
}

// Callback serves GET /callback
func Callback(w http.ResponseWriter, r *http.Request) {
	var user struct {
		ID    int
		Login string
	}

	cookie := http.Cookie{
		HttpOnly: true,
		Name:     "__Host-session",
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		Secure:   true,
	}

	var country, timeZone sql.NullString

	if tz, _ := time.LoadLocation(r.FormValue("time_zone")); tz != nil {
		country.String, country.Valid = zone.Country[tz.String()]

		timeZone = sql.NullString{
			String: tz.String(),
			Valid:  tz != time.Local && tz != time.UTC,
		}
	}

	// In dev mode, the username is selected by the "username" parameter
	if _, dev := os.LookupEnv("DEV"); dev && oauthConfig.ClientSecret == "" {
		user.Login = r.FormValue("username")
		if user.Login == "" {
			user.Login = "JRaspass"
		}

		if err := session.Database(r).QueryRow(
			`SELECT COALESCE((SELECT id FROM users WHERE login = $1), COUNT(*) + 1) FROM users`,
			user.Login,
		).Scan(&user.ID); err != nil {
			panic(err)
		}
	} else {
		if r.FormValue("code") == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		token, err := oauthConfig.Exchange(r.Context(), r.FormValue("code"))
		if err != nil {
			panic(err)
		}

		req, err := http.NewRequestWithContext(
			r.Context(), "GET", "https://api.github.com/user", nil)
		if err != nil {
			panic(err)
		}

		req.Header.Add("Authorization", "Bearer "+token.AccessToken)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()

		if err := json.NewDecoder(res.Body).Decode(&user); err != nil {
			panic(err)
		}
	}

	tx, err := session.Database(r).BeginTx(r.Context(), nil)
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()

	// Create or update a user. For now user ID == GitHub user ID.
	// This'll need to change for true multi connection OAuth support.
	if _, err := tx.Exec(
		`INSERT INTO users (id, login, country, time_zone)
		      VALUES       ($1,    $2,      $3,        $4)
		 ON CONFLICT       (id)
		   DO UPDATE SET login = excluded.login,
		     country = COALESCE(users.country,   excluded.country),
		   time_zone = COALESCE(users.time_zone, excluded.time_zone)`,
		user.ID, user.Login, country, timeZone,
	); err != nil {
		panic(err)
	}

	// Create or update a connection. For now user ID == GitHub user ID.
	// This'll need to change for true multi connection OAuth support.
	if _, err := tx.Exec(
		`INSERT INTO connections (connection, id, user_id, username)
		      VALUES             (  'github', $1,      $2,       $3)
		 ON CONFLICT             (connection, id)
		   DO UPDATE SET username = excluded.username`,
		user.ID, user.ID, user.Login,
	); err != nil {
		panic(err)
	}

	// Create a session, write cookie value.
	if err := tx.QueryRow(
		"INSERT INTO sessions (user_id) VALUES ($1) RETURNING id", user.ID,
	).Scan(&cookie.Value); err != nil {
		panic(err)
	}

	if err := tx.Commit(); err != nil {
		panic(err)
	}

	http.SetCookie(w, &cookie)

	uri := r.FormValue("redirect_uri")
	if uri == "" {
		uri = "/"
	}

	http.Redirect(w, r, uri, http.StatusSeeOther)
}
