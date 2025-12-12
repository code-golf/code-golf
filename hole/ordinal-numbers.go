package hole

import (
	"strconv"

	"github.com/code-golf/code-golf/pretty"
)

var _ = answerFunc("ordinal-numbers", func() []Answer {
	tests := make([]test, 1000)

	for i := range tests {
		tests[i] = test{
			strconv.Itoa(i),
			strconv.Itoa(i) + pretty.Ordinal(i),
		}
	}

	return outputTests(shuffle(tests))
})
