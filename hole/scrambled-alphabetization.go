package hole

import (
	"sort"
	"strings"
)

func alphabetizeWords(alphabet string, words []string) string {
	order := map[rune]int{}
	for i, ch := range alphabet {
		order[ch] = i
	}

	sort.Slice(words, func(i, j int) bool {
		for k := 0; k < len(words[i]) && k < len(words[j]); k++ {
			if order[rune(words[i][k])] < order[rune(words[j][k])] {
				return true
			} else if order[rune(words[i][k])] > order[rune(words[j][k])] {
				return false
			}
		}

		return len(words[i]) < len(words[j])
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

	return outputTests(shuffle(tests))
})
