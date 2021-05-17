package discord

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	Golfer "github.com/code-golf/code-golf/golfer"
	"github.com/code-golf/code-golf/hole"
	"github.com/code-golf/code-golf/lang"
	"github.com/code-golf/code-golf/pretty"
)

var (
	bot       *discordgo.Session
	channelID string
)

// Represents a new record announcement message
type RecAnnouncement struct {
	Message *discordgo.Message
	Updates [][]Golfer.RankUpdate
	Golfer  *Golfer.Golfer
	Hole    hole.Hole
	Lang    lang.Lang
}

var lastAnnouncement *RecAnnouncement

func init() {
	/* The authentication token of the bot and the ID of the announcement channel (800680710964903946)
	should be stored in the .env file. If either of them isn't found, the bot is skipped */
	var token string
	if token, channelID = os.Getenv("DISCORD_BOT_TOKEN"), os.Getenv("DISCORD_BOT_CHANNEL"); token != "" && channelID != "" {
		go func() {
			var err error
			if bot, err = discordgo.New("Bot " + token); err != nil {
				log.Println(err)
			} else if err := bot.Open(); err != nil {
				log.Println(err)
			} else {
				bot.AddHandler(handleMessage)
			}
		}()
	}
}

// recAnnounceToEmbed parses a recAnnouncement object and turns it into a Discord embed
func recAnnounceToEmbed(announce *RecAnnouncement) *discordgo.MessageEmbed {
	hole, lang, golfer := announce.Hole, announce.Lang, announce.Golfer
	imageURL := "https://avatars.githubusercontent.com/" + golfer.Name
	golferURL := "https://code.golf/golfers/" + golfer.Name

	// Creating the basic embed
	embed := &discordgo.MessageEmbed{
		Title:  fmt.Sprintf("New 🥇 on %s in %s!", hole.Name, lang.Name),
		URL:    "https://code.golf/rankings/holes/" + hole.ID + "/" + lang.ID + "/",
		Fields: make([]*discordgo.MessageEmbedField, 0, 2),
		Author: &discordgo.MessageEmbedAuthor{Name: golfer.Name, IconURL: imageURL, URL: golferURL},
	}

	// Now, we fill out the fields according to the updates of the announcement
	fieldValues := make(map[string]string)
	for _, pair := range announce.Updates {
		for _, update := range pair {
			if update.Beat.Valid {
				if fieldValues[update.Scoring] == "" {
					fieldValues[update.Scoring] = pretty.Comma(int(update.Beat.Int64))
				}
				fieldValues[update.Scoring] += "  →  "
			}
			fieldValues[update.Scoring] += pretty.Comma(int(update.To.Strokes.Int64))
		}
	}

	if fieldValues["bytes"] == fieldValues["chars"] {
		fieldValues = map[string]string{"bytes/chars": fieldValues["bytes"]}
	}
	// We iterate over the scorings rather than the map itself so that the order will be guaranteed
	for _, scoring := range []string{"bytes", "chars", "bytes/chars"} {
		if fieldValues[scoring] != "" {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   strings.Title(scoring),
				Value:  fieldValues[scoring],
				Inline: true,
			})
		}
	}

	// Find the dominant scoring (only "chars" if there were no improvements on bytes)
	if fieldValues["bytes"] == "" {
		embed.URL += "chars"
	} else {
		embed.URL += "bytes"
	}

	return embed
}

// LogNewRecord logs a record breaking solution in Discord.
func LogNewRecord(
	golfer *Golfer.Golfer, hole hole.Hole, lang lang.Lang, updates []Golfer.RankUpdate, db *sql.DB,
) {
	if bot == nil {
		return
	}

	announcement := &RecAnnouncement{
		Hole:    hole,
		Lang:    lang,
		Golfer:  golfer,
		Updates: [][]Golfer.RankUpdate{updates},
	}

	if lastAnnouncement != nil &&
		announcement.Lang.ID == lastAnnouncement.Lang.ID &&
		announcement.Hole.ID == lastAnnouncement.Hole.ID &&
		announcement.Golfer.ID == lastAnnouncement.Golfer.ID {
		lastAnnouncement.Updates = append(lastAnnouncement.Updates, updates)
		if _, err := bot.ChannelMessageEditEmbed(
			lastAnnouncement.Message.ChannelID,
			lastAnnouncement.Message.ID,
			recAnnounceToEmbed(lastAnnouncement),
		); err == nil { // Note that we only return if the embed was edited successfully;
			return // otherwise, we'll continue forward and send it as a new message
		}
	}

	var prevMessage string
	var newMessage *discordgo.Message
	var sendErr error

	if err := db.QueryRow(
		`SELECT message FROM discord_records WHERE hole = $1 AND lang = $2`,
		hole.ID, lang.ID,
	).Scan(&prevMessage); err == nil {
		newMessage, sendErr = bot.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
			Embed: recAnnounceToEmbed(announcement),
			Reference: &discordgo.MessageReference{
				MessageID: prevMessage,
				ChannelID: channelID,
			},
		})
	} else if errors.Is(err, sql.ErrNoRows) {
		newMessage, sendErr = bot.ChannelMessageSendEmbed(channelID, recAnnounceToEmbed(announcement))
	} else {
		log.Println(err)
	}

	if _, err := db.Exec(
		`INSERT INTO discord_records (hole, lang, message) VALUES
			($1, $2, $3)
			ON CONFLICT ON CONSTRAINT discord_records_pkey
			DO UPDATE SET message = $3`,
		hole.ID, lang.ID, newMessage.ID,
	); err != nil {
		log.Println(err)
	}

	if sendErr != nil {
		log.Println(sendErr)
	} else {
		lastAnnouncement = announcement
		lastAnnouncement.Message = newMessage
	}
}

// handleMessage handles a message received by the bot
func handleMessage(session *discordgo.Session, event *discordgo.MessageCreate) {
	if event.Author.Bot {
		return
	}

	// Discard the last announcement if another message was sent after it
	if event.ChannelID == channelID {
		lastAnnouncement = nil
	}
}
