package hole

import (
	"fmt"
	"math/rand"
	"strings"
)

type ZodiacSign struct {
	Month, StartDate, EndDate int
	Symbol                    string
}
type Test struct {
	Input, Output string
}

func RandomDate(sign ZodiacSign) Test {
	d := rand.Intn(sign.EndDate-sign.StartDate) + sign.StartDate
	return Test{fmt.Sprintf("%02d", sign.Month) + "-" + fmt.Sprintf("%02d", d), sign.Symbol}
}

func EdgeDate(sign ZodiacSign) Test {
	d := sign.StartDate
	if d == 1 {
		d = sign.EndDate
	}
	return Test{fmt.Sprintf("%02d", sign.Month) + "-" + fmt.Sprintf("%02d", d), sign.Symbol}
}

func zodiacSigns() ([]string, string) {
	signs := []ZodiacSign{
		{3, 21, 31, "♈"},
		{4, 1, 19, "♈"},

		{4, 20, 30, "♉"},
		{5, 1, 20, "♉"},

		{5, 21, 31, "♊"},
		{6, 1, 21, "♊"},

		{6, 22, 30, "♋"},
		{7, 1, 22, "♋"},

		{7, 23, 31, "♌"},
		{8, 1, 22, "♌"},

		{8, 23, 31, "♍"},
		{9, 1, 22, "♍"},

		{9, 23, 30, "♎"},
		{10, 1, 22, "♎"},

		{10, 23, 31, "♏"},
		{11, 1, 22, "♏"},

		{11, 23, 30, "♐"},
		{12, 1, 21, "♐"},

		{12, 22, 31, "♑"},
		{1, 1, 19, "♑"},

		{1, 20, 31, "♒"},
		{2, 1, 18, "♒"},

		{2, 19, 28, "♓"},
		{3, 1, 20, "♓"},
	}

	tests := []Test{}
	for _, sign := range signs {
		tests = append(tests, RandomDate(sign))
		tests = append(tests, EdgeDate(sign))
	}
	for i := 0; i < 20; i++ {
		sign := signs[rand.Intn(24)]
		tests = append(tests, RandomDate(sign))
	}

	rand.Shuffle(len(tests), func(i, j int) {
		tests[i], tests[j] = tests[j], tests[i]
	})

	args := make([]string, len(tests))
	var answer strings.Builder

	for i, x := range tests {
		args[i] = x.Input
		if i > 0 {
			answer.WriteByte('\n')
		}
		answer.WriteString(x.Output)
	}

	return args, answer.String()
}
