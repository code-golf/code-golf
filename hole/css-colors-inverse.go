package hole

var _ = answerFunc("css-colors-inverse", func() []Answer {
	nameInCodeOutTests := fixedTests("css-colors")
	codeToNames := make(map[string][]string)
	for _, nameInCodeOut := range nameInCodeOutTests {
		codeToNames[nameInCodeOut.out] = append(codeToNames[nameInCodeOut.out], nameInCodeOut.in)
	}
	var tests = []test{}
	for code := range codeToNames {
		tests = append(tests, test{code, ""})
	}
	return outputMultirunTests(tests)
})

func getCSSColorsInverseJudge() Judge {
	nameInCodeOutTests := fixedTests("css-colors")
	codeToNames := make(map[string][]string)
	for _, nameInCodeOut := range nameInCodeOutTests {
		codeToNames[nameInCodeOut.out] = append(codeToNames[nameInCodeOut.out], nameInCodeOut.in)
	}
	return oneOfPerOutputJudge(func(code string) []string {
		return codeToNames[code]
	}, true)
}

var _ = judge("css-colors-inverse", getCSSColorsInverseJudge())
