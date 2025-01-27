package hole

import "strconv"

func tutorial() []Run {
	runs := make([][]test, 5)

	for r := range runs {
		var tests []test

		for range randInt(1, 100) {
			word := randWord()
			tests = append(tests, test{word, word})
		}

		fixedOut := "Hello, World!\n"
		for i := range 10 {
			fixedOut += strconv.Itoa(i) + "\n"
		}
		tests[0].out = fixedOut + tests[0].out

		runs[r] = tests
	}

	return outputTests(runs...)
}
