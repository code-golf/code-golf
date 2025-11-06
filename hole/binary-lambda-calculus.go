package hole

import (
	"fmt"
	"strings"
)

func blcNum(n int) string {
	return fmt.Sprintf("0000%s10", strings.Repeat("01110", n))
}

var _ = answerFunc("binary-lambda-calculus", func() []Answer {
	tests := fixedTests("binary-lambda-calculus")

	// Addition
	for range 7 {
		i := randInt(1, 7)
		j := randInt(1, 7)
		in := fmt.Sprintf("0101000000000101111101100101111011010%s%s", blcNum(i), blcNum(j))
		out := blcNum(i + j)
		tests = append(tests, test{in, out})
	}

	// Multiplication
	for range 5 {
		i := randInt(1, 5)
		j := randInt(1, 4)
		in := fmt.Sprintf("01010000000111100111010%s%s", blcNum(i), blcNum(j))
		out := blcNum(i * j)
		tests = append(tests, test{in, out})
	}

	// Exponentiation
	i := randInt(2, 3)
	j := randInt(2, 3)
	in := fmt.Sprintf("01%s%s", blcNum(i), blcNum(j))
	outn := 1
	for range i {
		outn *= j
	}
	out := blcNum(outn)
	tests = append(tests, test{in, out})

	return outputTests(shuffle(tests))
})
