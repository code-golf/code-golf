package hole

import (
	"math/rand"
	"strconv"
	"strings"
)

// a fraction (fract) is defined by its
// numerator and denominator in simplified
// form (n, d) and its scale factor (s)
type fract struct{ n, d, s int }

// most of this code is lifted from the
// intersection hole
func strconvunsimplifiedfrac(frac fract) (out string) {
	var outs []string
	outs = append(outs, strconv.Itoa(frac.n*frac.s))
	outs = append(outs, strconv.Itoa(frac.d*frac.s))
	return strings.Join(outs, "/")
}

func strconvsimplifiedfrac(frac fract) (out string) {
	var outs []string
	outs = append(outs, strconv.Itoa(frac.n))
	outs = append(outs, strconv.Itoa(frac.d))
	return strings.Join(outs, "/")
}

// if numerator and denominator are
// divisible by a number greater than
// 1 then the fraction is reducible
func isIrreducible(frac fract) bool {
	gcd := 2
	for gcd <= frac.d {
		if frac.n%gcd == 0 && frac.d%gcd == 0 {
			return false
		}
		gcd++
	}
	return true
}

// Generates a fraction with numerator and denominator both less than 16
func smallFracGen() fract {
	numer := rand.Intn(16)
	denom := rand.Intn(15) + 1

	// Choose a random scale factor so that the
	// largest value in the fraction does not
	// exceed 250
	largerVal := numer
	if denom > largerVal {
		largerVal = denom
	}
	scale := rand.Intn(250/largerVal-1) + 1

	return fract{
		n: numer,
		d: denom,
		s: scale,
	}
}

// Generates a fraction with at least one value greater than 15
func largeFracGen() fract {
	numer := rand.Intn(251)
	denom := rand.Intn(250) + 1

	// If neither value is greater than 15,
	// reassign one chosen at random to a
	// larger number
	if numer <= 15 && denom <= 15 {
		newVal := rand.Intn(235) + 16
		chnged := rand.Intn(2)
		if chnged == 0 {
			numer = newVal
		} else {
			denom = newVal
		}
	}

	// Choose a random scale factor so that the
	// largest value in the fraction does not
	// exceed 250
	largerVal := numer
	if denom > largerVal {
		largerVal = denom
	}
	scale := 1
	if largerVal <= 125 {
		scale = rand.Intn(250/largerVal-1) + 1
	}

	return fract{
		n: numer,
		d: denom,
		s: scale,
	}
}

func fractions() (args []string, out string) {
	var outs []string

	//// default cases
	// hardcoded values
	ra1 := rand.Intn(248) + 2
	ra2 := rand.Intn(248) + 2

	hardcodes := []fract{
		fract{n: 1, d: 1, s: 1},
		fract{n: 1, d: 1, s: ra1},
		fract{n: 1, d: 2, s: 1},
		fract{n: 1, d: 2, s: 10},
		fract{n: 1, d: 2, s: 100},
		fract{n: 2, d: 1, s: 1},
		fract{n: 2, d: 1, s: 100},
		fract{n: 0, d: 1, s: 1},
		fract{n: 0, d: 1, s: ra2},
		fract{n: 15, d: 14, s: 15},
		fract{n: 250, d: 249, s: 1},
		fract{n: 127, d: 126, s: 2},
		fract{n: 126, d: 127, s: 2},
		fract{n: 51, d: 1, s: 5},
		fract{n: 1, d: 51, s: 5},
		fract{n: 249, d: 2, s: 1},
		fract{n: 2, d: 249, s: 1},
		fract{n: 15, d: 13, s: 9},
		fract{n: 15, d: 13, s: 5},
		fract{n: 13, d: 15, s: 9},
		fract{n: 15, d: 11, s: 13},
		fract{n: 11, d: 15, s: 13},
		fract{n: 17, d: 19, s: 13},
		fract{n: 19, d: 17, s: 13},
	}

	for _, testCase := range hardcodes {
		args = append(args, strconvunsimplifiedfrac(testCase))
		outs = append(outs, strconvsimplifiedfrac(testCase))
	}

	// generate 60 random cases with small reduced forms
	cases := 0
	for cases < 60 {
		f := smallFracGen()
		if isIrreducible(f) {
			args = append(args, strconvunsimplifiedfrac(f))
			outs = append(outs, strconvsimplifiedfrac(f))
			cases++
		}
	}

	// generate 40 random cases with at least 1 large number in their reduced forms
	// with at least 20 not being already simplified
	simpCases := 0
	unsimpCases := 0
	for unsimpCases < 20 || unsimpCases+simpCases < 40 {
		f := largeFracGen()
		if isIrreducible(f) {
			if simpCases < 20 && f.s == 1 {
				args = append(args, strconvunsimplifiedfrac(f))
				outs = append(outs, strconvsimplifiedfrac(f))
				simpCases++
			}
			if f.s > 1 {
				args = append(args, strconvunsimplifiedfrac(f))
				outs = append(outs, strconvsimplifiedfrac(f))
				unsimpCases++
			}
		}
	}

	// shuffle args and outputs in the same way
	rand.Shuffle(len(args), func(i, j int) {
		args[i], args[j] = args[j], args[i]
		outs[i], outs[j] = outs[j], outs[i]
	})
	out = strings.Join(outs, "\n")
	return
}
