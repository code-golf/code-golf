package golfer

import (
	"database/sql/driver"
	"encoding/json"
	"slices"
	"time"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/db"
	"github.com/code-golf/code-golf/null"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

const (
	FollowLimit        = 10
	FollowLimitSponsor = 24
)

// Golfer is the info of a logged in golfer we need on every request.
// TODO Some of this stuff isn't needed on every request but is needed to
// populate golfer settings, that should be fixed.
type Golfer struct {
	About, Name, Referrer, Theme          string
	Admin, HasNotes, ShowCountry, Sponsor bool
	BytesPoints, CharsPoints, ID          int
	Cheevos, Holes                        pq.StringArray
	Country                               *config.Country
	Delete                                null.Time
	FailingSolutions                      FailingSolutions
	Following                             pq.Int64Array
	Pronouns, TimeZone                    null.String
	Settings                              Settings
}

// GolferInfo is populated when looking at a /golfers/xxx route.
type GolferInfo struct {
	Golfer

	// Count of medals
	Unicorn, Diamond, Gold, Silver, Bronze int

	// Count of cheevos/holes/langs done
	Holes, Langs int

	HolesAuthored config.Holes

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
	Scoring                   string           `json:"scoring"`
	FailingStrokes            null.Int         `json:"failingStrokes"` // The length of the previous failing solution, if any,
	From                      RankUpdateFromTo `json:"from"`
	To                        RankUpdateFromTo `json:"to"`
	OldBestCurrentGolferCount null.Int         `json:"oldBestCurrentGolferCount"` // Number of golfers that previously held the gold medal (except current golfer).
	OldBestCurrentGolferID    null.Int         `json:"oldBestCurrentGolferID"`    // ID of the golfer that previously held the diamond (except current golfer), if there is exactly one.
	OldBestFirstGolferID      null.Int         `json:"oldBestFirstGolferID"`      // ID of the first golfer that obtained the previous diamond (including current golfer).
	OldBestStrokes            null.Int         `json:"oldBestStrokes"`            // Number of strokes for previous gold medal (including current golfer).
	OldBestSubmitted          null.Time        `json:"oldBestSubmitted"`          // Timestamp for previous diamond (including current golfer).
	NewSolutionCount          null.Int         `json:"newSolutionCount"`          // Number of golfers with solutions for this hole/lang/scoring (including current golfer).
}

// Settings is page → setting → value.
type Settings map[string]map[string]any

func (f *FailingSolutions) Scan(src any) error {
	return json.Unmarshal(src.([]byte), f)
}

// Earn the given cheevo, no-op if already earned.
func (g *Golfer) Earn(db db.Queryable, cheevoID string) (earned *config.Cheevo) {
	if rowsAffected, _ := db.MustExec(
		"INSERT INTO cheevos (user_id, cheevo) VALUES ($1, $2) ON CONFLICT DO NOTHING",
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

// SaveSettings saves golfer.Settings back to the DB.
func (g *Golfer) SaveSettings(db *sqlx.DB) {
	// Optimisation, trim the default values from the maps before saving.
	for page, settings := range config.Settings {
		for _, setting := range settings {
			if g.Settings[page][setting.ID] == setting.Default {
				delete(g.Settings[page], setting.ID)
			}
		}

		if len(g.Settings[page]) == 0 {
			delete(g.Settings, page)
		}
	}

	db.MustExec("UPDATE users SET settings = $1 WHERE id = $2", g.Settings, g.ID)
}

// Solved returns whether the golfer has solved that hole. Counts failing too.
func (g *Golfer) Solved(holeID string) bool {
	_, ok := slices.BinarySearch(g.Holes, holeID)
	return ok
}

func (g Golfer) SponsorOrAdmin() bool { return g.Sponsor || g.Admin }

func (g *Golfer) Value() (driver.Value, error) {
	if g == nil {
		return nil, nil
	}
	return int64(g.ID), nil
}

func (s *Settings) Scan(v any) error { return json.Unmarshal(v.([]byte), &s) }

func (s Settings) Value() (driver.Value, error) { return json.Marshal(s) }
