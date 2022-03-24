package hole

import (
	"math/rand"
	"strings"

	"golang.org/x/exp/constraints"
)

type test struct{ in, out string }

func outputTests(tests []test) ([]string, string) {
	args := make([]string, len(tests))
	var answer strings.Builder

	for i, t := range tests {
		args[i] = t.in

		if i > 0 {
			answer.WriteByte('\n')
		}
		answer.WriteString(t.out)
	}

	return args, answer.String()
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
