package hole

import (
	"math/rand"
	"strings"

	"golang.org/x/exp/constraints"
)

type test struct{ in, out string }

func outputTests(testRuns ...[]test) []Scorecard {
	scores := make([]Scorecard, len(testRuns))

	for i, tests := range testRuns {
		args := make([]string, len(tests))
		var answer strings.Builder

		for i, t := range tests {
			args[i] = t.in

			if i > 0 {
				answer.WriteByte('\n')
			}
			answer.WriteString(t.out)
		}

		scores[i] = Scorecard{Args: args, Answer: answer.String()}
	}

	return scores
}

func outputMultirunTests(tests []test) []Scorecard {
	shuffle(tests)
	mid := len(tests) / 2
	return outputTests(tests, tests[:mid], tests[mid:])
}

// Doesn't handle any special cases, will be in the stdlib/x one day.
func max[T constraints.Ordered](x, y T) T {
	if x > y {
		return x
	}
	return y
}

// Doesn't handle any special cases, will be in the stdlib/x one day.
func min[T constraints.Ordered](x, y T) T {
	if x < y {
		return x
	}
	return y
}

// Returning the slice is a convenience, the shuffle is still in-place.
func shuffle[E any](x []E) []E {
	rand.Shuffle(len(x), func(i, j int) { x[i], x[j] = x[j], x[i] })
	return x
}

func randChoice[E any](x []E) E {
	return x[rand.Intn(len(x))]
}

func randInt(a, b int) int { return rand.Intn(b-a+1) + a }
