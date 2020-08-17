package golfer

import (
	"database/sql"
	"errors"
	"time"

	"github.com/code-golf/code-golf/hole"
	"github.com/code-golf/code-golf/lang"
	"github.com/code-golf/code-golf/trophy"
)

type Golfer struct {
	Admin bool
	ID    int
	Name  string
}

type GolferInfo struct {
	Golfer

	Sponsor bool

	// Overall points
	Points int

	// Count of medals
	Gold, Silver, Bronze int

	// Count of holes/langs/trophies done
	Holes, Langs, Trophies int

	// Count of holes/langs/trophies available
	HolesTotal, LangsTotal, TrophiesTotal int

	// Start date
	TeedOff time.Time
}

func GetInfo(db *sql.DB, name string) *GolferInfo {
	info := GolferInfo{
		HolesTotal:    len(hole.List),
		LangsTotal:    len(lang.List),
		TrophiesTotal: len(trophy.List),
	}

	if err := db.QueryRow(
		`WITH ranked AS (
		    SELECT user_id,
		           RANK() OVER (PARTITION BY hole, lang ORDER BY chars)
		      FROM solutions
		      JOIN code ON code_id = id
		     WHERE NOT failing
		), medals AS (
		    SELECT user_id,
		           COUNT(*) FILTER (WHERE rank = 1) gold,
		           COUNT(*) FILTER (WHERE rank = 2) silver,
		           COUNT(*) FILTER (WHERE rank = 3) bronze
		      FROM ranked
		  GROUP BY user_id
		)  SELECT admin,
		          COALESCE(bronze, 0),
		          COALESCE(gold, 0),
		          (SELECT COUNT(DISTINCT hole)
		             FROM solutions
		            WHERE user_id = id AND NOT FAILING),
		          id,
		          (SELECT COUNT(DISTINCT lang)
		             FROM solutions
		            WHERE user_id = id AND NOT FAILING),
		          login,
		          COALESCE(points, 0),
		          COALESCE(silver, 0),
		          sponsor,
		          (SELECT COUNT(*) FROM trophies WHERE user_id = id)
		     FROM users
		LEFT JOIN medals ON id = medals.user_id
		LEFT JOIN points ON id = points.user_id
		    WHERE login = $1`,
		name,
	).Scan(
		&info.Admin,
		&info.Bronze,
		&info.Gold,
		&info.Holes,
		&info.ID,
		&info.Langs,
		&info.Name,
		&info.Points,
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
