package hole

import (
	"fmt"
	"math/rand/v2"
	"strings"
	"time"
)

func generate(s string) string {
	var calendar strings.Builder

	fields := strings.Fields(s)
	month, year := parseInt(fields[0]), parseInt(fields[1])

	calendar.WriteString("Mo Tu We Th Fr Sa Su\n")

	first, total := int((time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC).Weekday()-1+7)%7), time.Date(year, time.Month(month)+1, 0, 0, 0, 0, 0, time.UTC).Day()

	for range first {
		calendar.WriteString("   ")
	}

	for i := 1; i <= total; i++ {
		calendar.WriteString(fmt.Sprintf("%2d ", i))

		if (first+i)%7 == 0 {
			calendar.WriteByte('\n')
		}
	}

	if (first+total)%7 != 0 {
		calendar.WriteByte('\n')
	}

	return calendar.String()
}

func isLeap(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

var _ = answerFunc("calendar", func() []Answer {
	tests := make([]test, 0, 50)

	for i := 0; i < 50; {
		var argument strings.Builder

		month, year := rand.IntN(12), randInt(1800, 2400)

		if i < 12 {
			// Enforce at least one calendar for every month of the year in random non-leap years.
			month = i

			if isLeap(year) {
				continue
			}
		} else if i < 18 {
			// Enforce at least six calendars for months (four random, two February) in random leap years.
			if i < 14 {
				month = 1
			}

			if !isLeap(year) {
				continue
			}
		} else if i < 25 {
			// Enforce February in each century
			month = 1
			year = 1800 + 100*(i-18)
		}

		i++

		argument.WriteString(fmt.Sprintf("%.2d %d", month+1, year))

		tests = append(tests, test{
			argument.String(),
			generate(argument.String()),
		})
	}

	return outputTests(shuffle(tests))
})
