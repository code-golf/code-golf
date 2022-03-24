package hole

import (
	"math/rand"
	"strconv"
)

var (
	roman0 = []string{"", "I", "II", "III", "IV", "V", "VI", "VII", "VIII", "IX"}
	roman1 = []string{"", "X", "XX", "XXX", "XL", "L", "LX", "LXX", "LXXX", "XC"}
	roman2 = []string{"", "C", "CC", "CCC", "CD", "D", "DC", "DCC", "DCCC", "CM"}
	roman3 = []string{"", "M", "MM", "MMM"}
)

func roman(n int) string {
	return roman3[n%10000/1000] + roman2[n%1000/100] + roman1[n%100/10] + roman0[n%10]
}

func arabicToRoman(reverse bool) (args []string, out string) {
	// Start with the special cases.
	numbers := []int{4, 9, 40, 90, 400, 900}

	// Append another 14 random ints.
	for i := 0; i < 14; i++ {
		numbers = append(numbers, rand.Intn(3998)+1) // 1 - 3999 inclusive.
	}

	shuffle(numbers)

	if reverse {
		// Roman to Arabic.
		for _, number := range numbers {
			out += strconv.Itoa(number) + "\n"
			args = append(args, roman(number))
		}
	} else {
		// Arabic to Roman.
		for _, number := range numbers {
			out += roman(number) + "\n"
			args = append(args, strconv.Itoa(number))
		}
	}

	// Drop the trailing newline.
	out = out[:len(out)-1]

	return
}
