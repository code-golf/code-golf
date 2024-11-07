package hole

import (
	"strconv"

	"github.com/code-golf/code-golf/pretty"
)

func ordinalNumbers() []Run {
	tests := make([]test, 1000)

	for i := range tests {
		tests[i] = test{
			strconv.Itoa(i),
			strconv.Itoa(i) + pretty.Ordinal(i),
		}
	}

	return outputTests(shuffle(tests))
}
