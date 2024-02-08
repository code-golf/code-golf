package hole

import (
	"embed"
	"math/rand/v2"
	"path"
	"strings"
)

type test struct{ in, out string }

var fixedTestsMap = map[string][]test{}

//go:embed fixed-tests
var fixedTestsFS embed.FS

func init() {
	const dir = "fixed-tests"

	files, err := fixedTestsFS.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		name := file.Name()

		txt, err := fixedTestsFS.ReadFile(path.Join(dir, name))
		if err != nil {
			panic(err)
		}

		tokens := strings.Split(strings.Trim(string(txt), "\n"), "\n\n")
		tests := make([]test, 0, len(tokens)/2)

		for i := 0; i < len(tokens); i += 2 {
			tests = append(tests, test{tokens[i], tokens[i+1]})
		}

		fixedTestsMap[strings.TrimSuffix(name, path.Ext(name))] = tests
	}
}

// Return a copy so holes are free to append, shuffle, etc.
func fixedTests(holeID string) []test {
	// return fixedTestsMap[holeID]
	return append([]test(nil), fixedTestsMap[holeID]...)
}

func outputTests(testRuns ...[]test) []Run {
	runs := make([]Run, len(testRuns))

	for i, tests := range testRuns {
		args := make([]string, len(tests))
		var answer strings.Builder

		for i, t := range tests {
			args[i] = t.in

			if i > 0 {
				answer.WriteByte('\n')
			}
			answer.WriteString(t.out)
		}

		runs[i] = Run{Args: args, Answer: answer.String()}
	}

	return runs
}

func outputMultirunTests(tests []test) []Run {
	shuffle(tests)
	mid := len(tests) / 2
	return outputTests(tests, tests[:mid], tests[mid:])
}

// Return a random element from the given slice. Panics on empty slice.
func randChoice[E any](x []E) E { return x[rand.IntN(len(x))] }

// Return a random integer between min and max inclusive.
func randInt(min, max int) int { return min + rand.IntN(max-min+1) }

// Returning the slice is a convenience, the shuffle is still in-place.
func shuffle[E any](x []E) []E {
	rand.Shuffle(len(x), func(i, j int) { x[i], x[j] = x[j], x[i] })
	return x
}
