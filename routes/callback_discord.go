package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/code-golf/code-golf/session"
	"golang.org/x/oauth2"
)

var discordConfig = oauth2.Config{
	ClientID:     os.Getenv("DISCORD_CLIENT_ID"),
	ClientSecret: os.Getenv("DISCORD_CLIENT_SECRET"),
	Endpoint: oauth2.Endpoint{
		AuthURL:   discordgo.EndpointOauth2 + "authorize",
		TokenURL:  discordgo.EndpointOauth2 + "token",
		AuthStyle: oauth2.AuthStyleInParams,
	},
}

// CallbackDiscord serves GET /callback/discord
func CallbackDiscord(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("code") == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	golfer := session.Golfer(r)
	if golfer == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	discordConfig.RedirectURL = "https://" + r.Host + "/callback/discord"

	token, err := discordConfig.Exchange(r.Context(), r.FormValue("code"))
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequestWithContext(
		r.Context(), "GET", discordgo.EndpointUser("@me"), nil)
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
		ID string
	}

	if err := json.NewDecoder(res.Body).Decode(&user); err != nil {
		panic(err)
	}

	fmt.Println(user)

	if _, err := session.Database(r).Exec(
		`UPDATE users SET discord = $1`,
		user.ID,
	); err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/golfer/settings", http.StatusSeeOther)
}
