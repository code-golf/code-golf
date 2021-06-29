package discord

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

var Session *discordgo.Session

func init() {
	/* The authentication token of the bot and the ID of the announcement channel (800680710964903946)
	should be stored in the .env file. If either of them isn't found, the bot is skipped */
	var token string
	if token = os.Getenv("DISCORD_BOT_TOKEN"); token != "" {
		go func() {
			var err error
			if Session, err = discordgo.New("Bot " + token); err != nil {
				log.Println(err)
			} else if err := Session.Open(); err != nil {
				log.Println(err)
			}
		}()
	}
}
