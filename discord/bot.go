package discord

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/code-golf/code-golf/config"
	Golfer "github.com/code-golf/code-golf/golfer"
	"github.com/code-golf/code-golf/pretty"
)

var (
	bot *discordgo.Session

	// All the config keys!
	botToken      = os.Getenv("DISCORD_BOT_TOKEN")       // Caddie
	channelID     = os.Getenv("DISCORD_CHANNEL_ID")      // 🍇・sour-grapes
	guildID       = os.Getenv("DISCORD_GUILD_ID")        // Code Golf
	roleContribID = os.Getenv("DISCORD_ROLE_CONTRIB_ID") // Contributor
	roleSponsorID = os.Getenv("DISCORD_ROLE_SPONSOR_ID") // Sponsor
)

// Represents a new record announcement message
type RecAnnouncement struct {
	Message *discordgo.Message
	Updates [][]Golfer.RankUpdate
	Golfer  *Golfer.Golfer
	Hole    *config.Hole
	Lang    *config.Lang
}

var lastAnnouncement *RecAnnouncement

func init() {
	// Ensure we have all our config.
	switch "" {
	case botToken, channelID, guildID, roleContribID, roleSponsorID:
		return
	}

	// Connect to Discord off the main thread.
	go func() {
		var err error
		if bot, err = discordgo.New("Bot " + botToken); err != nil {
			log.Println(err)
		} else if err = bot.Open(); err != nil {
			log.Println(err)
			bot = nil
		} else {
			bot.AddHandler(handleMessage)
		}
	}()
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
			if update.Beat.Valid && fieldValues[update.Scoring] == "" {
				fieldValues[update.Scoring] = pretty.Comma(int(update.Beat.Int64))
			}
			if fieldValues[update.Scoring] != "" {
				fieldValues[update.Scoring] += "  →  "
			}
			fieldValues[update.Scoring] += pretty.Comma(int(update.To.Strokes.Int64))
		}
	}

	if fieldValues["bytes"] == fieldValues["chars"] {
		fieldValues = map[string]string{"bytes/chars": fieldValues["bytes"]}
	}

	// Find the dominant scoring (only "chars" if there were no improvements on bytes)
	if fieldValues["bytes"] == "" && fieldValues["bytes/chars"] == "" {
		embed.URL += "chars"
		fieldValues["bytes"] = "​" // Display the bytes column in any case, to avoid confusion
	} else {
		embed.URL += "bytes"
	}

	// We iterate over the scorings rather than the map itself so that the order will be guaranteed
	for _, scoring := range []string{"bytes", "chars", "bytes/chars"} {
		if fieldValues[scoring] != "" {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   pretty.Title(scoring),
				Value:  fieldValues[scoring],
				Inline: true,
			})
		}
	}

	return embed
}

// AwardRoles awards Discord roles based on cheevos etc.
func AwardRoles(db *sql.DB) error {
	if bot == nil {
		return nil
	}

	// Make maps of members with contributor or sponsor roles. TODO Paginate.
	var (
		contributors = map[string]any{}
		sponsors     = map[string]any{}
	)
	members, err := bot.GuildMembers(guildID, "", 1000)
	if err != nil {
		panic(err)
	}

	for _, member := range members {
		for _, role := range member.Roles {
			switch role {
			case roleContribID:
				contributors[member.User.ID] = nil
			case roleSponsorID:
				sponsors[member.User.ID] = nil
			}
		}
	}

	if err := awardRoles(
		db, contributors, roleContribID,
		`SELECT id
		   FROM connections JOIN trophies USING(user_id)
		  WHERE connection = 'discord' AND trophy = 'patches-welcome'`,
	); err != nil {
		return err
	}

	if err := awardRoles(
		db, sponsors, roleSponsorID,
		`SELECT c.id
		   FROM connections c JOIN users u ON user_id = u.id
		  WHERE connection = 'discord' AND sponsor`,
	); err != nil {
		return err
	}

	return nil
}

func awardRoles(db *sql.DB, members map[string]any, roleID, sql string) error {
	rows, err := db.Query(sql)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			return err
		}

		if _, ok := members[userID]; ok {
			delete(members, userID)
		} else if err := bot.GuildMemberRoleAdd(guildID, userID, roleID); err != nil {
			log.Println(err)
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	// Remove any stale roles.
	for userID := range members {
		if err := bot.GuildMemberRoleRemove(guildID, userID, roleID); err != nil {
			log.Println(err)
		}
	}

	return nil
}

func LogFailedRejudge(golfer *Golfer.Golfer, hole *config.Hole, lang *config.Lang, scoring string) {
	if bot == nil {
		return
	}

	imageURL := "https://avatars.githubusercontent.com/" + golfer.Name
	golferURL := "https://code.golf/golfers/" + golfer.Name

	if _, err := bot.ChannelMessageSendEmbed(channelID, &discordgo.MessageEmbed{
		Title:  fmt.Sprintf("%s in %s failed rejudge!", hole.Name, lang.Name),
		URL:    "https://code.golf/rankings/holes/" + hole.ID + "/" + lang.ID + "/" + scoring,
		Author: &discordgo.MessageEmbedAuthor{Name: golfer.Name, IconURL: imageURL, URL: golferURL},
	}); err != nil {
		log.Println(err)
	}
}

// LogNewRecord logs a record breaking solution in Discord.
func LogNewRecord(
	golfer *Golfer.Golfer, hole *config.Hole, lang *config.Lang, updates []Golfer.RankUpdate, db *sql.DB,
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
