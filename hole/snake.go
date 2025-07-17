package hole

import (
	"math/rand/v2"
	"strings"
)

var directions = [][]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}

func randTrail(length int) string {
	var trail []rune

	visited := map[[2]int]bool{{0, 0}: true}

	for d, x, y := 0, 0, 0; len(trail) < length; {
		move := []rune{'F', 'L', 'R'}[rand.IntN(3)]

		trail = append(trail, move)

		switch move {
		case 'F':
			dx, dy := x+directions[d][0], y+directions[d][1]

			if dx < 0 || dy < 0 || visited[[2]int{dx, dy}] {
				trail = trail[:len(trail)-1]
				continue
			}

			x, y, visited[[2]int{x, y}] = dx, dy, true
		case 'L':
			d = (d + 3) % 4
		case 'R':
			d = (d + 1) % 4
		}
	}

	return string(trail)
}

func printTrail(s string) string {
	var trail strings.Builder

	visited, d, x, y := map[[2]int]bool{{0, 0}: true}, 0, 0, 0

	for _, move := range s {
		switch move {
		case 'F':
			x, y = x+directions[d][0], y+directions[d][1]

			visited[[2]int{x, y}] = true
		case 'L':
			d = (d + 3) % 4
		case 'R':
			d = (d + 1) % 4
		}
	}

	xMin, xMax, yMin, yMax := 0, 0, 0, 0

	for i := range visited {
		xMin = min(xMin, i[0])
		xMax = max(xMax, i[0])
		yMin = min(yMin, i[1])
		yMax = max(yMax, i[1])
	}

	for i := xMin; i <= xMax; i++ {
		for j := yMin; j <= yMax; j++ {
			if visited[[2]int{i, j}] {
				trail.WriteByte('#')
			} else {
				trail.WriteByte(' ')
			}
		}

		trail.WriteByte('\n')
	}

	return trail.String()
}

var _ = answerFunc("snake", func() []Answer {
	tests := make([]test, 100)

	for i := range tests {
		argument := randTrail(randInt(10, randInt(20, 100)))

		tests[i] = test{argument, printTrail(argument)}
	}

	return outputTests(shuffle(tests))
})
