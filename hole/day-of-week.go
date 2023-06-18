package hole

import "time"

func dayOfWeek() []Run {
	tests := []test{
		{"1583-01-01", "Saturday"}, {"1999-12-31", "Friday"},
		{"2000-01-01", "Saturday"}, {"2000-12-31", "Sunday"},
		{"2001-01-01", "Monday"}, {"2001-01-02", "Tuesday"},
		{"2001-01-03", "Wednesday"}, {"2001-01-04", "Thursday"},
		{"9999-12-31", "Friday"},
	}

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
		date := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
		date = date.AddDate(0, 0, i%365)
		tests = append(tests, test{
			date.Format(time.DateOnly),
			date.Weekday().String(),
		})
	}

	return outputTests(shuffle(tests))
}
