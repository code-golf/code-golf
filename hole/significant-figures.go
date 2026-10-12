package hole

import (
	"fmt"
	"math/rand/v2"
	"strings"
)

var _ = answerFunc("significant-figures", func() []Answer {
	tests := make([]test, 100)

	for i := 0; i < len(tests); i++ {
		mantissalength := rand.IntN(16) + 1
		mantissa := ""

		x := rand.IntN(9) + 1
		mantissa += fmt.Sprint(x)

		for len(mantissa) < mantissalength - 1{
			x := rand.IntN(10)
			mantissa += fmt.Sprint(x)
		}

		for len(mantissa) < mantissalength{
			x := rand.IntN(9) + 1
			mantissa += fmt.Sprint(x)
		}

		var fullnumber string

		if i < 20 {
			fullnumber := mantissa + strings.Repeat("0", rand.IntN(16))
		} else if i < 40 {
			trailingZeroes := rand.IntN(15) + 1
			fullNumber := mantissa + strings.Repeat("0", trailingZeroes) + "."
			mantissalength += trailingZeroes
		} else if i < 60 {
			trailingAfterDecimal := rand.IntN(15) + 1
			trailingBeforeDecimal := rand.IntN(16 - trailingAfterDecimal)
			fullNumber := mantissa + strings.Repeat("0", trailingBeforeDecimal) + "." + strings.Repeat("0", trailingAfterDecimal)
			mantissalength += trailingBeforeDecimal + trailingAfterDecimal
		} else {
			magnitude := rand.IntN(14 + mantissalength) - 16
			trailingZeroes := rand.IntN(16)
			// todo
			mantissalength += trailingZeroes
		}

		tests[i] = tests{fullNumber, fmt.Sprint(mantissalength)}
	}
})
