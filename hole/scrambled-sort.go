package hole

import (
	"cmp"
	"slices"
	"strings"
)

var _ = answerFunc("scrambled-sort", func() []Answer {
	tests := make([]test, 100)

	fixedTests := []struct {
		alphabet []byte
		words    []string
	}{
		// One word, "eerie", sorts before the alphabet itself
		{
			[]byte("eyxmrgjbiunlwczqpsdtkhavfo"),
			shuffle([]string{"eerie", randWord(), randWord(), randWord(), randWord()}),
		},
		// One word is a prefix of another
		{
			shuffle([]byte("abcdefghijklmnopqrstuvwxyz")),
			shuffle([]string{"the", "then", "there", "they", randWord(), randWord()}),
		},
	}

	for i, t := range fixedTests {
		tests[i] = scrambledSortTest(t.alphabet, t.words)
	}

	for i := len(fixedTests); i < len(tests); i++ {
		// Generate a random alphabet and words.
		alphabet := shuffle([]byte("abcdefghijklmnopqrstuvwxyz"))

		words := make([]string, randInt(3, 9))
		for i := range words {
			words[i] = randWord()
		}

		tests[i] = scrambledSortTest(alphabet, words)
	}

	return outputTests(shuffle(tests))
})

func scrambledSortTest(alphabet []byte, words []string) test {
	var t test
	t.in = string(alphabet) + " " + strings.Join(words, " ")

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

	t.out = strings.Join(words, " ")
	return t
}
