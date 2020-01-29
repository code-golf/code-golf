package routes

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/code-golf/code-golf/cookie"
)

var clientSecret = os.Getenv("CLIENT_SECRET")

func callback(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	if code := query.Get("code"); code != "" {
		req, err := http.NewRequest(
			"POST",
			"https://github.com/login/oauth/access_token",
			strings.NewReader(url.Values{
				"client_id":     {"7f6709819023e9215205"},
				"client_secret": {clientSecret},
				"code":          {code},
			}.Encode()),
		)
		if err != nil {
			panic(err)
		}

		req.Header.Add("Accept", "application/json")

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()

		var accessToken struct {
			AccessToken string `json:"access_token"`
		}

		if err := json.NewDecoder(res.Body).Decode(&accessToken); err != nil {
			panic(err)
		}

		res, err = http.Get("https://api.github.com/user?access_token=" + accessToken.AccessToken)
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()

		var user struct {
			ID    int    `json:"id"`
			Login string `json:"login"`
		}

		if err := json.NewDecoder(res.Body).Decode(&user); err != nil {
			panic(err)
		}

		if _, err := db.Exec(
			`INSERT INTO users VALUES($1, $2)
				ON CONFLICT(id) DO UPDATE SET login = $2`,
			user.ID, user.Login,
		); err != nil {
			panic(err)
		}

		data := strconv.Itoa(user.ID) + ":" + user.Login

		w.Header().Set(
			"Set-Cookie",
			"__Host-user="+data+":"+cookie.Write(data)+
				";HttpOnly;Path=/;SameSite=Lax;Secure",
		)
	}

	if uri := query.Get("redirect_uri"); uri != "" {
		w.Header().Set("Location", uri)
	} else {
		w.Header().Set("Location", "/")
	}

	w.WriteHeader(302)
}
