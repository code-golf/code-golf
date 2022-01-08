package hole

import (
	"math/rand"
	"strconv"
)

func isbn() (args []string, out string) {
	for i := 0; i < 20; i++ {
		weightedDigitsSum := 0
		weight := 10

		// First digit of ISBN, not sticking with traditional 1 or 0, can't
		// let them exploit that.
		firstDigit := rand.Intn(10)

		weightedDigitsSum += firstDigit * weight
		weight--

		arg := strconv.Itoa(firstDigit) + "-"

		// This here logic is for varying the second two parts of the ISBN.
		// Sure, it's cosmetic, but it might mess some people up.
		difference := 6 - rand.Intn(5)
		for i := 0; i < difference; i++ {
			publisherDigit := rand.Intn(10)

			weightedDigitsSum += publisherDigit * weight
			weight--

			arg += strconv.Itoa(publisherDigit)
		}

		arg += "-"

		difference = 8 - difference
		for j := 0; j < difference; j++ {
			titleDigit := rand.Intn(10)

			weightedDigitsSum += titleDigit * weight
			weight--

			arg += strconv.Itoa(titleDigit)
		}

		arg += "-"
		args = append(args, arg)

		if check := (11 - (weightedDigitsSum % 11)) % 11; check == 10 {
			out += arg + "X\n"
		} else {
			out += arg + strconv.Itoa(check) + "\n"
		}
	}

	// Trim the trailing newline.
	out = out[:len(out)-1]

	return
}
