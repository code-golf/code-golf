package hole

import (
	"bufio"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/agnivade/levenshtein"
)

var words []string

func init() {
	file, err := os.Open("words.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}
}

func levenshteinDistance() ([]string, string) {
	const count = 20

	args := make([]string, count)
	outs := make([]string, count)
	perm := rand.Perm(count)

	for i := 0; i < count; i++ {
		a := words[rand.Intn(len(words))]

		if i == perm[0] {
			// Ensure we have at least one zero distance.
			args[i] = a + " " + a
			outs[i] = "0"
			continue
		} else if i == perm[1] {
			// Add a test case that blocks an incorrect simplification to the algorithm from working.
			args[i] = "open however"
			outs[i] = "5"
			continue
		}

		b := words[rand.Intn(len(words))]

		args[i] = a + " " + b
		outs[i] = strconv.Itoa(levenshtein.ComputeDistance(a, b))
	}

	return args, strings.Join(outs, "\n")
}
