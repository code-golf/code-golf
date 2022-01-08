package hole

import (
	"math/rand"
	"strconv"
	"time"
)

func isbn() (args []string, out string) {

	for m:=0; m < 20; m++ {


	// Initialize the ISBN string
	var ISBNarg = ""
	var weightedDigitsSum = 0
	var weight = 10

	//random number seeding
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	//first digit of ISBN, not sticking with traditional 1 or 0, can't let them exploit that.
	var firstDigit = r1.Intn(10)

	weightedDigitsSum += firstDigit * weight

	weight--

	ISBNarg = strconv.Itoa(firstDigit) + "-"

	// This here logic is for varying the second two parts of the ISBN. Sure, it's cosmetic, but it might mess some people up.
	difference := 6 - r1.Intn(5)
	for i := 0; i < difference; i++ {
		var publisherDigit = r1.Intn(10)

		weightedDigitsSum += publisherDigit * weight

		weight--

		ISBNarg += strconv.Itoa(publisherDigit)
	}

	ISBNarg += "-"

	secondDifference := 8 - difference
	for j := 0; j < secondDifference; j++ {
		var titleDigit = r1.Intn(10)

		weightedDigitsSum += titleDigit * weight

		weight--

		ISBNarg += strconv.Itoa(titleDigit)
	}

	ISBNarg += "-"

	args = append(args, ISBNarg)

	var checkDigit = (11 - (weightedDigitsSum % 11)) % 11

	if checkDigit == 10 {
		out += ISBNarg + "X" + "\n"
	} else {
		out+= ISBNarg + strconv.Itoa(checkDigit) + "\n"
	}


	}//end of m-loop

	//shamelessly stolen
	out = out[:len(out)-1]


	return

} //end of isbn func

