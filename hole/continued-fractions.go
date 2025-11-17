package hole

import (
	"fmt"
	"math/rand/v2"
)

var _ = answerFunc("repeating-decimals", func() []Answer {
	tests := make([]test, 100)

	for i := range 100 {
		tests[i] = continuedFractionsTest(rand.IntN(256), 1+rand.IntN(255))
	}

	return outputTests(shuffle(tests))
})

func continuedFractionsTest(n, d int) string {
	in := fmt.Sprint(n, "/", d)
	out := fmt.Sprint("[", n/d)
	n %= d
	if n > 0 {
		out = fmt.Sprint(out, ";")
		for n%d > 0 {
			out = fmt.Sprint(out, n/d, ",")
			n %= d
			n, d = d, n
		}
		out = out[:len(out)-1]
	}
	out = fmt.Sprint(out, "]")
	return test{in, out}
}
