package hole

import (
	_ "embed"
	"math/rand"
	"strconv"
	"strings"

	"github.com/agnivade/levenshtein"
)

var (
	//go:embed words.txt
	wordsTxt string
	words    = strings.Fields(wordsTxt)
)

func randWord() string { return words[rand.Intn(len(words))] }

func levenshteinDistance() ([]string, string) {
	const count = 20

	a := randWord()
	b := randWord()
	c := randWord()
	tests := []test{
		{a + " " + a, "0"},
		{"a " + b, levenshtein.ComputeDistance("a", b)},
		{"incomprehensible " + c, levenshtein.ComputeDistance("incomprehensible ", c},
		{"open however", "5"},
		{"however open", "5"},
		{"large hypothetical", "11"},
	}

	for i := len(tests); i < count; i++ {
		a := randWord()
		b := randWord()
		tests = append(tests, test{
			a + " " + b,
			strconv.Itoa(levenshtein.ComputeDistance(a, b)),
		})
	}

	return outputTests(shuffle(tests))
}
