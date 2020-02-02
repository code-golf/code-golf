package routes

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/code-golf/code-golf/cookie"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var config = oauth2.Config{
	ClientID:     "7f6709819023e9215205",
	ClientSecret: os.Getenv("CLIENT_SECRET"),
	Endpoint:     github.Endpoint,
}

// Callback serves GET /callback
func Callback(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("code") == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := config.Exchange(r.Context(), r.FormValue("code"))
	if err != nil {
		panic(err)
	}

	res, err := http.Get("https://api.github.com/user?access_token=" + token.AccessToken)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	var user struct {
		ID    int
		Login string
	}

	if err := json.NewDecoder(res.Body).Decode(&user); err != nil {
		panic(err)
	}

	if _, err := db.Exec(
		`INSERT INTO users VALUES($1, $2)
		     ON CONFLICT(id) DO UPDATE SET login = excluded.login`,
		user.ID, user.Login,
	); err != nil {
		panic(err)
	}

	data := strconv.Itoa(user.ID) + ":" + user.Login

	http.SetCookie(w, &http.Cookie{
		HttpOnly: true,
		Name:     "__Host-user",
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		Secure:   true,
		Value:    data + ":" + cookie.Write(data),
	})

	uri := r.FormValue("redirect_uri")
	if uri == "" {
		uri = "/"
	}

	http.Redirect(w, r, uri, http.StatusSeeOther)
}
