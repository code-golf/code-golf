package hole

import (
	"math/rand"
	"strconv"
	"strings"
)

type Unit struct {
	seconds  int
	singular string
	plural   string
}

var units = []Unit{
	{60 * 60 * 24 * 365 * 1000, "a millenium", "millenia"},
	{60 * 60 * 24 * 365, "a year", "years"},
	{60 * 60 * 24 * 30, "a month", "months"},
	{60 * 60 * 24 * 7, "a week", "weeks"},
	{60 * 60 * 24, "a day", "days"},
	{60 * 60, "an hour", "hours"},
	{60, "a minute", "minutes"},
	{1, "a second", "seconds"},
}

func formatDistance(secs int) string {
	if secs == 0 {
		return "now"
	}
	past := secs < 0
	if past {
		secs = -secs
	}
	result := ""
	for _, v := range units {
		if v.seconds <= secs {
			q := secs / v.seconds
			if q == 1 {
				result = v.singular
			} else {
				result = strconv.Itoa(q) + " " + v.plural
			}
			break
		}
	}
	if past {
		return result + " ago"
	}
	return "in " + result
}

func timeDistance() ([]string, string) {
	const rep = 2

	tests := []int{0}

	for j, v := range units {
		if j == 0 {
			continue
		}
		tests = append(tests, randInt(v.seconds, v.seconds*2-1))  // future singular
		tests = append(tests, -randInt(v.seconds, v.seconds*2-1)) // past singular
		for i := 0; i < rep; i++ {
			tests = append(tests, randInt(2*v.seconds, units[j-1].seconds-1))  // future plural
			tests = append(tests, -randInt(2*v.seconds, units[j-1].seconds-1)) // past plural
		}
	}

	rand.Shuffle(len(tests), func(i, j int) {
		tests[i], tests[j] = tests[j], tests[i]
	})

	args := make([]string, len(tests))
	var answer strings.Builder

	for i, secs := range tests {
		args[i] = strconv.Itoa(secs)
		if i > 0 {
			answer.WriteByte('\n')
		}
		answer.WriteString(formatDistance(secs))
	}

	return args, answer.String()
}
