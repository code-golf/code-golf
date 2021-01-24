package discord

import (
	"fmt"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/code-golf/code-golf/golfer"
	"github.com/code-golf/code-golf/hole"
	"github.com/code-golf/code-golf/lang"
)

var bot *discordgo.Session

const recordAnnounceChannel = "800680710964903946"
const prefix = "%"

func init() {
	// The authentication token of the bot should be stored in the .env file
	// If it's not found, the bot is skipped
	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token != "" {
		go func() {
			var err error
			bot, err = discordgo.New("Bot " + string(token))
			if err != nil {
				fmt.Println("Couldn't start bot:", err)
				return
			}

			// Start event
			bot.AddHandler(func(session *discordgo.Session, event *discordgo.Ready) {
				go func() {
					fmt.Println("Discord bot is now online!")
				}()
			})

			// Message event
			bot.AddHandler(handleMessage)

			// Start the bot!
			if bot.Open() != nil {
				fmt.Println("Error starting the Discord bot")
				return
			}
		}()
	}
}

// Log a record breaking solution in Discord
func LogNewRecord(golfer *golfer.Golfer, holeName string, langName string, scorings string, bytes int64, chars int64) {
	if bot != nil {
		go func() {
			holeName = hole.ByID[holeName].Name
			langName = lang.ByID[langName].Name

			scoring := ""
			if scorings != "byteschars" {
				scoring = scorings[:len(scorings)-1] + " "
			}

			embed := &discordgo.MessageEmbed{
				Title:       fmt.Sprintf("New record on %s in %s!", holeName, langName),
				Description: fmt.Sprintf("A new %srecord has been set by %s!", scoring, golfer.Name),
				Fields: []*discordgo.MessageEmbedField{
					{Name: "Bytes", Inline: true, Value: fmt.Sprint(bytes)},
					{Name: "Chars", Inline: true, Value: fmt.Sprint(chars)},
				},
			}

			bot.ChannelMessageSendEmbed(recordAnnounceChannel, embed)
		}()
	}
}

// Respond to a message on Discord (if necessary)
func handleMessage(session *discordgo.Session, event *discordgo.MessageCreate) {
	if strings.HasPrefix(event.Content, prefix) {
		// TODO
	}
}
