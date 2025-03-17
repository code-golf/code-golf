package hole

import (
	"bufio"
	"bytes"
	"embed"
	"io/fs"
	"math/rand/v2"
	"path/filepath"
	"slices"
	"strings"

	"github.com/code-golf/code-golf/config"
)

// Alias HoleAnswer into this package to reduce boilerplate in every callsite.
type Answer = config.HoleAnswer

type test struct{ in, out string }

var fixedTestsMap = map[string][]test{}

//go:embed fixed-tests
var fixedTestsFS embed.FS

//go:embed words.txt
var wordsTxt string
var words = strings.Fields(wordsTxt)

func init() {
	if err := fs.WalkDir(fixedTestsFS, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() || err != nil {
			return err
		}

		txt, err := fixedTestsFS.Open(path)
		if err != nil {
			return err
		}
		scanner := bufio.NewScanner(txt)

		var tests []test
		var in, out strings.Builder
		for scanner.Scan() {
			line := scanner.Bytes()

			if b, ok := bytes.CutPrefix(line, []byte{'<'}); ok {
				if in.Len() > 0 {
					in.WriteByte('\n')
				}
				in.Write(bytes.TrimPrefix(b, []byte{' '}))
			} else if b, ok := bytes.CutPrefix(line, []byte{'>'}); ok {
				if out.Len() > 0 {
					out.WriteByte('\n')
				}
				out.Write(bytes.TrimPrefix(b, []byte{' '}))
			} else {
				tests = append(tests, test{in.String(), out.String()})
				in.Reset()
				out.Reset()
			}
		}

		tests = append(tests, test{in.String(), out.String()})

		name := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
		fixedTestsMap[name] = tests

		return scanner.Err()
	}); err != nil {
		panic(err)
	}
}

// Return a copy so holes are free to append, shuffle, etc.
func fixedTests(holeID string) []test {
	return slices.Clone(fixedTestsMap[holeID])
}

// Set the answer func on the Hole by ID, return any to allow top-level call.
func answerFunc(holeID string, f config.HoleAnswerFunc) any {
	config.AllHoleByID[holeID].AnswerFunc = f
	return nil
}

func outputTests(tests ...[]test) []Answer { return outputTestsWithSep("\n", tests...) }

func outputTestsWithSep(sep string, testRuns ...[]test) []Answer {
	answers := make([]Answer, len(testRuns))

	for i, tests := range testRuns {
		args := make([]string, len(tests))
		var answer strings.Builder

		for i, t := range tests {
			args[i] = t.in

			if t.out != "" {
				answer.WriteString(t.out)
				answer.WriteString(sep)
			}
		}

		answers[i] = Answer{Args: args, Answer: answer.String()}
	}

	return answers
}

func outputMultirunTests(tests []test) []Answer {
	shuffle(tests)
	mid := len(tests) / 2
	return outputTests(tests, tests[:mid], tests[mid:])
}

// Return a random element from the given slice. Panics on empty slice.
func randChoice[E any](x []E) E { return x[rand.IntN(len(x))] }

// Return a random integer between min and max inclusive.
func randInt(min, max int) int { return min + rand.IntN(max-min+1) }

// Return a random word from words.txt.
func randWord() string { return randChoice(words) }

// Returning the slice is a convenience, the shuffle is still in-place.
func shuffle[E any](x []E) []E {
	rand.Shuffle(len(x), func(i, j int) { x[i], x[j] = x[j], x[i] })
	return x
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		a = -a
	}
	return a
}
