package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/jraspass/code-golf/cookie"
	"github.com/julienschmidt/httprouter"
)

var clientSecret = os.Getenv("CLIENT_SECRET")

func githubAuth(code string) (int, string) {
	req, _ := http.NewRequest(
		"POST",
		"https://github.com/login/oauth/access_token",
		bytes.NewBufferString(url.Values{
			"client_id":     {"7f6709819023e9215205"},
			"client_secret": {clientSecret},
			"code":          {code},
		}.Encode()),
	)

	req.Header.Add("Accept", "application/json")

	res, _ := http.DefaultClient.Do(req)

	type accessToken struct {
		AccessToken string `json:"access_token"`
	}

	var a accessToken

	json.NewDecoder(res.Body).Decode(&a)

	res, _ = http.Get("https://api.github.com/user?access_token=" + a.AccessToken)

	type user struct {
		ID    int    `json:id`
		Login string `json:login`
	}

	var u user

	if json.NewDecoder(res.Body).Decode(&u) == nil {
		if _, err := db.Exec(
			`INSERT INTO users VALUES($1, $2)
			 ON CONFLICT(id) DO UPDATE SET login = $2`,
			u.ID, u.Login,
		); err != nil {
			panic(err)
		}
	}

	return u.ID, u.Login
}

func callback(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if userID, login := githubAuth(r.FormValue("code")); userID != 0 {
		data := strconv.Itoa(userID) + ":" + login

		w.Header()["Set-Cookie"] = []string{
			"__Host-user=" + data + ":" + cookie.Write(data) +
				";HttpOnly;Path=/;SameSite=Lax;Secure",
		}
	}

	if uri := r.FormValue("redirect_uri"); uri != "" {
		w.Header().Set("Location", uri)
	} else {
		w.Header().Set("Location", "/")
	}

	w.WriteHeader(302)
}
