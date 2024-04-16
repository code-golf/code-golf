package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/code-golf/code-golf/null"
	"github.com/code-golf/code-golf/oauth"
	"github.com/code-golf/code-golf/session"
)

// GET /golfer/disconnect/{connection}
func golferDisconnectGET(w http.ResponseWriter, r *http.Request) {
	session.Database(r).MustExec(
		"DELETE FROM connections WHERE connection::text = $1 AND user_id = $2",
		param(r, "connection"), session.Golfer(r).ID,
	)

	http.Redirect(w, r, "/golfer/settings", http.StatusSeeOther)
}

// GET /golfer/connect/{connection}
func golferConnectGET(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	conn := param(r, "connection")
	config := oauth.Providers[conn]

	if code == "" || config == nil || config.Name == "GitHub" ||
		cookie(r, "__Host-oauth-state") != r.FormValue("state") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := config.Exchange(r.Context(), code)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequestWithContext(
		r.Context(), "GET", config.UserEndpoint, nil)
	if err != nil {
		panic(err)
	}

	// Stack Overflow expects the access token in the query string instead.
	if conn == "stack-overflow" {
		q := req.URL.Query()
		q.Add("access_token", token.AccessToken)
		req.URL.RawQuery = q.Encode()
	} else {
		req.Header.Add("Authorization", "Bearer "+token.AccessToken)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	var user struct{ Discriminator, ID, Username string }

	switch conn {
	case "discord":
		if err := json.NewDecoder(res.Body).Decode(&user); err != nil {
			panic(err)
		}
	case "gitlab":
		var info struct{ Nickname, Sub string }

		if err := json.NewDecoder(res.Body).Decode(&info); err != nil {
			panic(err)
		}

		user.ID = info.Sub
		user.Username = info.Nickname
	case "stack-overflow":
		var info struct {
			Items []struct {
				DisplayName string `json:"display_name"`
				UserID      int    `json:"user_id"`
			}
		}

		if err := json.NewDecoder(res.Body).Decode(&info); err != nil {
			panic(err)
		}

		user.ID = strconv.Itoa(info.Items[0].UserID)
		user.Username = info.Items[0].DisplayName
	}

	session.Database(r).MustExec(
		`INSERT INTO connections (connection, discriminator, id, user_id, username)
		      VALUES             (        $1,            $2, $3,      $4,       $5)
		 ON CONFLICT             (connection, id)
		   DO UPDATE SET discriminator = excluded.discriminator,
		                      username = excluded.username`,
		conn,
		null.New(user.Discriminator, user.Discriminator != ""),
		user.ID,
		session.Golfer(r).ID,
		user.Username,
	)

	http.Redirect(w, r, "/golfer/settings", http.StatusSeeOther)
}
