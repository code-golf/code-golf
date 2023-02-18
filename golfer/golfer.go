package golfer

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/code-golf/code-golf/config"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"golang.org/x/exp/slices"
	"gopkg.in/guregu/null.v4"
)

const (
	followLimit        = 10
	followLimitSponsor = 24
)

type FailingSolutions []struct{ Hole, Lang string }

func (f *FailingSolutions) Scan(src any) error {
	return json.Unmarshal(src.([]byte), f)
}

type Golfer struct {
	Admin, ShowCountry, Sponsor                    bool
	BytesPoints, CharsPoints, ID                   int
	Cheevos, Holes                                 []string
	Country, Layout, Keymap, Name, Referrer, Theme string
	Delete                                         sql.NullTime
	FailingSolutions                               FailingSolutions
	Following                                      []int64
	TimeZone                                       *time.Location
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
		return followLimitSponsor
	}
	return followLimit
}

// IsFollowing returns whether the golfer is following that golfer.
// FIXME Ideally we'd scan into a []int not a []int64 but pq won't.
func (g *Golfer) IsFollowing(userID int) bool {
	_, ok := slices.BinarySearch(g.Following, int64(userID))
	return ok
}

// Solved returns whether the golfer has solved that hole. Counts failing too.
func (g *Golfer) Solved(holeID string) bool {
	_, ok := slices.BinarySearch(g.Holes, holeID)
	return ok
}

type GolferInfo struct {
	Golfer

	// Count of medals
	Diamond, Gold, Silver, Bronze int

	// Count of cheevos/holes/langs done
	Holes, Langs int

	// Count of cheevos/holes/langs available
	CheevosTotal, HolesTotal, LangsTotal int

	// Slice of golfers referred
	Referrals []string

	// Start date
	TeedOff time.Time
}

type RankUpdateFromTo struct {
	Joint   null.Bool `json:"joint"`
	Rank    null.Int  `json:"rank"`
	Strokes null.Int  `json:"strokes"`
}

type RankUpdate struct {
	Scoring        string           `json:"scoring"`
	From           RankUpdateFromTo `json:"from"`
	To             RankUpdateFromTo `json:"to"`
	Beat           null.Int         `json:"beat"`
	OldBestJoint   null.Bool        `json:"oldBestJoint"`
	OldBestStrokes null.Int         `json:"oldBestStrokes"`
}

func GetInfo(db *sqlx.DB, name string) *GolferInfo {
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
		          ARRAY(
		            SELECT trophy
		              FROM trophies
		             WHERE user_id = users.id
		          ORDER BY trophy
		          ),
		          country_flag,
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
		          ARRAY(
		            SELECT login
		              FROM users u
		             WHERE referrer_id = users.id
		          ORDER BY login
		          ),
		          COALESCE(silver, 0),
		          sponsor
		     FROM users
		LEFT JOIN medals       ON id = medals.user_id
		LEFT JOIN points bytes ON id = bytes.user_id AND bytes.scoring = 'bytes'
		LEFT JOIN points chars ON id = chars.user_id AND chars.scoring = 'chars'
		    WHERE login = $1`,
		name,
	).Scan(
		&info.Admin,
		&info.Bronze,
		pq.Array(&info.Cheevos),
		&info.Country,
		&info.Diamond,
		&info.Gold,
		&info.Holes,
		&info.ID,
		&info.Langs,
		&info.Name,
		&info.BytesPoints,
		&info.CharsPoints,
		pq.Array(&info.Referrals),
		&info.Silver,
		&info.Sponsor,
	); errors.Is(err, sql.ErrNoRows) {
		return nil
	} else if err != nil {
		panic(err)
	}

	// TODO
	info.TeedOff = time.Date(2019, time.July, 15, 20, 13, 21, 0, time.UTC)

	return &info
}
