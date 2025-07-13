package hole

var _ = answerFunc("palindromemordnilap", func() []Answer {
	const alphabet = "qwertzuiopasdfghjklyxcvbnmQWERTZUIOPASDFGHJKLYXCVBNM0123456789"

	correctLength := func(input string) int {
		for length := len(input); length > 1; length-- {
			for i := (length + 1) / 2; i >= 0; i-- {
				if input[len(input)-1-i] != input[len(input)-length+i] {
					break
				}
				if i == 0 {
					return length
				}
			}
		}
		return 1
	}

	solve := func(input string) string {
		for i := len(input) - correctLength(input) - 1; i >= 0; i-- {
			input += input[i : i+1]
		}
		return input
	}

	fixedInputs := []string{
		"123456", "8989a",

		"a", "aA", "aa", "aaaaaaa", "ab", "aba", "abaaaba", "abb", "abc",
		"abca", "abcdc", "abcdcc", "sesphase",

		"better", "mississippi", "Palindrome",

		"ababcdcdefefg", "ABABCDCDEFEFG",
		"ghghijijklklm", "GHGHIJIJKLKLM",
		"mnmnopopqrqrs", "MNMNOPOPQRQRS",
		"ststuvuvwxwxy", "STSTUVUVWXWXY",
		"yzyz010123234", "YZYZ454567678",
	}

	tests := make([]test, len(fixedInputs))
	for i, input := range fixedInputs {
		tests[i] = test{input, solve(input)}
	}

	for baseLength := 1; baseLength <= 8; baseLength++ {
		for endIndex := baseLength - 1; endIndex <= baseLength; endIndex++ {
			for startIndex := 0; startIndex <= endIndex; startIndex++ {
				input := ""
				for range baseLength {
					j := randInt(0, 61)
					input += alphabet[j : j+1]
				}
				for i := endIndex - 1; i >= startIndex; i-- {
					input += input[i : i+1]
				}
				tests = append(tests, test{input, solve(input)})
			}
		}
	}

	tests = tests[:106] // Preserve original argc

	return outputTests(shuffle(tests))
})
