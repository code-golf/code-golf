package hole

import "strconv"

var (
	r0 = [...]string{"", "I", "II", "III", "IV", "V", "VI", "VII", "VIII", "IX"}
	r1 = [...]string{"", "X", "XX", "XXX", "XL", "L", "LX", "LXX", "LXXX", "XC"}
	r2 = [...]string{"", "C", "CC", "CCC", "CD", "D", "DC", "DCC", "DCCC", "CM"}
	r3 = [...]string{"", "M", "MM", "MMM"}
)

func arabicToRoman(reverse bool) []Scorecard {
	// Testing all ~4k is too slow and is too many arguments for J.
	const count = 2000

	// Start with the 6 special cases.
	numbers := make([]int, count)
	numbers[0] = 4
	numbers[1] = 9
	numbers[2] = 40
	numbers[3] = 90
	numbers[4] = 400
	numbers[5] = 900

	// Append another 1,994 random cases.
	for i := 6; i < count; i++ {
		numbers[i] = randInt(1, 3999)
	}

	tests := make([]test, count)
	for i, n := range shuffle(numbers) {
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
