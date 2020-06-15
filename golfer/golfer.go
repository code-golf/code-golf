package golfer

import (
	"database/sql"
	"time"
)

type Golfer struct {
	Admin bool
	Name  string
}

type GolferInfo struct {
	Golfer

	Sponsor bool

	// Overall points
	Points int

	// Count of medals
	Gold, Silver, Bronze int

	// Count of holes/trophies done
	Holes, Trophies int

	// Start date
	TeeOff time.Time
}

func GetInfo(db *sql.DB, name string) *GolferInfo {
	// TODO
	return &GolferInfo{
		Bronze:   28,
		Gold:     21,
		Golfer:   Golfer{true, "JRaspass"},
		Holes:    41,
		Points:   32418,
		Silver:   22,
		Sponsor:  true,
		TeeOff:   time.Date(2019, time.July, 15, 20, 13, 21, 0, time.UTC),
		Trophies: 11,
	}
}
