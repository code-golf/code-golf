package hole

import (
	"fmt"
	"math"
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
	const (
		aMin = 5
		aLen = 15
		bMin = 1
		bLen = 5
	)

	tests := make([]test, 0, aLen*bLen)

	for i := range aLen {
		for j := range bLen {
			a := aMin + i
			b := bMin + j
			tests = append(tests, test{fmt.Sprint(a, b), fmt.Sprint(int(perimeter(a, b)))})
		}
	}

	shuffle(tests)
	mid := len(tests) / 2
	return outputTests(tests[:mid], tests[mid:])
})
