package hole

import (
	"math/rand/v2"
	"strconv"
	"strings"
)

var _ = answerFunc("blackjack-dealer", func() []Answer {
	ranks := "A23456789TJQK"
	suits := "CDHS"
	deck := make([]string, 0, 52)
	for _, r := range ranks {
		for _, s := range suits {
			deck = append(deck, string([]byte{byte(r), byte(s)}))
		}
	}

	dealer := func(d []string) string {
		h := append([]string(nil), d[:2]...)
		value := func() (int, bool) {
			v, a := 0, 0
			for _, c := range h {
				switch c[0] {
				case 'A':
					v += 11
					a++
				case 'T', 'J', 'Q', 'K':
					v += 10
				default:
					v += int(c[0] - '0')
				}
			}
			for v > 21 && a > 0 {
				v -= 10
				a--
			}
			return v, a > 0
		}

		for i := 2; ; i++ {
			v, soft := value()
			if v > 21 || v > 17 || v == 17 && !soft {
				return strings.Join(h, " ") + " " + func() string {
					if v > 21 {
						return "BUST"
					}
					return strconv.Itoa(v)
				}()
			}
			h = append(h, d[i])
		}
	}

	makeDeck := func(prefix ...string) []string {
		d := append([]string(nil), prefix...)
		used := map[string]bool{}
		for _, c := range prefix {
			used[c] = true
		}
		for _, c := range deck {
			if !used[c] {
				d = append(d, c)
			}
		}
		return d
	}

	prefixes := [][]string{
		{"AH", "6C", "5S", "4D"},
		{"AH", "6C", "2S", "9D"},
		{"TH", "7C", "AS"},
		{"8C", "4C", "QD", "JD", "TH"},
		{"AS", "AS", "9C", "KH"},
		{"AC", "AD", "AH", "AS", "2C", "3D"},
	}

	tests := make([]test, 0, 106)
	for _, p := range prefixes {
		d := makeDeck(p...)
		tests = append(tests, test{strings.Join(d, " "), dealer(d)})
	}

	for len(tests) < 106 {
		d := append([]string(nil), deck...)
		rand.Shuffle(len(d), func(i, j int) {
			d[i], d[j] = d[j], d[i]
		})
		tests = append(tests, test{strings.Join(d, " "), dealer(d)})
	}

	shuffle(tests)
	return outputTests(tests[:53], tests[53:])
})
