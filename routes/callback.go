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
	ClientID:     "7f6709819023e9215205",
	ClientSecret: os.Getenv("CLIENT_SECRET"),
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

	// Only replace NULL countries and time_zones, never user chosen ones.
	if err := session.Database(r).QueryRow(
		`WITH golfer AS (
		    INSERT INTO users (id, login, country, time_zone)
		         VALUES       ($1,    $2,      $3,        $4)
		    ON CONFLICT       (id)
		  DO UPDATE SET login = excluded.login,
		              country = COALESCE(users.country,   excluded.country),
		            time_zone = COALESCE(users.time_zone, excluded.time_zone)
		      RETURNING id
		) INSERT INTO sessions (user_id) SELECT * FROM golfer RETURNING id`,
		user.ID, user.Login, country, timeZone,
	).Scan(&cookie.Value); err != nil {
		panic(err)
	}

	http.SetCookie(w, &cookie)

	uri := r.FormValue("redirect_uri")
	if uri == "" {
		uri = "/"
	}

	http.Redirect(w, r, uri, http.StatusSeeOther)
}
