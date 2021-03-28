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

func randWord() string {
	return words[rand.Intn(len(words))]
}

func levenshteinDistance() ([]string, string) {
	const count = 20

	args := make([]string, count)
	outs := make([]string, count)
	perm := rand.Perm(count)

	for i := 0; i < count; i++ {
		switch i {
		case perm[0]:
			// Ensure we have at least one zero distance.
			a := randWord()
			args[i] = a + " " + a
			outs[i] = "0"
		case perm[1]:
			// Add a test case that blocks an incorrect simplification to the algorithm from working.
			args[i] = "open however"
			outs[i] = "5"
		case perm[2]:
			args[i] = "however open"
			outs[i] = "5"
		// Ensure we have a double digit distance, TODO randomise the words?
		case perm[3]:
			args[i] = "large hypothetical"
			outs[i] = "11"
		default:
			a := randWord()
			b := randWord()
			args[i] = a + " " + b
			outs[i] = strconv.Itoa(levenshtein.ComputeDistance(a, b))
		}
	}

	return args, strings.Join(outs, "\n")
}
