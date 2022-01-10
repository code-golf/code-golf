package hole

import (
	"math/rand"
	"strconv"
)

func check_digit(digits [9]int) int {
	sum, weight := 0, 10

	for _, digit := range digits {
		sum += digit * weight
		weight--
	}

	return (11 - (sum % 11)) % 11
}

func isbn() (args []string, out string) {
	guaranteedTens := rand.Perm(100)
	for i := 0; i < 100; i++ {
		var digits [9]int

		for j := 0; j < 9; j++ {
			digits[j] = rand.Intn(10)
		}

		// Guarantee at least 5 arguments end with 'X'
		if guaranteedTens[i] < 5 {
			if digits[7] == 9 {
				digits[7] = 8
			}

			for check_digit(digits) != 2 && check_digit(digits) != 10 {
				digits[8] = (digits[8] + 1) % 10
			}

			if check_digit(digits) != 10 {
				digits[7] += 1
			}
		}

		arg := ""

		// This here logic is for varying the second two parts of the ISBN.
		// Sure, it's cosmetic, but it might mess some people up.
		difference := 7 - rand.Intn(5)
		for j, digit := range digits {
			arg += strconv.Itoa(digit)

			if j == 0 || j == difference {
				arg += "-"
			}
		}

		arg += "-"
		args = append(args, arg)

		if check := check_digit(digits); check == 10 {
			out += arg + "X\n"
		} else {
			out += arg + strconv.Itoa(check) + "\n"
		}
	}

	// Trim the trailing newline.
	out = out[:len(out)-1]

	return
}
