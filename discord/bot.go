package discord

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/code-golf/code-golf/config"
	Golfer "github.com/code-golf/code-golf/golfer"
	"github.com/code-golf/code-golf/pretty"
	"github.com/jmoiron/sqlx"
)

const minElapsedTimeToShowDate = 30 * 24 * time.Hour

var (
	bot              *discordgo.Session
	lastAnnouncement *RecAnnouncement
	mux              sync.Mutex

	// All the config keys!
	botToken      = os.Getenv("DISCORD_BOT_TOKEN")       // Caddie
	channelID     = os.Getenv("DISCORD_CHANNEL_ID")      // 🍇・sour-grapes
	guildID       = os.Getenv("DISCORD_GUILD_ID")        // Code Golf
	roleContribID = os.Getenv("DISCORD_ROLE_CONTRIB_ID") // Contributor
	roleSponsorID = os.Getenv("DISCORD_ROLE_SPONSOR_ID") // Sponsor
)

// Represents a new record announcement message
type RecAnnouncement struct {
	MessageChannelID string                `json:"messageChannelID"`
	MessageID        string                `json:"messageID"`
	Updates          [][]Golfer.RankUpdate `json:"updates"`
	Golfer           *Golfer.Golfer        `json:"-"`
	GolferID         int                   `json:"golfer"`
	Hole             *config.Hole          `json:"hole"`
	Lang             *config.Lang          `json:"lang"`
}

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
		}
	}()
}

func getUsername(id int, db *sqlx.DB) (name string) {
	err := db.Get(&name, "SELECT login FROM users WHERE id = $1", id)
	if err != nil {
		name = "unknown golfer"
		log.Println(err)
	}

	return
}

// recAnnounceToEmbed parses a recAnnouncement object and turns it into a Discord embed
func recAnnounceToEmbed(announce *RecAnnouncement, db *sqlx.DB) *discordgo.MessageEmbed {
	hole, lang, golfer := announce.Hole, announce.Lang, announce.Golfer
	imageURL := "https://code.golf/golfers/" + golfer.Name + "/avatar"
	golferURL := "https://code.golf/golfers/" + golfer.Name

	titlePrefix := "New Tied 🥇"
	isUnicorn := false

	// Creating the basic embed
	embed := &discordgo.MessageEmbed{
		URL:    "https://code.golf/rankings/holes/" + hole.ID + "/" + lang.ID + "/",
		Fields: make([]*discordgo.MessageEmbedField, 0, 2),
		Author: &discordgo.MessageEmbedAuthor{Name: golfer.Name, IconURL: imageURL, URL: golferURL},
	}

	// Now, we fill out the fields according to the updates of the announcement
	fieldValues := make(map[string]string)
	for _, pair := range announce.Updates {
		for _, update := range pair {
			if update.NewSolutionCount == 1 {
				// Once we detect a unicorn, we continue to show it when we edit messages.
				titlePrefix = "New 🦄"
				isUnicorn = true
			} else if !update.To.Joint.V && !isUnicorn {
				titlePrefix = "New 💎"
			}

			if update.OldBestStrokes.Valid && fieldValues[update.Scoring] == "" {
				fieldValues[update.Scoring] = pretty.Comma(update.OldBestStrokes.V)

				dateString := ""
				timestamp := update.OldBestSubmitted.V
				if time.Since(timestamp) > minElapsedTimeToShowDate {
					// Show the data using a locale-specific short date format.
					dateString = fmt.Sprintf("<t:%d:R>", timestamp.Unix())
				}

				// Determine the name or number of other golfers associated with the old record.
				othersString := ""
				if update.OldBestCurrentGolferCount.Valid && update.OldBestCurrentGolferCount.V > 1 {
					// Display the number of golfers, excluding the current golfer, that previously held this record.
					othersString = fmt.Sprintf("%d golfers", update.OldBestCurrentGolferCount.V)
				} else if update.OldBestCurrentGolferID.Valid && update.OldBestCurrentGolferID.V != golfer.ID {
					// Display the user name of the single golfer, excluding the current golfer, that previously held this record.
					othersString = getUsername(update.OldBestCurrentGolferID.V, db)
				}

				if othersString != "" && update.OldBestFirstGolferID.Valid && update.OldBestFirstGolferID.V == golfer.ID {
					// Report that the current golfer was the first to obtain the old record.
					othersString = fmt.Sprintf("%s, tied by %s", golfer.Name, othersString)
				}

				parenthetical := ""
				if othersString == "" {
					parenthetical = dateString
				} else if dateString == "" {
					parenthetical = othersString
				} else {
					parenthetical = fmt.Sprintf("%s by %s", dateString, othersString)
				}

				if parenthetical != "" {
					fieldValues[update.Scoring] += fmt.Sprintf(" (%s)", parenthetical)
				}
			}

			if !update.OldBestStrokes.Valid || update.To.Strokes.V < update.OldBestStrokes.V {
				if fieldValues[update.Scoring] != "" {
					fieldValues[update.Scoring] += "  →  "
				}
				fieldValues[update.Scoring] += pretty.Comma(update.To.Strokes.V)
			}

			if update.FailingStrokes.Valid && update.FailingStrokes.V <= update.To.Strokes.V {
				fieldValues[update.Scoring] += fmt.Sprintf(" (replaced failing %d)", update.FailingStrokes.V)
			}
		}
	}

	embed.Title = fmt.Sprintf("%s on %s in %s!", titlePrefix, hole.Name, lang.Name)

	if fieldValues["bytes"] == fieldValues["chars"] {
		fieldValues = map[string]string{"bytes/chars": fieldValues["bytes"]}
	}

	// Find the dominant scoring (only "chars" if there were no improvements on bytes)
	if fieldValues["bytes"] == "" && fieldValues["bytes/chars"] == "" {
		embed.URL += "chars"
		// Zero-width space to always show bytes column, to avoid confusion.
		fieldValues["bytes"] = "\u200b"
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
func AwardRoles(db *sqlx.DB) error {
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

func awardRoles(db *sqlx.DB, members map[string]any, roleID, sql string) error {
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

	imageURL := "https://code.golf/golfers/" + golfer.Name + "/avatar"
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
	golfer *Golfer.Golfer, hole *config.Hole, lang *config.Lang, updates []Golfer.RankUpdate, db *sqlx.DB,
) {
	mux.Lock()
	logNewRecord(golfer, hole, lang, updates, db)
	mux.Unlock()
}

func logNewRecord(
	golfer *Golfer.Golfer, hole *config.Hole, lang *config.Lang, updates []Golfer.RankUpdate, db *sqlx.DB,
) {
	if bot == nil {
		return
	}

	announcement := &RecAnnouncement{
		Hole:     hole,
		Lang:     lang,
		Golfer:   golfer,
		GolferID: golfer.ID,
		Updates:  [][]Golfer.RankUpdate{updates},
	}

	if lastAnnouncement == nil {
		loadLastAnnouncement(db)
	}

	if lastAnnouncement != nil {
		if channel, err := bot.Channel(channelID); err == nil {
			if channel.LastMessageID != lastAnnouncement.MessageID {
				// Discard the last announcement if another message was sent after it
				lastAnnouncement = nil
			}
		} else {
			log.Println(err)
		}
	}

	if lastAnnouncement != nil &&
		announcement.Lang.ID == lastAnnouncement.Lang.ID &&
		announcement.Hole.ID == lastAnnouncement.Hole.ID &&
		announcement.GolferID == lastAnnouncement.GolferID {
		lastAnnouncement.Golfer = golfer
		lastAnnouncement.Updates = append(lastAnnouncement.Updates, updates)
		if _, err := bot.ChannelMessageEditEmbed(
			lastAnnouncement.MessageChannelID,
			lastAnnouncement.MessageID,
			recAnnounceToEmbed(lastAnnouncement, db),
		); err == nil { // Note that we only return if the embed was edited successfully;
			saveLastAnnouncement(lastAnnouncement, db)
			return // otherwise, we'll continue forward and send it as a new message
		}
	}

	var prevMessage string
	var newMessage *discordgo.Message
	var sendErr error

	if err := db.QueryRow(
		"SELECT message FROM discord_records WHERE hole = $1 AND lang = $2",
		hole.ID, lang.ID,
	).Scan(&prevMessage); err == nil {
		newMessage, sendErr = bot.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
			Embed: recAnnounceToEmbed(announcement, db),
			Reference: &discordgo.MessageReference{
				MessageID: prevMessage,
				ChannelID: channelID,
			},
		})
	} else if errors.Is(err, sql.ErrNoRows) {
		newMessage, sendErr = bot.ChannelMessageSendEmbed(channelID, recAnnounceToEmbed(announcement, db))
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
		lastAnnouncement.MessageChannelID = newMessage.ChannelID
		lastAnnouncement.MessageID = newMessage.ID
		saveLastAnnouncement(lastAnnouncement, db)
	}
}

func loadLastAnnouncement(db *sqlx.DB) {
	var bytes []byte

	if err := db.QueryRow(
		"SELECT value FROM discord_state WHERE key = 'lastAnnouncement'",
	).Scan(&bytes); err != nil {
		log.Println(err)
		return
	}

	var announcement RecAnnouncement
	if err := json.Unmarshal(bytes, &announcement); err != nil {
		log.Println(err)
	} else {
		lastAnnouncement = &announcement
	}
}

func saveLastAnnouncement(announce *RecAnnouncement, db *sqlx.DB) {
	bytes, err := json.Marshal(announce)
	if err != nil {
		log.Println(err)
		return
	}

	if _, err := db.Exec(
		`INSERT INTO discord_state (key, value) VALUES
			('lastAnnouncement', $1)
			ON CONFLICT ON CONSTRAINT discord_state_pkey
			DO UPDATE SET value = $1`,
		bytes,
	); err != nil {
		log.Println(err)
	}
}
