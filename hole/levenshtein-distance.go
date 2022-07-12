package hole

import (
	_ "embed"
	"strconv"
	"strings"

	"github.com/agnivade/levenshtein"
)

var (
	//go:embed words.txt
	wordsTxt string
	words    = strings.Fields(wordsTxt)
)

func randWord() string { return randChoice(words) }

func levenshteinTest(a, b string) test {
	return test{a + " " + b, strconv.Itoa(levenshtein.ComputeDistance(a, b))}
}

func levenshteinDistance() []Scorecard {
	word := randWord()
	tests := []test{
		levenshteinTest(word, word),
		levenshteinTest("a", randWord()),
		levenshteinTest("incomprehensible", randWord()),
		levenshteinTest("open", "however"),
		levenshteinTest("however", "open"),
		levenshteinTest("large", "hypothetical"),
	}

	for i := len(tests); i < 20; i++ {
		tests = append(tests, levenshteinTest(randWord(), randWord()))
	}

	return outputTests(shuffle(tests))
}
