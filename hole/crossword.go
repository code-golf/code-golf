package hole

func init() {
	crosswords := make(map[string][]string)

	for _, words := range fixedTests("crossword") {
		crosswords[words.in] = append(crosswords[words.in], words.out)
	}

	answerFunc("crossword", func() []Answer {
		var tests []test

		for crossword := range crosswords {
			tests = append(tests, test{crossword, ""})
		}

		return outputMultirunTests(tests)
	})

	judge("crossword", oneOfPerOutputJudge(func(words string) []string {
		return crosswords[words]
	}, false))
}
