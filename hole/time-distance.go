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

func formatDistance(secs int) string {
	units := []Unit{
		Unit{60 * 60 * 24 * 365, "a year", "years"},
		Unit{60 * 60 * 24 * 30, "a month", "months"},
		Unit{60 * 60 * 24 * 7, "a week", "weeks"},
		Unit{60 * 60 * 24, "a day", "days"},
		Unit{60 * 60, "an hour", "hours"},
		Unit{60, "a minute", "minutes"},
		Unit{1, "a second", "seconds"},
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

	buckets := []int{1, 60, 60 * 60, 60 * 60 * 24, 60 * 60 * 24 * 7, 60 * 60 * 24 * 30, 60 * 60 * 24 * 365}
	tests := [rep * 7 * 4]int{}

	for j, span := range buckets {
		for i := 0; i < rep; i++ {
			tests[j*4*rep+4*i] = randInt(2*span, span*100)    // future plural
			tests[j*4*rep+4*i+1] = randInt(span, span*2)      // future singular
			tests[j*4*rep+4*i+2] = -randInt(2*span, span*100) // past plural
			tests[j*4*rep+4*i+3] = -randInt(span, span*2)     // past singular
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
