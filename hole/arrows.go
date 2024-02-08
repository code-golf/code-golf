package hole

import (
	"fmt"
	"math/rand/v2"
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

var arrowMapDownAndLeft = [...]string{"←", "↓", "↙", "↲", "⇐", "⇓", "⇙", "⇦", "⇩"}
var arrowMapUpAndLeft = [...]string{"←", "↑", "↖", "↰", "⇐", "⇑", "⇖", "⇦", "⇧"}
var arrowMapUpAndRight = [...]string{"↑", "→", "↗", "↱", "⇑", "⇒", "⇗", "⇧", "⇨"}
var arrowMapDownAndRight = [...]string{"→", "↓", "↘", "↳", "⇒", "⇓", "⇘", "⇨", "⇩"}

func arrows() []Run {
	args := make([]string, 0, 3*len(arrowMap))

	// 1-3 of each arrow.
	for arrow := range arrowMap {
		for times := rand.IntN(3); times >= 0; times-- {
			args = append(args, arrow)
		}
	}

	// Additional test to force all Cartesian quadrants
	argsDL := make([]string, 0, 4*len(arrowMapDownAndLeft))
	argsUL := make([]string, 0, 4*len(arrowMapUpAndLeft))
	argsDR := make([]string, 0, 4*len(arrowMapDownAndRight))
	argsUR := make([]string, 0, 4*len(arrowMapUpAndRight))

	timesDL := 2 + rand.IntN(2)
	timesUL := 5 - timesDL
	timesDR := 2 + rand.IntN(2)
	timesUR := 5 - timesDR

	// 3 or 4 of each arrow.
	for _, arrow := range arrowMapDownAndLeft {
		for times := timesDL; times >= 0; times-- {
			argsDL = append(argsDL, arrow)
		}
	}
	for _, arrow := range arrowMapUpAndLeft {
		for times := timesUL; times >= 0; times-- {
			argsUL = append(argsUL, arrow)
		}
	}
	for _, arrow := range arrowMapDownAndRight {
		for times := timesDR; times >= 0; times-- {
			argsDR = append(argsDR, arrow)
		}
	}
	for _, arrow := range arrowMapUpAndRight {
		for times := timesUR; times >= 0; times-- {
			argsUR = append(argsUR, arrow)
		}
	}

	var testRuns [][]test
	for _, arrows := range [][]string{args, argsDL, argsUL, argsDR, argsUR} {
		tests := make([]test, len(arrows))
		testRuns = append(testRuns, tests)

		var pos [2]int8
		for i, arrow := range shuffle(arrows) {
			coord := arrowMap[arrow]
			pos[0] += coord[0]
			pos[1] += coord[1]

			tests[i] = test{arrow, fmt.Sprint(pos[0], pos[1])}
		}
	}

	return outputTests(testRuns...)
}
