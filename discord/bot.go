package discord

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
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

// Receive code for a given hole and send a message if it breaks that hole's record
func LogNewRecord(db *sql.DB, code string, userName string, holeName string, langName string) {
	if bot != nil {
		// Save the current top records to compare against the new solution
		oldBytes, oldChars := 0, 0
		err := db.QueryRow(
			`SELECT MIN(bytes), MIN(chars) FROM code
				JOIN solutions ON code_id = id
				WHERE hole=$1 AND lang=$2;`,
			holeName, langName,
		).Scan(&oldBytes, &oldChars)

		if err == nil {
			go func() {
				newBytes, newChars := len(code), len([]rune(code))
				newByteRecord := newBytes < oldBytes
				newCharRecord := newChars < oldChars

				// Make sure the user has broken at least one of the records
				if !newByteRecord && !newCharRecord {
					return
				}

				holeName = hole.ByID[holeName].Name
				langName = lang.ByID[langName].Name

				scoring := ""
				if !newCharRecord {
					scoring = "byte "
				}
				if !newByteRecord {
					scoring = "char "
				}

				bot.ChannelMessageSend(recordAnnounceChannel,
					fmt.Sprintf("New %srecord on %s in %s by %s!\n\t%d bytes / %d chars",
						scoring, holeName, langName, userName, newBytes, newChars))
			}()
		}
	}
}

// Respond to a message on Discord (if necessary)
func handleMessage(session *discordgo.Session, event *discordgo.MessageCreate) {
	if strings.HasPrefix(event.Content, prefix) {
		// TODO
	}
}
