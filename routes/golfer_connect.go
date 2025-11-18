package routes

import (
	"encoding/json/v2"
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

	var user struct{ AvatarURL, Discriminator, ID, Username string }

	switch conn {
	case "discord":
		var info struct {
			Avatar        string `json:"avatar"`
			Discriminator string `json:"discriminator"`
			ID            string `json:"id"`
			Username      string `json:"username"`
		}

		if err := json.UnmarshalRead(res.Body, &info); err != nil {
			panic(err)
		}

		if info.Avatar != "" {
			user.AvatarURL = "https://cdn.discordapp.com/avatars/" +
				info.ID + "/" + info.Avatar
		}

		user.Discriminator = info.Discriminator
		user.ID = info.ID
		user.Username = info.Username
	case "gitlab":
		var info struct {
			Nickname string `json:"nickname"`
			Picture  string `json:"picture"`
			Sub      string `json:"sub"`
		}

		if err := json.UnmarshalRead(res.Body, &info); err != nil {
			panic(err)
		}

		user.AvatarURL = info.Picture
		user.ID = info.Sub
		user.Username = info.Nickname
	case "stack-overflow":
		var info struct {
			Items []struct {
				DisplayName  string `json:"display_name"`
				ProfileImage string `json:"profile_image"`
				UserID       int    `json:"user_id"`
			} `json:"items"`
		}

		if err := json.UnmarshalRead(res.Body, &info); err != nil {
			panic(err)
		}

		user.AvatarURL = info.Items[0].ProfileImage
		user.ID = strconv.Itoa(info.Items[0].UserID)
		user.Username = info.Items[0].DisplayName
	}

	session.Database(r).MustExec(
		`INSERT INTO connections (avatar_url, connection, discriminator, id, user_id, username)
		      VALUES             (        $1,         $2,            $3, $4,      $5,       $6)
		 ON CONFLICT             (connection, id)
		   DO UPDATE SET avatar_url = excluded.avatar_url,
		              discriminator = excluded.discriminator,
		                   username = excluded.username`,
		null.NullIfZero(user.AvatarURL),
		conn,
		null.NullIfZero(user.Discriminator),
		user.ID,
		session.Golfer(r).ID,
		user.Username,
	)

	http.Redirect(w, r, "/golfer/settings", http.StatusSeeOther)
}
