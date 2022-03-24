package hole

import (
	"fmt"
	"math/rand"
	"strings"
)

type fraction struct{ numerator, denominator, scale int }

// If numerator and denominator are divisible by a number greater than 1 then
// the fraction is reducible.
func (f fraction) isIrreducible() bool {
	for i := 2; i <= f.denominator; i++ {
		if f.numerator%i == 0 && f.denominator%i == 0 {
			return false
		}
	}
	return true
}

// Generates a fraction with numerator and denominator both less than 16
func smallFracGen() fraction {
	f := fraction{numerator: rand.Intn(16), denominator: rand.Intn(15) + 1}

	// Choose a random scale factor so that the largest value in the fraction
	// does not exceed 250.
	max := max(f.numerator, f.denominator)
	f.scale = rand.Intn(250/max-1) + 1

	return f
}

// Generates a fraction with at least one value greater than 15
func largeFracGen() fraction {
	f := fraction{numerator: rand.Intn(251), denominator: rand.Intn(250) + 1}

	// If neither value is greater than 15, reassign one chosen at random to a
	// larger number.
	if f.numerator <= 15 && f.denominator <= 15 {
		newVal := rand.Intn(235) + 16
		if rand.Intn(2) == 0 {
			f.numerator = newVal
		} else {
			f.denominator = newVal
		}
	}

	// Choose a random scale factor so that the largest value in the fraction
	// does not exceed 250.
	f.scale = 1
	if max := max(f.numerator, f.denominator); max <= 125 {
		f.scale = rand.Intn(250/max-1) + 1
	}

	return f
}

func fractions() ([]string, string) {
	// Default cases.
	ra1 := rand.Intn(249) + 2
	ra2 := rand.Intn(249) + 2
	ra3 := rand.Intn(249) + 2
	ra4 := rand.Intn(151) + 100

	fractions := []fraction{
		{numerator: 1, denominator: 1, scale: 1},
		{numerator: 1, denominator: 1, scale: ra1},
		{numerator: 1, denominator: ra3, scale: 1},
		{numerator: ra3, denominator: 1, scale: 1},
		{numerator: 1, denominator: ra4, scale: 1},
		{numerator: ra4, denominator: 1, scale: 1},
		{numerator: 1, denominator: 2, scale: 1},
		{numerator: 1, denominator: 2, scale: 10},
		{numerator: 1, denominator: 2, scale: 100},
		{numerator: 2, denominator: 1, scale: 1},
		{numerator: 2, denominator: 1, scale: 100},
		{numerator: 0, denominator: 1, scale: 1},
		{numerator: 0, denominator: 1, scale: ra2},
		{numerator: 15, denominator: 14, scale: 15},
		{numerator: 250, denominator: 249, scale: 1},
		{numerator: 125, denominator: 124, scale: 2},
		{numerator: 124, denominator: 125, scale: 2},
		{numerator: 50, denominator: 1, scale: 5},
		{numerator: 1, denominator: 50, scale: 5},
		{numerator: 249, denominator: 2, scale: 1},
		{numerator: 2, denominator: 249, scale: 1},
		{numerator: 15, denominator: 13, scale: 9},
		{numerator: 15, denominator: 13, scale: 5},
		{numerator: 13, denominator: 15, scale: 9},
		{numerator: 15, denominator: 11, scale: 13},
		{numerator: 11, denominator: 15, scale: 13},
		{numerator: 17, denominator: 19, scale: 13},
		{numerator: 19, denominator: 17, scale: 13},
	}

	// Generate 60 random cases with small reduced forms.
	for i := 0; i < 60; {
		if f := smallFracGen(); f.isIrreducible() {
			fractions = append(fractions, f)
			i++
		}
	}

	// Generate 40 random cases with at least 1 large number in their reduced
	// forms with at least 20 not being already simplified.
	var simpCases, unsimpCases int
	for unsimpCases < 20 || unsimpCases+simpCases < 40 {
		if f := largeFracGen(); f.isIrreducible() {
			if simpCases < 20 && f.scale == 1 {
				fractions = append(fractions, f)
				simpCases++
			} else if f.scale > 1 {
				fractions = append(fractions, f)
				unsimpCases++
			}
		}
	}

	args := make([]string, len(fractions))
	outs := make([]string, len(fractions))

	for i, f := range shuffle(fractions) {
		args[i] = fmt.Sprintf("%d/%d", f.numerator*f.scale, f.denominator*f.scale)
		outs[i] = fmt.Sprintf("%d/%d", f.numerator, f.denominator)
	}

	return args, strings.Join(outs, "\n")
}
