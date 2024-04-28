package golfer

import (
	"database/sql/driver"
	"encoding/json"
	"slices"
	"time"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/null"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

const (
	FollowLimit        = 10
	FollowLimitSponsor = 24
)

// Golfer is the info of a logged in golfer we need on every request.
type Golfer struct {
	Admin, ShowCountry, Sponsor           bool
	BytesPoints, CharsPoints, ID          int
	Cheevos, Holes                        pq.StringArray
	Country                               config.NullCountry
	Delete                                null.Time
	FailingSolutions                      FailingSolutions
	Following                             pq.Int64Array
	Keymap, Layout, Name, Referrer, Theme string
	Pronouns, TimeZone                    null.String
	Settings                              Settings
}

// GolferInfo is populated when looking at a /golfers/xxx route.
type GolferInfo struct {
	Golfer

	// Count of medals
	Diamond, Gold, Silver, Bronze int

	// Count of cheevos/holes/langs done
	Holes, Langs int

	// Count of cheevos/holes/langs available
	CheevosTotal, HolesTotal, LangsTotal int

	// Slice of golfers referred
	Referrals pq.StringArray

	// Start date
	Started time.Time
}

type FailingSolutions []struct{ Hole, Lang string }

// FIXME I'm not sure these RankUpdate structs belong here.
type RankUpdateFromTo struct {
	Joint   null.Bool `json:"joint"`
	Rank    null.Int  `json:"rank"`
	Strokes null.Int  `json:"strokes"`
}

type RankUpdate struct {
	Scoring            string           `json:"scoring"`
	From               RankUpdateFromTo `json:"from"`
	To                 RankUpdateFromTo `json:"to"`
	OldBestGolferCount null.Int         `json:"oldBestGolferCount"` // Number of golfers that previously held the gold medal (except current golfer).
	OldBestGolferID    null.Int         `json:"oldBestGolferID"`    // ID of the golfer that previously held the diamond (except current golfer).
	OldBestStrokes     null.Int         `json:"oldBestStrokes"`     // Number of strokes for previous diamond (including current golfer).
	OldBestSubmitted   null.Time        `json:"oldBestSubmitted"`   // Timestamp for previous diamond (including current golfer).
}

// Settings is page → setting → value.
type Settings map[string]map[string]string

func (f *FailingSolutions) Scan(src any) error {
	return json.Unmarshal(src.([]byte), f)
}

// Earn the given cheevo, no-op if already earned.
func (g *Golfer) Earn(db *sqlx.DB, cheevoID string) (earned *config.Cheevo) {
	if rowsAffected, _ := db.MustExec(
		"INSERT INTO trophies VALUES (DEFAULT, $1, $2) ON CONFLICT DO NOTHING",
		g.ID,
		cheevoID,
	).RowsAffected(); rowsAffected == 1 {
		earned = config.CheevoByID[cheevoID]
	}

	// Update g.Cheevos if necessary.
	if i, ok := slices.BinarySearch(g.Cheevos, cheevoID); !ok {
		g.Cheevos = slices.Insert(g.Cheevos, i, cheevoID)
	}

	return
}

// Earned returns whether the golfer has that cheevo.
func (g *Golfer) Earned(cheevoID string) bool {
	_, ok := slices.BinarySearch(g.Cheevos, cheevoID)
	return ok
}

// FollowLimit returns the max number of golfers this golfer can follow.
func (g *Golfer) FollowLimit() int {
	if g.Sponsor {
		return FollowLimitSponsor
	}
	return FollowLimit
}

// IsFollowing returns whether the golfer is following that golfer.
// FIXME Ideally we'd scan into a []int not a []int64 but pq won't.
func (g *Golfer) IsFollowing(userID int) bool {
	_, ok := slices.BinarySearch(g.Following, int64(userID))
	return ok
}

func (g *Golfer) Location() (loc *time.Location) {
	if loc, _ = time.LoadLocation(g.TimeZone.V); loc == nil {
		loc = time.UTC
	}
	return
}

// Solved returns whether the golfer has solved that hole. Counts failing too.
func (g *Golfer) Solved(holeID string) bool {
	_, ok := slices.BinarySearch(g.Holes, holeID)
	return ok
}

func (s *Settings) Scan(v any) error { return json.Unmarshal(v.([]byte), &s) }

func (s Settings) Value() (driver.Value, error) { return json.Marshal(s) }
