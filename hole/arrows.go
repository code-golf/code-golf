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

var arrowMapDownAndLeft = map[string][2]int8{
	"←": {-1, 0}, "↓": {0, -1}, "↔": {0, 0},
	"↕": {0, 0}, "↙": {-1, -1},
	"↲": {-1, -1},
	"⇐": {-1, 0}, "⇓": {0, -1}, "⇔": {0, 0},
	"⇕": {0, 0}, "⇙": {-1, -1},
	"⇦": {-1, 0}, "⇩": {0, -1},
	"⥀": {0, 0}, "⥁": {0, 0},
}

func arrows() []Scorecard {
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

		outs[i] = fmt.Sprint(pos[0], pos[1])
	}

	// Additional test to force Cartesian quadrant 3
	argsDL := make([]string, 0, 3*len(arrowMapDownAndLeft))
	posDL := [2]int8{}

	// 7-9 of each arrow.
	for arrow := range arrowMapDownAndLeft {
		for times := 6 + rand.Intn(3); times >= 0; times-- {
			argsDL = append(argsDL, arrow)
		}
	}

	// Calculate the outs from the args.
	outsDL := make([]string, len(argsDL))
	for i, arrow := range shuffle(argsDL) {
		coord := arrowMap[arrow]
		posDL[0] += coord[0]
		posDL[1] += coord[1]

		outsDL[i] = fmt.Sprint(posDL[0], posDL[1])
	}

	return []Scorecard{
		{Args: args, Answer: strings.Join(outs, "\n")},
		{Args: argsDL, Answer: strings.Join(outsDL, "\n")},
	}
}
