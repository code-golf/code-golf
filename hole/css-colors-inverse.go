package hole

var cssFixedTests []test
var cssCodesToNames = make(map[string][]string)

func initCSS() {
	if len(cssCodesToNames) < 1 {
		cssFixedTests = fixedTests("css-colors")
		for _, nameInCodeOut := range cssFixedTests {
			cssCodesToNames[nameInCodeOut.out] = append(cssCodesToNames[nameInCodeOut.out], nameInCodeOut.in)
		}
	}
}

var _ = answerFunc("css-colors-inverse", func() []Answer {
	initCSS()
	var tests = []test{}
	for code := range cssCodesToNames {
		tests = append(tests, test{code, ""})
	}
	return outputMultirunTests(tests)
})

var _ = judge("css-colors-inverse", oneOfPerOutputJudge(func(code string) []string {
	initCSS()
	return cssCodesToNames[code]
}, true))
