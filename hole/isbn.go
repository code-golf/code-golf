package hole

import (
	"math/rand"
	"strings"
)

func checkDigit(digits [9]int) int {
	sum, weight := 0, 10

	for _, digit := range digits {
		sum += digit * weight
		weight--
	}

	return (11 - (sum % 11)) % 11
}

func isbn() []Run {
	args := make([]string, 100)
	outs := make([]string, 100)

	for i, perm := range rand.Perm(100) {
		var digits [9]int

		for j := range digits {
			digits[j] = rand.Intn(10)
		}

		// Guarantee at least 5 arguments end with 'X'
		if perm < 5 {
			if digits[7] == 9 {
				digits[7] = 8
			}

			for checkDigit(digits) != 2 && checkDigit(digits) != 10 {
				digits[8] = (digits[8] + 1) % 10
			}

			if checkDigit(digits) != 10 {
				digits[7]++
			}
		}

		var id strings.Builder

		// This here logic is for varying the second two parts of the ISBN.
		// Sure, it's cosmetic, but it might mess some people up.
		difference := 7 - rand.Intn(5)
		for j, digit := range digits {
			id.WriteByte(byte('0' + digit))

			if j == 0 || j == difference {
				id.WriteByte('-')
			}
		}

		id.WriteByte('-')

		args[i] = id.String()

		if digit := checkDigit(digits); digit == 10 {
			id.WriteByte('X')
		} else {
			id.WriteByte(byte('0' + digit))
		}

		outs[i] = id.String()
	}

	return []Run{{Args: args, Answer: strings.Join(outs, "\n")}}
}
