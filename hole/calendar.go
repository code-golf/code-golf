package hole

import (
	"fmt"
	"math/rand/v2"
	"strings"
	"time"
)

func generate(month, year int) string {
	first := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	days := first.AddDate(0, 1, -1).Day()
	offset := (int(first.Weekday()) + 6) % 7

	var calendar strings.Builder
	calendar.WriteString("Mo Tu We Th Fr Sa Su\n")
	calendar.WriteString(strings.Repeat("   ", offset))

	for day := 1; day <= days; day++ {
		fmt.Fprintf(&calendar, "%2d ", day)

		if (day+offset)%7 == 0 || day == days {
			calendar.WriteByte('\n')
		}
	}

	return calendar.String()
}

func isLeap(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

var _ = answerFunc("calendar", func() []Answer {
	tests := make([]test, 50)

	for i := 0; i < len(tests); {
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

		month++

		tests[i] = test{
			fmt.Sprintf("%.2d %d", month, year),
			generate(month, year),
		}

		i++
	}

	return outputTests(shuffle(tests))
})
