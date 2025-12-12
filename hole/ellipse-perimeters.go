package hole

import (
	"fmt"
	"math"
	"math/rand/v2"
)

func perimeter(ai, bi int) (p float64) {
	a, b := float64(ai), float64(bi)
	h := math.Pow(a-b, 2) / math.Pow(a+b, 2)
	for ni := range 100 {
		n := float64(ni)
		bin := math.Gamma(1.5) / (math.Gamma(1.0+n) * math.Gamma(1.5-n))
		p += math.Pow(bin, 2) * math.Pow(h, n)
	}
	p *= math.Pi * (a + b)
	return
}

var _ = answerFunc("ellipse-perimeters", func() []Answer {
	tests := make([]test, 10)

	// some random tests
	for i := range tests {
		a := rand.IntN(15) + 5
		b := rand.IntN(5) + 1
		tests[i] = test{fmt.Sprint(a, b), fmt.Sprint(int(perimeter(a, b)))}
	}

	return outputTests(tests)
})
