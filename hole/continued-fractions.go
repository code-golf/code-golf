package hole

import (
	"fmt"
	"math/rand/v2"
)

var _ = answerFunc("continued-fractions", func() []Answer {
	tests := make([]test, 100)

	for i := range 100 {
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
