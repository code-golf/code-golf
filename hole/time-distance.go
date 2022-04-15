package hole

import (
	"strconv"
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

	inputs := []int{0}

	unitsChosen := []int{1, 2, 3, 4, 5, 6, 7}
	for i := 0; i <= rep; i++ {
		unitsChosen = append(unitsChosen, randInt(1, 7)) // randomize which units will appear
	}

	for _, j := range unitsChosen {
		secs := units[j].seconds
		secsLarger := units[j-1].seconds
		inputs = append(inputs, randInt(secs, secs*2-1))        // future singular
		inputs = append(inputs, -randInt(secs, secs*2-1))       // past singular
		inputs = append(inputs, randInt(2*secs, secsLarger-1))  // future plural
		inputs = append(inputs, -randInt(2*secs, secsLarger-1)) // past plural
		inputs = append(inputs, 2*secs)                         // future exactly 2
		inputs = append(inputs, -2*secs)                        // past exactly 2
		blimit := secs - 1
		if blimit > 1000 {
			blimit = 1000
		}
		a := randInt(2, 6)
		b := randInt(-blimit, blimit)
		inputs = append(inputs, a*secs+b) // future plural antiapproximation
		a = -randInt(2, 6)
		b = randInt(-blimit, blimit)
		inputs = append(inputs, a*secs+b) // past plural antiapproximation
	}

	tests := make([]test, len(inputs))
	for i, inp := range inputs {
		tests[i] = test{strconv.Itoa(inp), formatDistance(inp)}
	}

	return outputTests(shuffle(tests))
}
