package hole

func init() {
	cssCodesToNames := make(map[string][]string)
	for _, nameInCodeOut := range fixedTests("css-colors") {
		cssCodesToNames[nameInCodeOut.out] = append(cssCodesToNames[nameInCodeOut.out], nameInCodeOut.in)
	}

	answerFunc("css-colors-inverse", func() []Answer {
		var tests = []test{}
		for code := range cssCodesToNames {
			tests = append(tests, test{code, ""})
		}
		return outputMultirunTests(tests)
	})

	judge("css-colors-inverse", oneOfPerOutputJudge(func(code string) []string {
		return cssCodesToNames[code]
	}, true))
}
