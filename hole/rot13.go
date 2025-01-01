package hole

import "strings"

func rot13() []Run {
	tests := make([]test, 0, 100)

	for range 100 {
		argument := generateSequence()

		tests = append(tests, test{
			argument,
			encodeDecodeSequence(argument),
		})
	}

	return outputTests(shuffle(tests))
}

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
