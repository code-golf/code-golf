package hole

import (
	"math"
	"math/rand"
	"strconv"
)

func perimeter(ai, bi int) (p float64) {
	var n float64
	a, b := float64(ai), float64(bi)
	h := math.Pow(a-b, 2) / math.Pow(a+b, 2)
	var bin float64
	for ni := 0; ni < 100; ni++ {
		n = float64(ni)
		bin = math.Gamma(1.5) / (math.Gamma(1.0+n) * math.Gamma(1.5-n))
		p += math.Pow(bin, 2) * math.Pow(h, n)
	}
	p *= math.Pi * (a + b)
	return
}

func ellipsePerimeters() []Run {
	const randomCases = 10
	tests := make([]test, randomCases)

	// some random tests
	var a, b int
	for i := 0; i < randomCases; i++ {
		a = rand.Intn(15) + 5
		b = rand.Intn(5) + 1
		tests[i] = test{
			strconv.Itoa(a) + " " + strconv.Itoa(b),
			strconv.Itoa(int(perimeter(a, b))),
		}
	}

	return outputTests(shuffle(tests))
}
