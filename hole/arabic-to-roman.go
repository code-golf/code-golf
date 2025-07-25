package hole

import (
	"math/rand/v2"
	"strconv"
)

var r0 = [...]string{"", "I", "II", "III", "IV", "V", "VI", "VII", "VIII", "IX"}
var r1 = [...]string{"", "X", "XX", "XXX", "XL", "L", "LX", "LXX", "LXXX", "XC"}
var r2 = [...]string{"", "C", "CC", "CCC", "CD", "D", "DC", "DCC", "DCCC", "CM"}
var r3 = [...]string{"", "M", "MM", "MMM"}

var _ = answerFunc("arabic-to-roman", func() []Answer { return arabicToRoman(false) })
var _ = answerFunc("roman-to-arabic", func() []Answer { return arabicToRoman(true) })

func arabicToRoman(reverse bool) []Answer {
	// The max roman numeral is 3,999. Test all of them.
	const count = 3999

	tests := make([]test, count)
	for i, n := range rand.Perm(count) {
		n++
		arabic := strconv.Itoa(n)
		roman := r3[n%10000/1000] + r2[n%1000/100] + r1[n%100/10] + r0[n%10]

		if reverse {
			tests[i] = test{roman, arabic}
		} else {
			tests[i] = test{arabic, roman}
		}
	}

	const argc = 2000 // Preserve original argc
	return outputTests(tests[:argc], tests[len(tests)-argc:])
}
