package hole

import (
	"fmt"
	"math/rand"
	"strings"
)

var arrowMap = map[string][2]int8{
	// U+2190 - U+2199
	"←": {-1, 0}, "↑": {0, 1}, "→": {1, 0}, "↓": {0, -1}, "↔": {0, 0},
	"↕": {0, 0}, "↖": {-1, 1}, "↗": {1, 1}, "↘": {1, -1}, "↙": {-1, -1},

	// U+21B0 - U+21B3
	"↰": {-1, 1}, "↱": {1, 1}, "↲": {-1, -1}, "↳": {1, -1},

	// U+21D0 - U+21D9
	"⇐": {-1, 0}, "⇑": {0, 1}, "⇒": {1, 0}, "⇓": {0, -1}, "⇔": {0, 0},
	"⇕": {0, 0}, "⇖": {-1, 1}, "⇗": {1, 1}, "⇘": {1, -1}, "⇙": {-1, -1},

	// U+21E6 - U+21E9
	"⇦": {-1, 0}, "⇧": {0, 1}, "⇨": {1, 0}, "⇩": {0, -1},

	// U+2940 - U+2941
	"⥀": {0, 0}, "⥁": {0, 0},
}

func arrows() ([]string, string) {
	args := make([]string, 0, 3*len(arrowMap))
	pos := [2]int8{}

	// 1-3 of each arrow.
	for arrow := range arrowMap {
		for times := rand.Intn(3); times >= 0; times-- {
			args = append(args, arrow)
		}
	}

	// Calculate the outs from the args.
	outs := make([]string, len(args))
	for i, arrow := range shuffle(args) {
		coord := arrowMap[arrow]
		pos[0] += coord[0]
		pos[1] += coord[1]

		outs[i] = fmt.Sprintf("%d %d", pos[0], pos[1])
	}

	return args, strings.Join(outs, "\n")
}
