package hole

import (
	"fmt"
	"math/rand"
)

func genRepeatingDecimalTest(p, q int) test {
	in := fmt.Sprint(p, "/", q)
	out := fmt.Sprint(p / q)
	p %= q
	if p > 0 {
		buf := []byte{}
		occurs := make([]byte, q)
		var pos byte = 1
		for occurs[p] == 0 {
			occurs[p] = pos
			p *= 10
			buf = append(buf, byte(p/q+'0'))
			p %= q
			pos++
		}
		if p > 0 {
			start := occurs[p] - 1
			out = fmt.Sprintf("%s.%s(%s)", out, buf[:start], buf[start:])
		} else {
			out = fmt.Sprintf("%s.%s", out, buf[:len(buf)-1])
		}
	}

	return test{in, out}
}

var repeatingDecimalFractions = []struct{ p, q int }{
	{0, 1},
	{0, 7},
	{0, 32},
	{1, 1},
	{1, 2},
	{1, 3},
	{1, 4},
	{1, 5},
	{1, 7},
	{1, 10},
	{1, 22},
	{1, 25},
	{1, 28},
	{1, 35},
	{1, 40},
	{1, 50},
	{1, 64},
	{1, 70},
	{1, 96},
	{1, 97},
	{2, 7},
	{2, 14},
	{2, 52},
	{3, 94},
	{4, 2},
	{4, 94},
	{5, 13},
	{5, 65},
	{5, 75},
	{5, 94},
	{10, 7},
	{17, 10},
	{20, 7},
	{40, 2},
	{83, 60},
	{97, 6},
	{97, 7},
	{98, 49},
	{99, 1},
	{99, 3},
}

func repeatingDecimal() []Scorecard {
	tests := make([]test, 0, 100)
	for _, fraction := range repeatingDecimalFractions {
		tests = append(tests, genRepeatingDecimalTest(fraction.p, fraction.q))
	}
	for len(tests) < 100 {
		tests = append(tests, genRepeatingDecimalTest(rand.Intn(100), 1+rand.Intn(99)))
	}
	return outputTests(shuffle(tests))
}
