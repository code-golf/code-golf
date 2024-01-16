package hole

import (
	"fmt"
	"math"
)

func perimeter(ai, bi int) (p float64) {
	a, b := float64(ai), float64(bi)
	h := math.Pow(a-b, 2) / math.Pow(a+b, 2)
	for ni := 0; ni < 100; ni++ {
		n := float64(ni)
		bin := math.Gamma(1.5) / (math.Gamma(1.0+n) * math.Gamma(1.5-n))
		p += math.Pow(bin, 2) * math.Pow(h, n)
	}
	p *= math.Pi * (a + b)
	return
}

func ellipsePerimeters() []Run {
	const aMin = 5
	const aMax = 19
	const bMin = 1
	const bMax = 5
	const cases = (aMax - aMin + 1) * (bMax - bMin + 1)
	tests := make([]test, cases)

	// test every a,b combination
	for i := range tests {
		a := aMin + i / (bMax - bMin + 1)
		b := bMin + i % (bMax - bMin + 1)
		tests[i] = test{fmt.Sprint(a, b), fmt.Sprint(int(perimeter(a, b)))}
	}

	return outputTests(shuffle(tests))
}
