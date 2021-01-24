package discord

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/code-golf/code-golf/golfer"
	"github.com/code-golf/code-golf/hole"
	"github.com/code-golf/code-golf/lang"
)

var bot *discordgo.Session

var channelID string

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
			}
		}()
	}
}

// LogNewRecord logs a record breaking solution in Discord.
func LogNewRecord(
	golfer *golfer.Golfer, hole hole.Hole, lang lang.Lang, updates []golfer.RankUpdate,
) {
	if bot == nil {
		return
	}

	imageURL := "https://avatars.githubusercontent.com/" + golfer.Name
	golferURL := "https://code.golf/golfers/" + golfer.Name

	embed := &discordgo.MessageEmbed{
		Title: fmt.Sprintf("New :first_place: on %s in %s!",
			hole.Name,
			lang.Name,
		),
		URL:    "https://code.golf/scores/" + hole.ID + "/" + lang.ID + "/" + updates[0].Scoring,
		Fields: make([]*discordgo.MessageEmbedField, 0, 2),
		Author: &discordgo.MessageEmbedAuthor{Name: golfer.Name, IconURL: imageURL, URL: golferURL},
	}

	// Add in the scorings (as necessary)
	for _, update := range updates {
		improveString := fmt.Sprint(update.To.Strokes.Int64)
		if update.From.Strokes.Valid {
			improveString = fmt.Sprint(update.From.Strokes.Int64) + "  →  " + improveString
		}
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   strings.Title(update.Scoring),
			Value:  improveString,
			Inline: true,
		})
	}

	if _, err := bot.ChannelMessageSendEmbed(channelID, embed); err != nil {
		log.Println(err)
	}
}
