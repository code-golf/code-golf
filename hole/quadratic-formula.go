package hole

import (
	"fmt"
	"math"
	"math/rand/v2"
	"strings"
)

type quadraticSolution struct {
	n1, d1, n2, d2 int
	sq, im         bool
}

func reduce(n, d int) (int, int) {
	if d < 0 {
		d = -d
		n = -n
	}
	for i := d; i > 0; i-- {
		if n%i == 0 && d%i == 0 {
			return n / i, d / i
		}
	}
	return n, d
}

func sqReduce(n, d int) (int, int) {
	if d < 0 {
		d = -d
	}
	for i := d; i > 0; i-- {
		if n%(i*i) == 0 && d%i == 0 {
			return n / (i * i), d / i
		}
	}
	return n, d
}

func solve(a, b, c int) quadraticSolution {
	if a != 0 {
		disc := b*b - 4*a*c
		im := disc < 0
		if im {
			disc = -disc
		}
		root := math.Sqrt(float64(disc))
		sq := math.Floor(root) != root
		n1, d1 := reduce(-b, 2*a)
		var n2, d2 int
		if sq {
			n2, d2 = sqReduce(disc, 2*a)
		} else {
			disc = int(root)
			n2, d2 = reduce(disc, 2*a)
		}
		return quadraticSolution{
			n1: n1,
			d1: d1,
			n2: n2,
			d2: d2,
			sq: sq,
			im: im,
		}
	} else {
		n, d := reduce(-c, b)
		return quadraticSolution{
			n1: n,
			d1: d,
			n2: 0,
			d2: 0,
			sq: false,
			im: false,
		}
	}
}

func (s quadraticSolution) String() string {
	if s.d1 == 0 {
		if s.n1 == 0 {
			return "indeterminate"
		} else {
			return "undefined"
		}
	}
	st := ""
	if s.n1 != 0 {
		st = fmt.Sprint(s.n1)
		if s.d1 != 1 {
			st = fmt.Sprint(st, "/", s.d1)
		}
	}
	if s.n2 != 0 {
		if st == "" {
			st = "±"
		} else {
			st += " ± "
		}
	}
	if s.im {
		st += "i"
	}
	if s.sq {
		st += "√"
	}
	if s.n2 != 0 {
		if !(s.n2 == 1 && s.d2 == 1 && s.im) {
			st = fmt.Sprint(st, s.n2)
		}
		if s.d2 != 1 {
			st = fmt.Sprint(st, "/", s.d2)
		}
	}
	if st == "" {
		st = "0"
	}
	return st
}

func quadraticFormula() []Run {
	const numTests = 200
	args := make([]string, numTests)
	solstrings := make([]string, numTests)
	for i := range numTests {
		a := rand.IntN(20) - 10
		b := rand.IntN(20) - 10
		c := rand.IntN(50) - 25
		args[i] = fmt.Sprint(a, b, c)
		solstrings[i] = solve(a, b, c).String()
	}
	return []Run{{Args: args, Answer: strings.Join(solstrings, "\n")}}
}
