package hole

import (
	"fmt"
	"strings"
	"time"
)

func generate(s string) string {
	var calendar strings.Builder

	fields := strings.Fields(s)
	month, year := parseInt(fields[0]), parseInt(fields[1])

	calendar.WriteString("Mo Tu We Th Fr Sa Su\n")

	first, total := int((time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC).Weekday()-1+7)%7), time.Date(year, time.Month(month)+1, 0, 0, 0, 0, 0, time.UTC).Day()

	for i := 0; i < first; i++ {
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

func calendar() []Run {
	tests := make([]test, 0, 31)

	for range 31 {
		argument := fmt.Sprintf("%.2d %d", randInt(1, 12), randInt(1800, 2400))

		tests = append(tests, test{
			argument,
			generate(argument),
		})
	}

	return outputTests(shuffle(tests))
}
