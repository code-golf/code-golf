package hole

import (
	"bufio"
	"bytes"
	"embed"
	"io/fs"
	"math/rand/v2"
	"path/filepath"
	"strings"
)

type test struct{ in, out string }

var fixedTestsMap = map[string][]test{}

//go:embed fixed-tests
var fixedTestsFS embed.FS

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
	return append([]test(nil), fixedTestsMap[holeID]...)
}

func outputTests(tests ...[]test) []Run { return outputTestsWithSep("\n", tests...) }

func outputTestsWithSep(sep string, testRuns ...[]test) []Run {
	runs := make([]Run, len(testRuns))

	for i, tests := range testRuns {
		args := make([]string, len(tests))
		var answer strings.Builder

		for i, t := range tests {
			args[i] = t.in

			if i > 0 {
				answer.WriteString(sep)
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
