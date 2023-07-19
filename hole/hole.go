package hole

import (
	"math/rand"
	"strings"
)

type test struct{ in, out string }

func outputTests(testRuns ...[]test) []Run {
	runs := make([]Run, len(testRuns))

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

		runs[i] = Run{Args: args, Answer: answer.String()}
	}

	return runs
}

func outputMultirunTests(tests []test) []Run {
	shuffle(tests)
	mid := len(tests) / 2
	return outputTests(tests, tests[:mid], tests[mid:])
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
