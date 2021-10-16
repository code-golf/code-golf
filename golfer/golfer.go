package golfer

import (
	"database/sql"
	"encoding/json"
	"errors"
	"sort"
	"time"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/oauth"
	"gopkg.in/guregu/null.v4"
)

type FailingSolutions []struct{ Hole, Lang string }

func (f *FailingSolutions) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), f)
}

type Golfer struct {
	Admin, ShowCountry                     bool
	Cheevos                                []string
	Country, Keymap, Name, Referrer, Theme string
	Delete                                 sql.NullTime
	FailingSolutions                       FailingSolutions
	ID                                     int
	TimeZone                               *time.Location
}

// Earn the given cheevo, no-op if already earnt.
func (g *Golfer) Earn(db *sql.DB, cheevoID string) (earned *config.Cheevo) {
	if res, err := db.Exec(
		"INSERT INTO trophies VALUES (DEFAULT, $1, $2) ON CONFLICT DO NOTHING",
		g.ID,
		cheevoID,
	); err != nil {
		panic(err)
	} else if rowsAffected, _ := res.RowsAffected(); rowsAffected == 1 {
		earned = config.CheevoByID[cheevoID]
	}

	// Update g.Cheevos if necessary.
	if i := sort.SearchStrings(
		g.Cheevos, cheevoID,
	); i == len(g.Cheevos) || g.Cheevos[i] != cheevoID {
		g.Cheevos = append(g.Cheevos, "")
		copy(g.Cheevos[i+1:], g.Cheevos[i:])
		g.Cheevos[i] = cheevoID
	}

	return
}

// Earnt returns whether the golfer has that cheevo.
func (g *Golfer) Earnt(cheevoID string) bool {
	i := sort.SearchStrings(g.Cheevos, cheevoID)
	return i < len(g.Cheevos) && g.Cheevos[i] == cheevoID
}

type Connection struct {
	Connection, Username string
	Discriminator        sql.NullInt16
	ID                   int
	Public               bool
}

type GolferInfo struct {
	Golfer

	Sponsor bool

	// Overall points
	BytesPoints, CharsPoints int

	// Count of medals
	Diamond, Gold, Silver, Bronze int

	// Count of cheevos/holes/langs done
	Cheevos, Holes, Langs int

	// Count of cheevos/holes/langs available
	CheevosTotal, HolesTotal, LangsTotal int

	// Start date
	TeedOff time.Time

	Connections []Connection
}

// MissingConnections returns which connections are missing.
func (g *GolferInfo) MissingConnections() map[string]*oauth.Config {
	missing := map[string]*oauth.Config{}

	for id, config := range oauth.Connections {
		i := sort.Search(len(g.Connections),
			func(i int) bool { return g.Connections[i].Connection >= id })

		if i >= len(g.Connections) || g.Connections[i].Connection != id {
			missing[id] = config
		}
	}

	return missing
}

type RankUpdate struct {
	Scoring  string
	From, To struct {
		Joint         null.Bool
		Rank, Strokes null.Int
	}
	Beat null.Int
}

func GetInfo(db *sql.DB, name string) *GolferInfo {
	info := GolferInfo{
		CheevosTotal: len(config.CheevoList),
		HolesTotal:   len(config.HoleList),
		LangsTotal:   len(config.LangList),
	}

	if err := db.QueryRow(
		`WITH medals AS (
		   SELECT user_id,
		          COUNT(*) FILTER (WHERE medal = 'diamond') diamond,
		          COUNT(*) FILTER (WHERE medal = 'gold'   ) gold,
		          COUNT(*) FILTER (WHERE medal = 'silver' ) silver,
		          COUNT(*) FILTER (WHERE medal = 'bronze' ) bronze
		     FROM medals
		 GROUP BY user_id
		)  SELECT admin,
		          COALESCE(bronze, 0),
		          COALESCE(CASE WHEN show_country THEN country END, ''),
		          COALESCE(diamond, 0),
		          COALESCE(gold, 0),
		          (SELECT COUNT(DISTINCT hole)
		             FROM solutions
		            WHERE user_id = id AND NOT FAILING),
		          id,
		          (SELECT COUNT(DISTINCT lang)
		             FROM solutions
		            WHERE user_id = id AND NOT FAILING),
		          login,
		          COALESCE(bytes.points, 0),
		          COALESCE(chars.points, 0),
		          COALESCE(silver, 0),
		          sponsor,
		          (SELECT COUNT(*) FROM trophies WHERE user_id = id)
		     FROM users
		LEFT JOIN medals       ON id = medals.user_id
		LEFT JOIN points bytes ON id = bytes.user_id AND bytes.scoring = 'bytes'
		LEFT JOIN points chars ON id = chars.user_id AND chars.scoring = 'chars'
		    WHERE login = $1`,
		name,
	).Scan(
		&info.Admin,
		&info.Bronze,
		&info.Country,
		&info.Diamond,
		&info.Gold,
		&info.Holes,
		&info.ID,
		&info.Langs,
		&info.Name,
		&info.BytesPoints,
		&info.CharsPoints,
		&info.Silver,
		&info.Sponsor,
		&info.Cheevos,
	); errors.Is(err, sql.ErrNoRows) {
		return nil
	} else if err != nil {
		panic(err)
	}

	// TODO
	info.TeedOff = time.Date(2019, time.July, 15, 20, 13, 21, 0, time.UTC)

	rows, err := db.Query(
		` SELECT connection, discriminator, id, public, username
		    FROM connections
		   WHERE user_id = $1
		ORDER BY connection`,
		info.ID,
	)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var c Connection

		if err := rows.Scan(
			&c.Connection,
			&c.Discriminator,
			&c.ID,
			&c.Public,
			&c.Username,
		); err != nil {
			panic(err)
		}

		info.Connections = append(info.Connections, c)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	return &info
}
