package hole

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

type Unit struct {
	seconds  int
	singular string
	plural   string
}

func randInt(a, b int) int { return rand.Intn(b-a+1) + a }

func formatDistance(secs int) string {
	units := []Unit{
		Unit{seconds: 60 * 60 * 24 * 365, singular: "a year", plural: "years"},
		Unit{seconds: 60 * 60 * 24 * 30, singular: "a month", plural: "months"},
		Unit{seconds: 60 * 60 * 24 * 7, singular: "a week", plural: "weeks"},
		Unit{seconds: 60 * 60 * 24, singular: "a day", plural: "days"},
		Unit{seconds: 60 * 60, singular: "an hour", plural: "hours"},
		Unit{seconds: 60, singular: "a minute", plural: "minutes"},
		Unit{seconds: 1, singular: "a second", plural: "seconds"},
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
