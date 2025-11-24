package hole

import (
	"fmt"
	"math"
	"math/big"
	"strings"
)

hardcodedTests = [...]struct{ n, k int }{
	{2, 32768}, {3, 529},
}

var _ = answerFunc("smooth-numbers", func() []Answer {
	answers := make([]Answer, 5)

	for j, test := range hardcodedTests {
		n, k := test.n, test.k

		answers[j] = Answer{
			Args:   []string{fmt.Sprint(n), fmt.Sprint(k)},
			Answer: smoothOutput(n, k),
		}
	}

	for j := len(hardcodedTests); j < len(answers); j++ {
		n, k := randInt(1, math.MaxInt8), randInt(1, math.MaxInt16)

		answers[j] = Answer{
			Args:   []string{fmt.Sprint(n), fmt.Sprint(k)},
			Answer: smoothOutput(n, k),
		}
	}

	return shuffle(answers)
})

func smoothOutput(n, k int) string {
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
	
	return expected.String()
}
