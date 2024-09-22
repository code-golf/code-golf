package hole

import (
	"strconv"
	"strings"
)

func cardNumberValidation() []Run {
	type Case struct {
		arg   string
		valid bool
	}

	digitSum := func(n int) int {
		return n - n/10*9
	}

	getCheckDigit := func(digits []int) int {
		s := 0
		for i, x := range digits {
			s += digitSum(x + (len(digits)-i)%2*x)
		}
		return 9 - (s+9)%10
	}

	formatDigits := func(digits []int) string {
		result := ""
		for i, x := range digits {
			if i > 0 && i%4 < 1 {
				result += " "
			}
			result += strconv.Itoa(x)
		}
		return result
	}

	cases := []Case{
		{"4242 4242 4242 4242", true},
		{"4242 4244 4242 4242", false},
		{"5555 5555 5555 4444", true},
		{"5555 5555 5555 5444", false},
		{"3566 0020 2036 0505", true},
		{"3656 0020 2036 0505", false},
	}

	for i := range 100 {
		digits := make([]int, 15)
		for j := range digits {
			digits[j] = randInt(0, 9)
		}
		checkDigit := getCheckDigit(digits)
		lastDigit := checkDigit
		if i < 50 {
			lastDigit = randInt(0, 9)
		}
		cases = append(cases, Case{formatDigits(append(digits, lastDigit)), lastDigit == checkDigit})
	}

	shuffle(cases)
	args := make([]string, len(cases))
	var answer strings.Builder

	for i, c := range cases {
		args[i] = c.arg

		if c.valid {
			if answer.Len() > 0 {
				answer.WriteByte('\n')
			}
			answer.WriteString(c.arg)
		}
	}

	return []Run{{Args: args, Answer: answer.String()}}
}
