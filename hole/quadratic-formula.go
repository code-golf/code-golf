package hole

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
)

const DEBUG = false

type solution struct {
	n1, d1, n2, d2 int
	sq, im         bool
}

func reduce(n int, d int) (int, int) {
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

func sq_reduce(n int, d int) (int, int) {
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

func solve(a int, b int, c int) solution {
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
			n2, d2 = sq_reduce(disc, 2*a)
		} else {
			disc = int(root)
			n2, d2 = reduce(disc, 2*a)
		}
		return solution{
			n1: n1,
			d1: d1,
			n2: n2,
			d2: d2,
			sq: sq,
			im: im,
		}
	} else {
		n, d := reduce(-c, b)
		return solution{
			n1: n,
			d1: d,
			n2: 0,
			d2: 0,
			sq: false,
			im: false,
		}
	}
}

func (s solution) String() string {
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
			st = fmt.Sprint(st, "±")
		} else {
			st = fmt.Sprint(st, " ± ")
		}

	}

	if s.im {
		st = fmt.Sprint(st, "i")
	}
	if s.sq {
		st = fmt.Sprint(st, "√")
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
	if DEBUG {
		st = fmt.Sprint(st, "  ", s.n1, s.d1, s.n2, s.d2, s.sq, s.im)
	}
	return st
}

const NUM_TESTS = 200

func quadraticFormula() []Scorecard {
	args := make([]string, NUM_TESTS)
	solstrings := make([]string, NUM_TESTS)
	for i := 0; i < NUM_TESTS; i++ {
		a := rand.Intn(20) - 10
		b := rand.Intn(20) - 10
		c := rand.Intn(50) - 25
		s := solve(a, b, c)
		args[i] = fmt.Sprintf("%d %d %d", a, b, c)
		solstrings[i] = fmt.Sprint(s)
	}

	return []Scorecard{{Args: args, Answer: strings.Join(solstrings, "\n")}}
}
