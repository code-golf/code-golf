package hole

import (
	"math/rand/v2"
	"strings"
)

func blackjackDealer(d []string) string {
	v, a := 0, 0

	for i, c := range d {
		switch c[0] {
		case 'A':
			v += 11
			a++
		case 'T', 'J', 'Q', 'K':
			v += 10
		default:
			v += int(c[0] - '0')
		}

		for v > 21 && a > 0 {
			v -= 10
			a--
		}

		if v > 21 || v > 17 || v == 17 && a == 0 {
			return strings.Join(d[:i+1], " ")
		}
	}

	return strings.Join(d, " ")
}

var _ = answerFunc("blackjack-dealer", func() []Answer {
	const ranks = "A23456789TJQK"
	const suits = "CDHS"

	deck := make([]string, 0, 52)
	for _, r := range ranks {
		for _, s := range suits {
			deck = append(deck, string([]byte{byte(r), byte(s)}))
		}
	}

	makeDeck := func(prefix string) []string {
		d := make([]string, 0, 52)
		used := map[string]bool{}

		for _, c := range strings.Fields(prefix) {
			d = append(d, c)
			used[c] = true
		}

		for _, c := range deck {
			if !used[c] {
				d = append(d, c)
			}
		}

		return d
	}

	tests := make([]test, 0, 100)

	prefixes := []string{
		"AH 6C 2S",
		"TH 7C AS",
		"9S 2H 4D AD TH",
		"AS 6C 2D",
		"AS AD 5C KD",
		"AC AD AH AS 2C 3D 4H 5S 6C",
		"2C 2D 2H 2S 9C 9D",
		"5C 5D 5H 5S 6C",
		"KC QD JH TS",
		"AC 2C 3C 4C 5C 6C 7C",
		"7C 4D 6H",
		"8C 8D",
		"9C 7D",
		"AC 5D",
		"AD 5H",
		"AH 2C 2D 2H 2S",
		"AS 9C",
		"AS 8C",
		"AS 7C 2D",
		"AS 6C 2D 3H",
	}

	for _, p := range prefixes {
		d := makeDeck(p)
		tests = append(tests, test{
			strings.Join(d, " "),
			blackjackDealer(d),
		})
	}

	for len(tests) < 100 {
		d := append([]string(nil), deck...)

		rand.Shuffle(len(d), func(i, j int) {
			d[i], d[j] = d[j], d[i]
		})

		tests = append(tests, test{
			strings.Join(d, " "),
			blackjackDealer(d),
		})
	}

	return outputTests(shuffle(tests))
})
