package hole

var _ = answerFunc("floyd-steinberg-dithering", func() []Answer {
	tests := fixedTests("floyd-steinberg-dithering")
	runs := make([]Answer, len(tests))
	for i, test := range tests {
		runs[i] = Answer{Args: []string{test.in}, Answer: test.out}
	}
	return shuffle(runs)
})
