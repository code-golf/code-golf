package hole

import (
	"fmt"
	"math/rand/v2"
	"strings"
)

var deck = [...]string{
	"1GED", "1GEO", "1GEW", "1GHD", "1GHO", "1GHW", "1GSD", "1GSO", "1GSW",
	"1PED", "1PEO", "1PEW", "1PHD", "1PHO", "1PHW", "1PSD", "1PSO", "1PSW",
	"1RED", "1REO", "1REW", "1RHD", "1RHO", "1RHW", "1RSD", "1RSO", "1RSW",
	"2GED", "2GEO", "2GEW", "2GHD", "2GHO", "2GHW", "2GSD", "2GSO", "2GSW",
	"2PED", "2PEO", "2PEW", "2PHD", "2PHO", "2PHW", "2PSD", "2PSO", "2PSW",
	"2RED", "2REO", "2REW", "2RHD", "2RHO", "2RHW", "2RSD", "2RSO", "2RSW",
	"3GED", "3GEO", "3GEW", "3GHD", "3GHO", "3GHW", "3GSD", "3GSO", "3GSW",
	"3PED", "3PEO", "3PEW", "3PHD", "3PHO", "3PHW", "3PSD", "3PSO", "3PSW",
	"3RED", "3REO", "3REW", "3RHD", "3RHO", "3RHW", "3RSD", "3RSO", "3RSW",
}

func isSet(a, b, c string) bool {
	for i := range 4 {
		if (a[i] == b[i] && b[i] == c[i]) || (a[i] != b[i] && b[i] != c[i] && a[i] != c[i]) {
			continue
		}

		return false
	}

	return true
}

var _ = answerFunc("set", func() []Answer {
	tests := make([]test, 100)

outer:
	for i := 0; i < 100; {

		// Pick 12 random cards from deck to form a hand.
		cards := make([]string, 12)
		for j, k := range rand.Perm(len(deck))[:len(cards)] {
			cards[j] = deck[k]
		}

		// Use this hand if it contains exactly one set.
		var expected string
		for j := range cards {
			for k := j + 1; k < len(cards); k++ {
				for l := k + 1; l < len(cards); l++ {
					if isSet(cards[j], cards[k], cards[l]) {
						// Bail early if we've already found a set.
						if expected != "" {
							continue outer
						}

						expected = fmt.Sprintf("%s %s %s", cards[j], cards[k], cards[l])
					}
				}
			}
		}

		// Only use this hand if we found a set.
		if expected != "" {
			tests[i] = test{strings.Join(cards, " "), expected}
			i++
		}
	}

	return outputTests(tests)
})
