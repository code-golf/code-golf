package oauth

import (
	"net/url"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/code-golf/code-golf/db"
	"github.com/code-golf/code-golf/null"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
)

type Config struct {
	oauth2.Config
	Name, UserEndpoint string
}

type Connection struct {
	Connection, Username string
	Discriminator        null.Int
	ID                   int
	Public               bool
}

var Providers = map[string]*Config{
	// https://discord.com/developers/applications
	// https://discord.com/developers/docs/resources/user
	"discord": {
		Name:         "Discord",
		UserEndpoint: discordgo.EndpointUser("@me"),
		Config: oauth2.Config{
			Endpoint: endpoints.Discord,
			Scopes:   []string{"identify"},
		},
	},

	// https://github.com/settings/developers
	// https://docs.github.com/en/rest/users/users
	"github": {
		Name:   "GitHub",
		Config: oauth2.Config{Endpoint: endpoints.GitHub},
	},

	// https://gitlab.com/-/profile/applications
	// https://docs.gitlab.com/integration/openid_connect_provider/
	"gitlab": {
		Name:         "GitLab",
		UserEndpoint: "https://gitlab.com/oauth/userinfo",
		Config: oauth2.Config{
			Endpoint: endpoints.GitLab,
			Scopes:   []string{"openid"},
		},
	},

	// https://stackapps.com/apps/oauth
	// https://api.stackexchange.com/docs/me
	"stack-overflow": {
		Name:         "Stack Overflow",
		Config:       oauth2.Config{Endpoint: endpoints.StackOverflow},
		UserEndpoint: "https://api.stackexchange.com/me?site=stackoverflow",
	},
}

func init() {
	host := "code.golf"
	if _, dev := os.LookupEnv("DEV"); dev {
		host = "localhost"
	}

	for id, config := range Providers {
		prefix := strings.ReplaceAll(strings.ToUpper(id), "-", "_")

		config.ClientID = os.Getenv(prefix + "_CLIENT_ID")
		config.ClientSecret = os.Getenv(prefix + "_CLIENT_SECRET")
		config.RedirectURL = "https://" + host + "/golfer/connect/" + id

		// Add a key to UserEndpoint if we have one.
		if key := os.Getenv(prefix + "_KEY"); key != "" {
			u, err := url.Parse(config.UserEndpoint)
			if err != nil {
				panic(err)
			}

			q := u.Query()
			q.Set("key", key)
			u.RawQuery = q.Encode()
			config.UserEndpoint = u.String()
		}
	}
}

func GetConnections(db db.Queryable, golferID int, onlyPublic bool) (c []Connection) {
	if err := db.Select(
		&c,
		` SELECT connection, discriminator, id, public, username
		    FROM connections
		   WHERE user_id = $1 AND public IN (true, $2)
		ORDER BY connection`,
		golferID,
		onlyPublic,
	); err != nil {
		panic(err)
	}

	return
}
