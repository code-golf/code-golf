package hole

import (
	"cmp"
	"slices"
	"strings"
)

func alphabetizeWords(alphabet string, words []string) string {
	order := map[byte]int{}
	for i := range alphabet {
		order[alphabet[i]] = i
	}

	slices.SortFunc(words, func(a, b string) int {
		for i := range min(len(a), len(b)) {
			if c := cmp.Compare(order[a[i]], order[b[i]]); c != 0 {
				return c
			}
		}

		return cmp.Compare(len(a), len(b))
	})

	return strings.Join(words, " ")
}

var _ = answerFunc("scrambled-alphabetization", func() []Answer {
	tests := make([]test, 100)

	for i := range tests {
		alphabet := string(shuffle([]byte("abcdefghijklmnopqrstuvwxyz")))

		words := make([]string, randInt(3, 9))
		for i := range words {
			words[i] = randWord()
		}

		tests[i] = test{
			in:  alphabet + " " + strings.Join(words, " "),
			out: alphabetizeWords(alphabet, words),
		}
	}

	return outputTests(tests)
})
