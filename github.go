package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/url"
	"os"

	_ "github.com/lib/pq"
)

var clientSecret = os.Getenv("CLIENT_SECRET")
var dbDSN = os.Getenv("DB_DSN")

type user struct {
	ID    int    `json:id`
	Login string `json:login`
}

func githubAuth(code string) user {
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

	var u user

	if json.NewDecoder(res.Body).Decode(&u) == nil {
		if db, err := sql.Open("postgres", dbDSN); err != nil {
			panic(err)
		} else if _, err := db.Exec(
			"INSERT INTO users VALUES($1, $2) ON CONFLICT(id) DO UPDATE SET login = $2",
			u.ID, u.Login,
		); err != nil {
			panic(err)
		}
	}

	return u
}
