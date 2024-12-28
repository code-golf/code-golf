package hole

import (
	"math/rand/v2"
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

	argument.WriteString(scrambleAlphabet())

	for range randInt(3, 9) {
		argument.WriteByte(' ')
		argument.WriteString(randWord())
	}

	return argument.String()
}

func scrambleAlphabet() string {
	alphabet := strings.Split("abcdefghijklmnopqrstuvwxyz", "")

	rand.Shuffle(len(alphabet), func(i, j int) {
		alphabet[i], alphabet[j] = alphabet[j], alphabet[i]
	})

	return strings.Join(alphabet, "")
}

func scrambledAlphabetization() []Run {
	tests := make([]test, 0, 100)

	for range 100 {
		argument := generateArgument()

		tests = append(tests, test{
			argument,
			alphabetizeWords(argument),
		})
	}

	return outputTests(shuffle(tests))
}
