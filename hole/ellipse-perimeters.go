package hole

import (
	"math"
	"math/rand"
	"strconv"
	"strings"
)

func perimeter(ai, bi int) (p float64) {
	a, b := float64(ai), float64(bi)
	h := math.Pow(a-b, 2) / math.Pow(a+b, 2)
	p = math.Pi * (a + b) * (1.0 + (3.0 * h / (10.0 + math.Sqrt(4.0-(3.0*h)))))
	return
}

func ellipsePerimeters() (args []string, out string) {
	var outs []string

	// some random tests
	var a, b int
	var p float64
	var ps string
	for i := 0; i < 50; i++ {
		a = rand.Intn(50) + 1
		b = rand.Intn(50) + 1
		args = append(args, strconv.Itoa(a)+" "+strconv.Itoa(b))

		p = perimeter(a, b)
		ps = strconv.FormatFloat(p, 'f', 40, 64)[0:11]
		outs = append(outs, ps)
	}

	rand.Shuffle(len(args), func(i, j int) {
		args[i], args[j] = args[j], args[i]
		outs[i], outs[j] = outs[j], outs[i]
	})

	out = strings.Join(outs, "\n")
	return
}
