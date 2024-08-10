package hole

import (
	_ "embed"
	"strconv"

	"github.com/agnivade/levenshtein"
)

func levenshteinTest(a, b string) test {
	return test{a + " " + b, strconv.Itoa(levenshtein.ComputeDistance(a, b))}
}

func levenshteinDistance() []Run {
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
