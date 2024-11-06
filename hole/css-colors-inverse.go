package hole

func cssColorsInverse() ([]Run, Judge) {
	nameInCodeOutTests := fixedTests("css-colors")
	codeToNames := make(map[string][]string)
	for _, nameInCodeOut := range nameInCodeOutTests {
		codeToNames[nameInCodeOut.out] = append(codeToNames[nameInCodeOut.out], nameInCodeOut.in)
	}
	var tests = []test{}
	for code := range codeToNames {
		tests = append(tests, test{code, ""})
	}
	return outputMultirunTests(tests), oneOfPerOutputJudge(func(code string) []string {
		return codeToNames[code]
	}, true)
}
