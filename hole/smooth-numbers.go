package hole

import (
	"fmt"
	"math"
	"strings"
)

var _ = answerFunc("smooth-numbers", func() []Answer {
	var tests []test

	for range 100 {
		n, k := randInt(0, math.MaxInt8), randInt(0, math.MaxInt16)

		flags := make([]bool, n+1)

		var primes []int

		for i := 2; i <= n; i++ {
			if !flags[i] {
				primes = append(primes, i)

				for j := i; i*j <= n; j++ {
					flags[i*j] = true
				}
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
				expected.WriteString(fmt.Sprintf("%d ", i))
			}
		}

		tests = append(tests, test{fmt.Sprintf("%d %d", n, k), expected.String()})
	}

	return outputTests(shuffle(tests))
})
