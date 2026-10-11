package hole

import (
	"strconv"
	"strings"
)

const narcissisticMax = 1000

var _ = answerFunc("narcissistic-numbers", narcissisticNumbers)

func narcissisticNumbers() []Answer {
	var out strings.Builder

	for n := 1; n <= narcissisticMax; n++ {
		if isNarcissisticAnyBase(n) {
			if out.Len() > 0 {
				out.WriteByte('\n')
			}
			out.WriteString(strconv.Itoa(n))
		}
	}

	return []Answer{{Args: []string{}, Answer: out.String()}}
}

func isNarcissisticAnyBase(n int) bool {
	for b := 2; b <= n; b++ {
		if isNarcissisticInBase(n, b) {
			return true
		}
	}
	return false
}

func isNarcissisticInBase(n, b int) bool {
	if n <= 0 || b < 2 || b > n {
		return false
	}

	var digits [16]int
	numDigits := 0
	for x := n; x > 0; x /= b {
		digits[numDigits] = x % b
		numDigits++
	}

	sum := 0
	for i := range numDigits {
		sum += intPow(digits[i], numDigits)
		if sum > n {
			return false
		}
	}

	return sum == n
}

func intPow(base, exp int) int {
	result := 1
	for exp > 0 {
		if exp&1 == 1 {
			result *= base
		}
		base *= base
		exp >>= 1
	}
	return result
}
