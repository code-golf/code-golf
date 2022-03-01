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
	const rep = 32

	tests := []int{0}

	unitsChosen := []int{1, 2, 3, 4, 5, 6, 7}
	for i := 0; i <= rep; i++ {
		unitsChosen = append(unitsChosen, randInt(1, 7)) // randomize which units will appear
	}

	for _, j := range unitsChosen {
		secs := units[j].seconds
		secsLarger := units[j-1].seconds
		tests = append(tests, randInt(secs, secs*2-1))        // future singular
		tests = append(tests, -randInt(secs, secs*2-1))       // past singular
		tests = append(tests, randInt(2*secs, secsLarger-1))  // future plural
		tests = append(tests, -randInt(2*secs, secsLarger-1)) // past plural
		blimit := secs - 1
		if blimit > 1000 {
			blimit = 1000
		}
		a := randInt(2, 6)
		b := randInt(-blimit, blimit)
		tests = append(tests, a*secs+b) // future plural antiapproximation
		a = -randInt(2, 6)
		b = randInt(-blimit, blimit)
		tests = append(tests, a*secs+b) // past plural antiapproximation
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
