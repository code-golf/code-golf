package hole

import (
	"math/rand/v2"
	"strings"
)

var _ = answerFunc("rot13", func() []Answer {
	tests := make([]test, 0, 100)

	// Enforce one sequence consisting of only one alphabetic character, randomly generated.
	argument := string(rune(randInt(65, 90) + 32*rand.IntN(2)))

	for i := 0; i < 100; {
		if i > 0 {
			argument = generateSequence()
		}

		// Prevent arguments from starting with '@' because some languages cannot handle them correctly.
		if argument[0] == '@' {
			continue
		}

		i++

		tests = append(tests, test{
			argument,
			encodeDecodeSequence(argument),
		})
	}

	return outputTests(shuffle(tests))
})

func generateSequence() string {
	var sequence strings.Builder

	for range randInt(25, 100) {
		sequence.WriteRune(rune(randInt(32, 126)))
	}

	return sequence.String()
}

func encodeDecodeSequence(s string) string {
	var sequence strings.Builder

	for _, ch := range s {
		switch {
		case ch >= 'A' && ch <= 'Z':
			sequence.WriteRune('A' + (ch-'A'+13)%26)
		case ch >= 'a' && ch <= 'z':
			sequence.WriteRune('a' + (ch-'a'+13)%26)
		default:
			sequence.WriteRune(ch)
		}
	}

	return sequence.String()
}
