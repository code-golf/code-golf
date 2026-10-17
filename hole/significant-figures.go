package hole

import (
	"fmt"
	"math/rand/v2"
	"strings"
)

var _ = answerFunc("significant-figures", func() []Answer {
	tests := make([]test, 100)

	tests[0] = test{"48721", "5"}
	tests[1] = test{"1", "1"}
	tests[2] = test{"10.0", "3"}
	tests[3] = test{"830.970", "6"}
	tests[4] = test{"12321", "5"}

	for i := 5; i < len(tests); i++ {
		mantissalength := rand.IntN(16) + 1
		mantissa := ""

		x := rand.IntN(9) + 1
		mantissa += fmt.Sprint(x)

		for len(mantissa) < mantissalength-1 {
			x := rand.IntN(10)
			mantissa += fmt.Sprint(x)
		}

		for len(mantissa) < mantissalength {
			x := rand.IntN(9) + 1
			mantissa += fmt.Sprint(x)
		}

		var fullNumber string

		if i < 20 {
			fullNumber = mantissa + strings.Repeat("0", rand.IntN(16))
		} else if i < 40 {
			trailingZeroes := rand.IntN(15) + 1
			fullNumber = mantissa + strings.Repeat("0", trailingZeroes) + "."
			mantissalength += trailingZeroes
		} else if i < 60 {
			trailingAfterDecimal := rand.IntN(15) + 1
			trailingBeforeDecimal := rand.IntN(16 - trailingAfterDecimal)
			fullNumber = mantissa + strings.Repeat("0", trailingBeforeDecimal) + "." + strings.Repeat("0", trailingAfterDecimal)
			mantissalength += trailingBeforeDecimal + trailingAfterDecimal
		} else {
			magnitude := rand.IntN(14+mantissalength) - 16
			trailingZeroes := rand.IntN(8)
			if magnitude < 0 {
				fullNumber = "0." + strings.Repeat("0", -1-magnitude) + mantissa + strings.Repeat("0", trailingZeroes)
			} else {
				fullNumber = mantissa[:mantissalength-1-magnitude] + "." + mantissa[mantissalength-1-magnitude:] + strings.Repeat("0", trailingZeroes)
			}
			mantissalength += trailingZeroes
		}

		tests[i] = test{fullNumber, fmt.Sprint(mantissalength)}
	}

	return outputTests(shuffle(tests))
})
