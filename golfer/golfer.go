package golfer

import (
	"database/sql"
	"errors"
	"time"
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

	// Start date
	TeedOff time.Time
}

func GetInfo(db *sql.DB, name string) *GolferInfo {
	var info GolferInfo

	if err := db.QueryRow(
		`SELECT admin,
		        (SELECT COUNT(DISTINCT hole)
		           FROM solutions
		          WHERE user_id = id AND NOT FAILING),
		        id,
		        (SELECT COUNT(DISTINCT lang)
		           FROM solutions
		          WHERE user_id = id AND NOT FAILING),
		        login,
		        (SELECT points FROM points WHERE user_id = id),
		        sponsor,
		        (SELECT COUNT(*) FROM trophies WHERE user_id = id)
		   FROM users
		  WHERE login = $1`,
		name,
	).Scan(
		&info.Admin,
		&info.Holes,
		&info.ID,
		&info.Langs,
		&info.Name,
		&info.Points,
		&info.Sponsor,
		&info.Trophies,
	); errors.Is(err, sql.ErrNoRows) {
		return nil
	} else if err != nil {
		panic(err)
	}

	// TODO
	info.Bronze = 28
	info.Gold = 21
	info.Silver = 22
	info.TeedOff = time.Date(2019, time.July, 15, 20, 13, 21, 0, time.UTC)

	return &info
}
