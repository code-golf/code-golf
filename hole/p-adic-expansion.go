package hole

import "fmt"

func reverseBytes(s []byte) []byte {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

const pAdicDigits = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXY"

func pAdicTest(p, n, d int) test {
	in := fmt.Sprint(p, n, "/", d)
	buf := []byte{}
	occurs := make(map[int]int)
	pos := 0
	for {
		if prev, ok := occurs[n]; ok {
			out := fmt.Sprintf("%s'%s", reverseBytes(buf[prev:]), reverseBytes(buf[:prev]))
			return test{in, out}
		}
		occurs[n] = pos
		digit := 0
		for n%p != 0 {
			n -= d
			digit++
		}
		n /= p
		buf = append(buf, pAdicDigits[digit])
		pos++
	}
}

var pAdicTests = [...]struct{ p, n, d int }{
	{2, 26, 89}, {2, 77, 13}, {3, -98, 89}, {3, 0, 1}, {3, 3, 91}, {5, 17, 1},
	{5, 75, 97}, {7, -1, 1}, {61, 1, 71},
}
var pAdicPrimes = []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		a = -a
	}
	return a
}

func pAdicExpansion() []Scorecard {
	tests := make([]test, len(pAdicTests), 100)

	for i, test := range pAdicTests {
		tests[i] = pAdicTest(test.p, test.n, test.d)
	}

	for len(tests) < 100 {
		p := randChoice(pAdicPrimes)
		n := randInt(-99, 99)
		d := randInt(1, 99)
		if d%p != 0 && gcd(n, d) == 1 {
			tests = append(tests, pAdicTest(p, n, d))
		}
	}

	return outputTests(shuffle(tests))
}
