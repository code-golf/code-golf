package hole

import (
	"fmt"
	"math"
	"math/big"
	"strings"
)

var _ = answerFunc("smooth-numbers", func() []Answer {
	answers := make([]Answer, 5)

	for j := range answers {
		n, k := randInt(1, math.MaxInt8), randInt(1, math.MaxInt16)

		var primes []int

		for i := range n + 1 {
			if big.NewInt(int64(i)).ProbablyPrime(0) {
				primes = append(primes, i)
			}
		}

		var expected strings.Builder

		for i := 1; i <= k; i++ {
			x := i

			for _, p := range primes {
				for x%p == 0 {
					x /= p
				}
			}

			if x == 1 {
				if i != 1 {
					expected.WriteByte('\n')
				}
				fmt.Fprint(&expected, i)
			}
		}

		answers[j] = Answer{
			Args:   []string{fmt.Sprint(n, k)},
			Answer: expected.String(),
		}
	}

	return answers
})
