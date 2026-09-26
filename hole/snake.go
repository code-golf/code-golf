package hole

import (
	"math/rand/v2"
	"strings"
)

var directions = [][]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}

func randTrail(length int) string {
	var trail []rune

	visited := map[[2]int]bool{{0, 0}: true}

	for d, y, x := 0, 0, 0; len(trail) < length; {
		move := []rune{'F', 'L', 'R'}[rand.IntN(3)]

		trail = append(trail, move)

		switch move {
		case 'F':
			dy, dx := y+directions[d][0], x+directions[d][1]

			if dx < 0 || dy < 0 || dx > 19 || dy > 10 || visited[[2]int{dy, dx}] {
				trail = trail[:len(trail)-1]
				continue
			}

			y, x, visited[[2]int{y, x}] = dy, dx, true
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

	visited, d, y, x := map[[2]int]bool{{0, 0}: true}, 0, 0, 0

	for _, move := range s {
		switch move {
		case 'F':
			y, x = y+directions[d][0], x+directions[d][1]

			visited[[2]int{y, x}] = true
		case 'L':
			d = (d + 3) % 4
		case 'R':
			d = (d + 1) % 4
		}
	}

	yMin, yMax, xMin, xMax := 0, 0, 0, 0

	for i := range visited {
		yMin = min(yMin, i[0])
		yMax = max(yMax, i[0])
		xMin = min(xMin, i[1])
		xMax = max(xMax, i[1])
	}

	for i := yMin; i <= yMax; i++ {
		for j := xMin; j <= xMax; j++ {
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
