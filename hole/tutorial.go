package hole

import "strconv"

var _ = answerFunc("tutorial", func() []Answer {
	runs := make([][]test, 5)

	for r := range runs {
		var tests []test

		for range randInt(1, 99) {
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
})
