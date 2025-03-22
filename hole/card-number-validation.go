package hole

import "strconv"

var _ = answerFunc("card-number-validation", func() []Answer {
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

	tests := []test{
		{"4242 4242 4242 4242", "4242 4242 4242 4242"},
		{"4242 4244 4242 4242", ""},
		{"5555 5555 5555 4444", "5555 5555 5555 4444"},
		{"5555 5555 5555 5444", ""},
		{"3566 0020 2036 0505", "3566 0020 2036 0505"},
		{"3656 0020 2036 0505", ""},
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

		cardNumber := formatDigits(append(digits, lastDigit))
		t := test{in: cardNumber}
		if lastDigit == checkDigit {
			t.out = cardNumber
		}

		tests = append(tests, t)
	}

	return outputTests(shuffle(tests))
})
