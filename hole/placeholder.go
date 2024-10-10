package hole

import "strconv"

func placeholder() []Run {
	runs := make([][]test, 5)

	for r := range runs {
		var tests []test

		numArgs := randInt(1, 100)
		for i := 0; i < numArgs; i++ {
			word := randWord()
			tests = append(tests, test{word, word})
		}

		fixedOut := "Hello, World!\n"
		for i := 0; i < 10; i++ {
			fixedOut += strconv.Itoa(i) + "\n"
		}
		tests[0].out = fixedOut + tests[0].out

		runs[r] = tests
	}

	return outputTests(runs...)
}
