package hole

import (
	"cmp"
	"slices"
	"strings"
)

var _ = answerFunc("scrambled-sort", func() []Answer {
	tests := make([]test, 100)

	for i := range tests {
		// Generate a random alphabet and words.
		alphabet := shuffle([]byte("abcdefghijklmnopqrstuvwxyz"))

		words := make([]string, randInt(3, 9))
		for i := range words {
			words[i] = randWord()
		}

		tests[i].in = string(alphabet) + " " + strings.Join(words, " ")

		// Sort words according to the random alphabet.
		order := map[byte]int{}
		for i, b := range alphabet {
			order[b] = i
		}

		slices.SortFunc(words, func(a, b string) int {
			for i := range min(len(a), len(b)) {
				if c := cmp.Compare(order[a[i]], order[b[i]]); c != 0 {
					return c
				}
			}

			return cmp.Compare(len(a), len(b))
		})

		tests[i].out = strings.Join(words, " ")
	}

	return outputTests(tests)
})
