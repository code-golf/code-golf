package hole

var _ = answerFunc("quine", func() []Answer {
	return []Answer{{Args: []string{}, Answer: ""}}
})

var _ = judge("quine", func(run Run) string {
	return run.Code
})
