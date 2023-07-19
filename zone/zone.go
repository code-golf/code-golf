package zone

import (
	"cmp"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"slices"
	"time"
)

type Zone struct {
	Name   string
	Offset int
}

var locations []*time.Location

// ByID TODO Actually point to something.
var ByID = map[string]bool{}

var Country = map[string]string{}

func init() {
	file, err := os.Open("/usr/share/zoneinfo/zone1970.tab")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comment = '#'
	reader.Comma = '\t'
	reader.FieldsPerRecord = -1

	for {
		switch row, err := reader.Read(); err {
		case nil:
			location, err := time.LoadLocation(row[2])
			if err != nil {
				panic(err)
			}
			locations = append(locations, location)

			ByID[row[2]] = true

			// If the zone maps exactly to one country.
			if len(row[0]) == 2 {
				Country[row[2]] = row[0]
			}
		case io.EOF:
			return
		default:
			panic(err)
		}
	}
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func (z Zone) String() string {
	mins := z.Offset / 60
	return fmt.Sprintf("(%+03d:%02d) %s", mins/60, abs(mins%60), z.Name)
}

func List() []Zone {
	now := time.Now()
	zones := make([]Zone, len(locations))

	for i, location := range locations {
		_, offset := now.In(location).Zone()
		zones[i] = Zone{location.String(), offset}
	}

	slices.SortFunc(zones, func(a, b Zone) int {
		if c := cmp.Compare(a.Offset, b.Offset); c != 0 {
			return c
		}
		return cmp.Compare(a.Name, b.Name)
	})

	return zones
}
