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
	same := rand.Intn(count)

	for i := 0; i < count; i++ {
		a := words[rand.Intn(len(words))]

		// Ensure we have at least one zero distance.
		if i == same {
			args[i] = a + " " + a
			outs[i] = "0"
			continue
		}

		b := words[rand.Intn(len(words))]

		args[i] = a + " " + b
		outs[i] = strconv.Itoa(levenshtein.ComputeDistance(a, b))
	}

	return args, strings.Join(outs, "\n")
}
