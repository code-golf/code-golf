package hole

import "strings"

var multitapMap = map[rune]string{
	'A': "2",
	'B': "22",
	'C': "222",
	'D': "3",
	'E': "33",
	'F': "333",
	'G': "4",
	'H': "44",
	'I': "444",
	'J': "5",
	'K': "55",
	'L': "555",
	'M': "6",
	'N': "66",
	'O': "666",
	'P': "7",
	'Q': "77",
	'R': "777",
	'S': "7777",
	'T': "8",
	'U': "88",
	'V': "888",
	'W': "9",
	'X': "99",
	'Y': "999",
	'Z': "9999",
	' ': "0",
}

var _ = answerFunc("multitap-input", func() []Answer {
	tests := make([]test, 50)

	fixedTests := []struct{ words string }{
		{"XYLOPHONE INGREDIENT EACH A VERY"},
		{"JUST ZERO TRAIN OPEN GUIDELINE WE"},
		{"DISAMBIGUATE PUBLIC IF MAKE QUIZ"},
	}

	for i, t := range fixedTests {
		tests[i] = multiTapConvert(t.words)
	}

	for i := len(fixedTests); i < len(tests); i++ {
		text := ""
		for j := range randInt(3, 9) {
			if j != 0 {
				text += " "
			}
			text += randWord()
		}

		tests[i] = multiTapConvert(strings.ToUpper(text))
	}

	return outputTests(shuffle(tests))
})

func multiTapConvert(text string) test {

	result := ""

	var lastChar byte

	for i, char := range text {
		nextInput := multitapMap[char]
		if i != 0 && (lastChar == nextInput[0] || randInt(0, 3) < 1) {
			result += " "
		}
		result += nextInput
		lastChar = result[len(result)-1]
	}

	return test{result, text}
}


