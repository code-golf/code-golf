package hole

import "time"

var _ = answerFunc("day-of-week", func() []Answer {
	runs := make([][]test, 4)

	for r := range runs {
		tests := fixedTests("day-of-week")

		var years []int
		for year := 1583; year <= 2082; year++ {
			years = append(years, year)
		}
		for year := 2100; year <= 9900; year += 100 {
			years = append(years, year)
		}
		for year := 9901; year <= 9999; year++ {
			years = append(years, year)
		}

		for i, year := range shuffle(years) {
			date := time.Date(year, time.January, i%365+1, 0, 0, 0, 0, time.UTC)
			tests = append(tests, test{
				date.Format(time.DateOnly),
				date.Weekday().String(),
			})
		}

		runs[r] = shuffle(tests)
	}

	return outputTests(runs...)
})
