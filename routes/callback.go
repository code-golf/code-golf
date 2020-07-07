package routes

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/code-golf/code-golf/session"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var config = oauth2.Config{
	ClientID:     "7f6709819023e9215205",
	ClientSecret: os.Getenv("CLIENT_SECRET"),
	Endpoint:     github.Endpoint,
}

var hmacKey []byte

func init() {
	var err error
	if hmacKey, err = base64.RawURLEncoding.DecodeString(os.Getenv("HMAC_KEY")); err != nil {
		panic(err)
	}
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

	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		panic(err)
	}

	req.Header.Add("Authorization", "Bearer "+token.AccessToken)

	res, err := http.DefaultClient.Do(req)
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

	if _, err := session.Database(r).Exec(
		`INSERT INTO users VALUES($1, $2)
		     ON CONFLICT(id) DO UPDATE SET login = excluded.login`,
		user.ID, user.Login,
	); err != nil {
		panic(err)
	}

	data := strconv.Itoa(user.ID) + ":" + user.Login

	mac := hmac.New(sha256.New, hmacKey)
	mac.Write([]byte(data))

	http.SetCookie(w, &http.Cookie{
		HttpOnly: true,
		Name:     "__Host-user",
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		Secure:   true,
		Value:    data + ":" + base64.RawURLEncoding.EncodeToString(mac.Sum(nil)),
	})

	uri := r.FormValue("redirect_uri")
	if uri == "" {
		uri = "/"
	}

	http.Redirect(w, r, uri, http.StatusSeeOther)
}
