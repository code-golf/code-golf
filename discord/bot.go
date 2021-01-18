package discord

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/code-golf/code-golf/hole"
	"github.com/code-golf/code-golf/lang"
)

var bot *discordgo.Session

const RECORD_ANNOUNCE_CHANNEL = "755435773096099992"
const PREFIX = "%"

func init() {
	// Stop event (on SIGINT)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		// The authentication token is kept in botAuth.txt (untracked on Git)
		token, err := ioutil.ReadFile("/go/discord/botAuth.txt")
		if err != nil {
			return
		}

		bot, err = discordgo.New("Bot " + string(token))
		if err != nil {
			fmt.Println(err)
			return
		}

		// Start event
		bot.AddHandler(func(session *discordgo.Session, event *discordgo.Ready) {
			go func() {
				fmt.Println("Discord bot is now online!")
				session.UpdateStatus(0, "Code Golf")
			}()
		})

		// Message event
		bot.AddHandler(handleMessage)

		// Start the bot!
		if bot.Open() != nil {
			fmt.Println("Error starting the Discord bot")
			return
		}

		// Wait for sigint, then disconnect
		<-sig
		bot.Close()
		fmt.Println("Discord bot disconnected")
	}()
}

// Receive code for a given hole and send a message if it breaks that hole's record
func LogNewRecord(db *sql.DB, code string, userName string, holeName string, langName string) {

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

			bot.ChannelMessageSend(RECORD_ANNOUNCE_CHANNEL,
				fmt.Sprintf("New %srecord on %s in %s by %s!\n\t%d bytes / %d chars",
					scoring, holeName, langName, userName, newBytes, newChars))
		}()
	}
}

// Respond to a message on Discord (if necessary)
func handleMessage(session *discordgo.Session, event *discordgo.MessageCreate) {

}
