package hole

import (
	"fmt"
	"time"
)

func dayOfWeek() []Scorecard {
	tests := []test{}
	years := []int{}
	for year := 1583; year < 2083; year++ {
		years = append(years, year)
	}
	for year := 2100; year < 9999; year += 100 {
		years = append(years, year)
	}
	for year := 9901; year <= 9999; year++ {
		years = append(years, year)
	}
	for i, y := range shuffle(years) {
		date := time.Date(y, time.January, 1, 0, 0, 0, 0, time.UTC)
		date = date.AddDate(0, 0, i%365)
		tests = append(tests, test{
			date.Format("2006-01-02"),
			fmt.Sprint(date.Weekday()),
		})
	}
	tests = append(tests, test{"1583-01-01", "Saturday"})
	tests = append(tests, test{"1999-12-31", "Friday"})
	tests = append(tests, test{"2000-01-01", "Saturday"})
	tests = append(tests, test{"2000-12-31", "Sunday"})
	tests = append(tests, test{"2001-01-01", "Monday"})
	tests = append(tests, test{"2001-01-02", "Tuesday"})
	tests = append(tests, test{"2001-01-03", "Wednesday"})
	tests = append(tests, test{"2001-01-04", "Thursday"})
	tests = append(tests, test{"9999-12-31", "Friday"})
	return outputTests(shuffle(tests))
}
