package hole

import (
	"math/rand"
	"strconv"
	"strings"
)

type ticket struct {
	digits, base int
	result       int64
}

var data = []ticket{
	{8, 2, 70},
	{14, 5, 454805755},
	{6, 6, 4332},
	{12, 7, 786588243},
	{4, 8, 344},
	{4, 9, 489},
	{8, 9, 2306025},
	{2, 10, 10},
	{4, 10, 670},
	{6, 10, 55252},
	{8, 10, 4816030},
	{12, 9, 12434998005},
	{12, 10, 39581170420},
	{12, 11, 112835748609},
	{6, 13, 204763},
	{4, 15, 2255},
	{6, 15, 418503},
	{8, 15, 82073295},
	{10, 15, 16581420835},
}

func iPow(a, b int64) int64 {
	var result int64 = 1

	for b != 0 {
		if b&1 != 0 {
			result *= a
		}
		b >>= 1
		a *= a
	}

	return result
}

func sumDigits(number int64, base int64) int64 {
	var result int64

	for number > 0 {
		result += number % base
		number /= base
	}

	return result
}

func luckyTickets() ([]string, string) {
	// make a random selection of 4 from the fixed cases
	rand.Shuffle(len(data), func(i, j int) {
		data[i], data[j] = data[j], data[i]
	})
	data = data[0:4]
	tickets := make([]ticket, len(data))
	// always add case 14 12
	tickets = append(tickets, ticket{14, 12, 39222848622984})
	copy(tickets, data)

	// Randomly generate additional test cases.
	for i := 0; i < 15; i++ {
		digits := 2 + 2*rand.Intn(5)
		base := 2 + rand.Intn(15)

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

		tickets = append(tickets, ticket{digits, base, result})
	}

	args := make([]string, len(tickets))
	outs := make([]string, len(tickets))

	for i, item := range tickets {
		args[i] = strconv.Itoa(item.digits) + " " + strconv.Itoa(item.base)
		outs[i] = strconv.FormatInt(item.result, 10)
	}

	rand.Shuffle(len(args), func(i, j int) {
		args[i], args[j] = args[j], args[i]
		outs[i], outs[j] = outs[j], outs[i]
	})

	return args, strings.Join(outs, "\n")
}
