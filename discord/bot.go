package discord

import (
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/code-golf/code-golf/golfer"
	"github.com/code-golf/code-golf/hole"
	"github.com/code-golf/code-golf/lang"
)

var bot *discordgo.Session

const channelID = "800680710964903946"

func init() {
	// The authentication token of the bot should be stored in the .env file
	// If it's not found, the bot is skipped
	if token := os.Getenv("DISCORD_BOT_TOKEN"); token != "" {
		go func() {
			var err error
			if bot, err = discordgo.New("Bot " + token); err != nil {
				log.Println(err)
			} else if err := bot.Open(); err != nil {
				log.Println(err)
			}
		}()
	}
}

// LogNewRecord logs a record breaking solution in Discord.
func LogNewRecord(
	golfer *golfer.Golfer, hole hole.Hole, lang lang.Lang, scorings string,
	bytes int64, chars int64,
) {
	if bot == nil {
		return
	}

	scoring := ""
	if scorings != "byteschars" {
		scoring = scorings[:len(scorings)-1] + " "
	}

	_, err := bot.ChannelMessageSendEmbed(channelID, &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("New record on %s in %s!", hole.Name, lang.Name),
		Description: fmt.Sprintf("A new %srecord has been set by %s!", scoring, golfer.Name),
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Bytes", Inline: true, Value: fmt.Sprint(bytes)},
			{Name: "Chars", Inline: true, Value: fmt.Sprint(chars)},
		},
	})
	if err != nil {
		log.Println(err)
	}
}
