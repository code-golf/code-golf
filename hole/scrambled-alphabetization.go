package hole

import (
	"sort"
	"strings"
)

func alphabetizeWords(s string) string {
	fields := strings.Fields(s)
	order, words := make(map[rune]int), fields[1:]

	for i, ch := range fields[0] {
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

func generateArgument() string {
	var argument strings.Builder

	alphabet := string(shuffle([]byte("abcdefghijklmnopqrstuvwxyz")))
	argument.WriteString(alphabet)

	for range randInt(3, 9) {
		argument.WriteByte(' ')
		argument.WriteString(randWord())
	}

	return argument.String()
}

var _ = answerFunc("scrambled-alphabetization", func() []Answer {
	tests := make([]test, 100)

	for i := range tests {
		argument := generateArgument()

		tests[i] = test{argument, alphabetizeWords(argument)}
	}

	return outputTests(shuffle(tests))
})
