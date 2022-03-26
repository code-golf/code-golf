package hole

import (
	"strconv"

	"github.com/code-golf/code-golf/pretty"
)

func ordinalNumbers() ([]string, string) {
	const count = 1000
	tests := make([]test, count)

	for i := 0; i < count; i++ {
		tests[i] = test{
			strconv.Itoa(i),
			strconv.Itoa(i) + pretty.Ordinal(i),
		}
	}

	return outputTests(shuffle(tests))
}
