package hole

import (
	"fmt"
	"math/rand/v2"
	"strings"
)

var hardCodedFractions = [...]struct{ n, d int }{
	{0, 1}, {0, 7}, {233, 144}, {249, 2}, {4, 1}, {87, 1}, {50, 10}, {3, 14},
	{5, 75}, {7, 147}, {242, 11}, {5, 3}, {205, 102}, {234, 233}, {1, 211},
}

var _ = answerFunc("continued-fractions", func() []Answer {
	tests := make([]test, 100)

	for i, fraction := range hardCodedFractions {
		tests[i] = continuedFractionsTest(fraction.n, fraction.d)
	}

	for i := len(hardCodedFractions); i < len(tests); i++ {
		tests[i] = continuedFractionsTest(rand.IntN(251), 1+rand.IntN(250))
	}
	return outputTests(shuffle(tests))
})

func continuedFractionsTest(n, d int) test {
	in := fmt.Sprint(n, "/", d)
	var out strings.Builder

	prefixes := [...]string{"[", "; ", ", "}
	for i := 0; d > 0; i++ {
		fmt.Fprint(&out, prefixes[min(i, len(prefixes)-1)], n/d)

		n %= d
		n, d = d, n
	}

	out.WriteByte(']')

	return test{in, out.String()}
}
