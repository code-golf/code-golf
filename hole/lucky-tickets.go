package hole

import (
	"fmt"
	"math/rand/v2"
)

type ticket struct {
	digits, base int
	result       int64
}

var data = [...]ticket{
	{2, 10, 10},
	{4, 8, 344},
	{4, 9, 489},
	{4, 10, 670},
	{4, 15, 2255},
	{4, 16, 2736},
	{6, 6, 4332},
	{6, 10, 55252},
	{6, 13, 204763},
	{6, 15, 418503},
	{8, 2, 70},
	{8, 9, 2306025},
	{8, 10, 4816030},
	{8, 15, 82073295},
	{10, 15, 16581420835},
	{12, 7, 786588243},
	{12, 9, 12434998005},
	{12, 10, 39581170420},
	{12, 11, 112835748609},
	{14, 5, 454805755},
	{14, 7, 35751527189},
	{14, 12, 39222848622984},
}

func iPow(a, b int64) int64 {
	result := int64(1)

	for b != 0 {
		if b&1 != 0 {
			result *= a
		}
		b >>= 1
		a *= a
	}

	return result
}

func sumDigits(number, base int64) (result int64) {
	for number > 0 {
		result += number % base
		number /= base
	}

	return result
}

func luckyTickets() []Run {
	var tickets [40]ticket

	// Add all fixed test cases.
	for i, j := range rand.Perm(len(data)) {
		tickets[i] = data[j]
	}

	// Randomly generate additional test cases.
	for i := 22; i < 40; i++ {
		digits := 2 + 2*rand.IntN(5)
		base := 2 + rand.IntN(15)

		halfValue := iPow(int64(base), int64(digits/2))
		maxSum := (base - 1) * digits / 2
		counts := make([]int64, maxSum+1)
		for j := int64(0); j < halfValue; j++ {
			counts[sumDigits(j, int64(base))]++
		}

		var result int64
		for _, count := range counts {
			result += count * count
		}

		tickets[i] = ticket{digits, base, result}
	}

	tests := make([]test, len(tickets))

	for i, item := range tickets {
		tests[i] = test{
			fmt.Sprint(item.digits, item.base),
			fmt.Sprint(item.result),
		}
	}

	tests = shuffle(tests)

	const argc = 20 // Preserve original argc
	return outputTests(tests[:argc], tests[len(tests)-argc:])
}
