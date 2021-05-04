package golfer

import (
	"database/sql"
	"encoding/json"
	"errors"
	"sort"
	"time"

	"github.com/code-golf/code-golf/hole"
	"github.com/code-golf/code-golf/lang"
	"github.com/code-golf/code-golf/trophy"
	"gopkg.in/guregu/null.v4"
)

type FailingSolutions []struct{ Hole, Lang string }

func (f *FailingSolutions) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), f)
}

type Golfer struct {
	Admin, ShowCountry              bool
	Cheevos                         []string
	Country, Keymap, Name, Referrer string
	Delete                          sql.NullTime
	FailingSolutions                FailingSolutions
	ID                              int
	TimeZone                        *time.Location
}

// Earn the given cheevo, no-op if already earnt.
func (g *Golfer) Earn(db *sql.DB, cheevoID string) {
	if _, err := db.Exec(
		"INSERT INTO trophies VALUES (DEFAULT, $1, $2) ON CONFLICT DO NOTHING",
		g.ID,
		cheevoID,
	); err != nil {
		panic(err)
	}

	// Update g.Cheevos if necessary.
	if i := sort.SearchStrings(
		g.Cheevos, cheevoID,
	); i == len(g.Cheevos) || g.Cheevos[i] != cheevoID {
		g.Cheevos = append(g.Cheevos, "")
		copy(g.Cheevos[i+1:], g.Cheevos[i:])
		g.Cheevos[i] = cheevoID
	}
}

// Earnt returns whether the golfer has that cheevo.
func (g *Golfer) Earnt(cheevoID string) bool {
	i := sort.SearchStrings(g.Cheevos, cheevoID)
	return i < len(g.Cheevos) && g.Cheevos[i] == cheevoID
}

type GolferInfo struct {
	Golfer

	Sponsor bool

	// Overall points
	BytesPoints, CharsPoints int

	// Count of medals
	Diamond, Gold, Silver, Bronze int

	// Count of holes/langs/trophies done
	Holes, Langs, Trophies int

	// Count of holes/langs/trophies available
	HolesTotal, LangsTotal, TrophiesTotal int

	// Start date
	TeedOff time.Time
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
		HolesTotal:    len(hole.List),
		LangsTotal:    len(lang.List),
		TrophiesTotal: len(trophy.List),
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
		          COALESCE(bytes_points, 0),
		          COALESCE(chars_points, 0),
		          COALESCE(silver, 0),
		          sponsor,
		          (SELECT COUNT(*) FROM trophies WHERE user_id = id)
		     FROM users
		LEFT JOIN medals ON id = medals.user_id
		LEFT JOIN bytes_points ON id = bytes_points.user_id
		LEFT JOIN chars_points ON id = chars_points.user_id
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
		&info.Trophies,
	); errors.Is(err, sql.ErrNoRows) {
		return nil
	} else if err != nil {
		panic(err)
	}

	// TODO
	info.TeedOff = time.Date(2019, time.July, 15, 20, 13, 21, 0, time.UTC)

	return &info
}
