package golfer

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/code-golf/code-golf/config"
	"github.com/lib/pq"
	"golang.org/x/exp/slices"
	"gopkg.in/guregu/null.v4"
)

type FailingSolutions []struct{ Hole, Lang string }

func (f *FailingSolutions) Scan(src any) error {
	return json.Unmarshal(src.([]byte), f)
}

type Golfer struct {
	Admin, ShowCountry                     bool
	Cheevos                                []string
	Country, Keymap, Name, Referrer, Theme string
	Delete                                 sql.NullTime
	FailingSolutions                       FailingSolutions
	Following                              []int64
	ID                                     int
	TimeZone                               *time.Location
}

// Earn the given cheevo, no-op if already earned.
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

// IsFollowing returns whether the golfer is following that golfer.
// FIXME Ideally we'd scan into a []int not a []int64 but pq won't.
func (g *Golfer) IsFollowing(userID int) bool {
	_, ok := slices.BinarySearch(g.Following, int64(userID))
	return ok
}

type GolferInfo struct {
	Golfer

	Sponsor bool

	// Overall points
	BytesPoints, CharsPoints int

	// Count of medals
	Diamond, Gold, Silver, Bronze int

	// Count of cheevos/holes/langs done
	Holes, Langs int

	// Count of cheevos/holes/langs available
	CheevosTotal, HolesTotal, LangsTotal int

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
