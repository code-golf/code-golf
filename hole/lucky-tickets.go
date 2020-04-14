package hole

import (
	"math/rand"
	"strconv"
	"strings"
)

type Ticket struct {
	digits int
	base   int
	result int64
}

var data = []Ticket{
	{8, 2, 70},
	{4, 8, 344},
	{2, 10, 10},
	{4, 10, 670},
	{6, 10, 55252},
	{14, 12, 39222848622984},
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
	tickets := make([]Ticket, len(data))
	copy(tickets, data)

	// Randomly generate additional test cases.
	for i := 0; i < 5; i++ {
		digits := 2 + 2*rand.Intn(5)
		base := 2 + rand.Intn(15)

		halfValue := iPow(int64(base), int64(digits/2))
		maxSum := (base - 1) * digits / 2
		counts := make([]int64, maxSum+1)
		var j int64
		for ; j < halfValue; j++ {
			counts[sumDigits(j, int64(base))] += 1
		}

		var result int64
		for _, count := range counts {
			result += count * count
		}

		tickets = append(tickets, Ticket{digits, base, result})
	}

	args := make([]string, len(tickets))
	outs := make([]string, len(tickets))

	for i, item := range tickets {
		args[i] = strconv.Itoa(item.digits) + " " + strconv.Itoa(item.base)
		outs[i] = strconv.FormatInt(item.result, 10)
	}

	// Shuffle
	for i := range args {
		j := rand.Intn(i + 1)
		args[i], args[j] = args[j], args[i]
		outs[i], outs[j] = outs[j], outs[i]
	}

	return args, strings.Join(outs, "\n")
}
