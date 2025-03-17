package hole

import (
	_ "embed"
	"strconv"

	"github.com/agnivade/levenshtein"
)

func levenshteinTest(a, b string) test {
	return test{a + " " + b, strconv.Itoa(levenshtein.ComputeDistance(a, b))}
}

var _ = answerFunc("levenshtein-distance", func() []Answer {
	word := randWord()
	tests := []test{
		levenshteinTest(word, word),
		levenshteinTest("a", randWord()),
		levenshteinTest(randWord(), "a"),
		levenshteinTest("incomprehensible", randWord()),
		levenshteinTest("open", "however"),
		levenshteinTest("however", "open"),
		levenshteinTest("large", "hypothetical"),
		levenshteinTest("hypothetical", "set"),
		levenshteinTest("very", "incomprehensible"),
		levenshteinTest("apprentice", "point"),
		levenshteinTest("school", "school"),
	}

	for i := len(tests); i < 40; i++ {
		tests = append(tests, levenshteinTest(randWord(), randWord()))
	}

	tests = shuffle(tests)

	const argc = 20 // Preserve original argc
	return outputTests(tests[:argc], tests[len(tests)-argc:])
})
