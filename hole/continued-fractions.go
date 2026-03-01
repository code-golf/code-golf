package hole

import (
	"fmt"
	"math/rand/v2"
)

var hardCodedFractions = [...]struct{ n, d int }{
	{0, 1}, {0, 7}, {233, 144}, {255, 2}, {4, 1}, {87, 1}, {50, 10}, {3, 14}, {5, 75},
	{7, 147}, {253, 11}, {5, 3}, {205, 102}, {234, 233}, {1, 211},
}

var _ = answerFunc("continued-fractions", func() []Answer {
	tests := make([]test, 100)

	for i, fraction := range hardCodedFractions {
		tests[i] = continuedFractionsTest(fraction.n, fraction.d)
	}

	for i := len(hardCodedFractions); i < len(tests); i++ {
		tests[i] = continuedFractionsTest(rand.IntN(256), 1+rand.IntN(255))
	}
	return outputTests(shuffle(tests))
})

func continuedFractionsTest(n, d int) test {
	in := fmt.Sprint(n, "/", d)
	out := fmt.Sprint("[", n/d)
	n %= d
	d, n = n, d
	if d > 0 {
		out = fmt.Sprint(out, "; ")
		for d > 0 {
			out = fmt.Sprint(out, n/d, ", ")
			n %= d
			n, d = d, n
		}
		out = out[:len(out)-2]
	}
	out = fmt.Sprint(out, "]")
	return test{in, out}
}
