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
		if i[0] < xMin {
			xMin = i[0]
		}

		if i[0] > xMax {
			xMax = i[0]
		}

		if i[1] < yMin {
			yMin = i[1]
		}

		if i[1] > yMax {
			yMax = i[1]
		}
	}

	for i := xMin; i <= xMax; i++ {
		for j := yMin; j <= yMax; j++ {
			if visited[[2]int{i, j}] {
				trail.WriteRune('#')
			} else {
				trail.WriteRune(' ')
			}
		}

		trail.WriteRune('\n')
	}

	return trail.String()
}

func snake() []Run {
	tests := make([]test, 0, 100)

	for range 100 {
		argument := randTrail(randInt(10, randInt(20, 100)))

		tests = append(tests, test{
			argument,
			printTrail(argument),
		})
	}

	return outputTests(shuffle(tests))
}
