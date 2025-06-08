package hole

import (
	"fmt"
	"math/rand/v2"
	"strconv"
	"strings"
)

func computeZeckendorf(n int32) string {
	// a list of fibonacci values less than 2^31
	fib := []int32{1836311903, 1134903170, 701408733, 433494437, 267914296, 165580141, 102334155, 63245986, 39088169, 24157817, 14930352, 9227465, 5702887, 3524578, 2178309, 1346269, 832040, 514229, 317811, 196418, 121393, 75025, 46368, 28657, 17711, 10946, 6765, 4181, 2584, 1597, 987, 610, 377, 233, 144, 89, 55, 34, 21, 13, 8, 5, 3, 2, 1}

	out := make([]int32, 0)

	for _, f := range fib {
		if f > n {
			continue
		}
		n -= f
		out = append(out, f)
	}

	return strings.Trim(strings.ReplaceAll(fmt.Sprint(out), " ", " + "), "[]")
}

var _ = answerFunc("zeckendorf-representation", func() []Answer {
	fixedCases := []test{
		{"64", "55 + 8 + 1"},
		{"89", "89"},
		{"144", "144"},
		{"701408733", "701408733"},
		{"7", "5 + 2"},
		{"1836311905", "1836311903 + 2"},
		{"1568397607", "1134903170 + 433494437"},
		{"165580140", "102334155 + 39088169 + 14930352 + 5702887 + 2178309 + 832040 + 317811 + 121393 + 46368 + 17711 + 6765 + 2584 + 987 + 377 + 144 + 55 + 21 + 8 + 3 + 1"},
	}

	tests := make([]test, 40)

	for i := range 40 {
		n := int32(0)
		for n == 0 {
			n = rand.Int32()
		}
		tests[i] = test{
			strconv.FormatInt(int64(n), 10),
			computeZeckendorf(n),
		}
	}

	tests = append(fixedCases, tests...)

	return outputTests(shuffle(tests))
})
