package routes

import (
	"math/rand"
	"strings"
	"unicode"
)

var segments = [][]string{
	{" _ ", "   ", " _ ", " _ ", "   ", " _ ", " _ ", " _ ", " _ ", " _ "},
	{"| |", "  |", " _|", " _|", "|_|", "|_ ", "|_ ", "  |", "|_|", "|_|"},
	{"|_|", "  |", "|_ ", " _|", "  |", " _|", "|_|", "  |", "|_|", " _|"},
}

func sevenSegment() (args []string, out string) {
	digits := []byte{'0', '0', '1', '1', '2', '2', '3', '3', '4', '4', '5', '5', '6', '6', '7', '7', '8', '8', '9', '9'}

	// Shuffle the digits
	for i := range digits {
		j := rand.Intn(i + 1)
		digits[i], digits[j] = digits[j], digits[i]
	}

	args = append(args, string(digits))

	for row := 0; row < 3; row++ {
		for _, digit := range digits {
			out += segments[row][digit-'0']
		}

		out = strings.TrimRightFunc(out, unicode.IsSpace) + "\n"
	}

	out = strings.TrimRightFunc(out, unicode.IsSpace)

	return
}
