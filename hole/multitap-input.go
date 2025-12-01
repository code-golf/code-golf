
package hole

var _ = answerFunc("multitap-input", func() []Answer {
	tests := make([]test, 100)

	fixedTests := []struct{ words string }{}

	for i, t := range fixedTests {
		tests[i] = multiTapConvert(t.words)
	}

	for i := len(fixedTests); i < len(tests); i++ {
    words := " "
    for i := range randInt(3, 9) {
			if i != 0 {
			  words += " "
	    }
      words += randWord()
		}

		tests[i] = multiTapConvert(words)
	}

	return outputTests(shuffle(tests))
})

func multiTapConvert(text string) test {

	result := ""

	var lastChar byte

	for i, char := range text {
		nextInput := multitapMap[char]
		if i != 0 && (lastChar == nextInput[0] || rand.IntN(4) < 1) {
			result += " "
		}
		result += nextInput
		lastChar = result[len(result)-1]
	}

	return test{result, text}
}


