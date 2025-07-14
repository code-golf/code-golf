package hole

import (
	"fmt"
	"strings"
)

var _ = answerFunc("billiards", func() []Answer {
	tests := []test{}
	addTest := func(height, width int) {
		common := gcd(height, width)
		patternSize := common * 2
		var answer strings.Builder
		for y := range height {
			if y > 0 {
				answer.WriteByte('\n')
			}
			for x := range width {
				if x%patternSize == y%patternSize {
					answer.WriteByte('\\')
				} else if x%patternSize == patternSize-1-y%patternSize {
					answer.WriteByte('/')
				} else {
					answer.WriteByte(' ')
				}
			}
		}
		tests = append(tests, test{fmt.Sprint(height, width), answer.String()})
	}
	for i := range 32 {
		common := i%5 + 1
		for {
			height := common * randInt(1, 2)
			width := common * randInt(1, 5)
			if gcd(height, width) == common {
				addTest(height, width)
				break
			}
		}
	}
	for i := range 32 {
		height := i%8 + 1
		width := randInt(1, 20)
		addTest(height, width)
	}

	return outputTestsWithSep("\n\n", shuffle(tests))
})
