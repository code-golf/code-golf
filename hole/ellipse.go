package hole

import (
	"math/big"
	"math/rand"
	"strconv"
	"strings"
)

func perimeter(a, b int) (p float64) {
	//TODO: compute perimeter
}

func ellipse() (args []string, out string) {
	var outs []string

	// few random tests
	var a, b int
	for i := 0; i < 10; i++ {
		a = rand.Intn(10) + 1
		b = rand.Intn(10) + 1

		args = append(args, strconv.Itoa(a)+" "+strconv.Itoa(b))
		outs = append(outs, strconv.Itoa(perimeter(a, b)))
	}

	//TODO: add default tests

	rand.Shuffle(len(args), func(i, j int) {
		args[i], args[j] = args[j], args[i]
		outs[i], outs[j] = outs[j], outs[i]
	})

	out = strings.Join(outs, "\n")
	return
}
