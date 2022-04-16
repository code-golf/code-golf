package hole

import (
	"math/rand"
	"strconv"
)

var (
	r0 = [...]string{"", "I", "II", "III", "IV", "V", "VI", "VII", "VIII", "IX"}
	r1 = [...]string{"", "X", "XX", "XXX", "XL", "L", "LX", "LXX", "LXXX", "XC"}
	r2 = [...]string{"", "C", "CC", "CCC", "CD", "D", "DC", "DCC", "DCCC", "CM"}
	r3 = [...]string{"", "M", "MM", "MMM"}
)

func arabicToRoman(reverse bool) ([]string, string) {
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

	return outputTests(tests)
}
