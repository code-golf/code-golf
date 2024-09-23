package hole

import "strconv"

func helloWorld() []Run {
	var tests []test

	for i := 0; i < 100; i++ {
		word := randWord()
		tests = append(tests, test{word, word})
	}

	firstArgOut := tests[0].out

	fixedOut := "Hello, World!\n"
	for i := 0; i < 10; i++ {
		fixedOut += strconv.Itoa(i) + "\n"
	}
	tests[0].out = fixedOut + firstArgOut

	return outputTests(tests)
}
