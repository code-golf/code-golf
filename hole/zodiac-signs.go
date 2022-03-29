package hole

import (
	"fmt"
	"math/rand"
)

var signs = [...]zodiacSign{
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

type zodiacSign struct {
	month, startDate, endDate int
	symbol                    string
}

func (sign zodiacSign) randomDate() test {
	d := rand.Intn(sign.endDate-sign.startDate) + sign.startDate
	return test{fmt.Sprintf("%02d-%02d", sign.month, d), sign.symbol}
}

func (sign zodiacSign) edgeDate() test {
	d := sign.startDate
	if d == 1 {
		d = sign.endDate
	}
	return test{fmt.Sprintf("%02d-%02d", sign.month, d), sign.symbol}
}

func zodiacSigns() ([]string, string) {
	const (
		randomCases = 20
		totalcases  = randomCases + 2*len(signs)
	)

	tests := make([]test, randomCases, totalcases)

	for i := 0; i < randomCases; i++ {
		tests[i] = signs[rand.Intn(len(signs))].randomDate()
	}

	for _, sign := range signs {
		tests = append(tests, sign.edgeDate(), sign.randomDate())
	}

	return outputTests(shuffle(tests))
}
