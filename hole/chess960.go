package hole

import (
	"math/rand/v2"
	"strconv"
	"strings"
)

var _ = answerFunc("chess960-decoder", func() []Answer { return chess960(true) })
var _ = answerFunc("chess960-encoder", func() []Answer { return chess960(false) })

func chess960(reverse bool) []Answer {
	const count = 960

	knightArrangement := [10][2]int{{0, 1}, {0, 2}, {0, 3}, {0, 4}, {1, 2}, {1, 3}, {1, 4}, {2, 3}, {2, 4}, {3, 4}}

	tests := make([]test, count)
	for x, n := range rand.Perm(count) {
		spid := strconv.Itoa(n)
		grid := [8]rune{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '}

		n2, b1 := n/4, n%4
		grid[b1*2+1] = '♗'

		n3, b2 := n2/4, n2%4
		grid[b2*2] = '♗'

		var remains [6]int
		ind := 0
		for i, c := range grid {
			if c == ' ' {
				remains[ind] = i
				ind++
			}
		}
		n4, q := n3/6, n3%6
		grid[remains[q]] = '♕'

		var remains2 [5]int
		ind2 := 0
		for i, c := range grid {
			if c == ' ' {
				remains2[ind2] = i
				ind2++
			}
		}
		knights := knightArrangement[n4]
		grid[remains2[knights[0]]] = '♘'
		grid[remains2[knights[1]]] = '♘'

		var remains3 [3]int
		ind3 := 0
		for i, c := range grid {
			if c == ' ' {
				remains3[ind3] = i
				ind3++
			}
		}
		grid[remains3[0]] = '♖'
		grid[remains3[1]] = '♔'
		grid[remains3[2]] = '♖'

		var layout strings.Builder
		for _, piece := range grid {
			layout.WriteRune(piece)
		}

		if reverse {
			tests[x] = test{spid, layout.String()}
		} else {
			tests[x] = test{layout.String(), spid}
		}
	}

	return outputTests(tests[:480], tests[480:])
}
