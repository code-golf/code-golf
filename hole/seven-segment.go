package hole

import (
	"strings"
	"unicode"
)

var segments = [][]string{
	{" _ ", "   ", " _ ", " _ ", "   ", " _ ", " _ ", " _ ", " _ ", " _ "},
	{"| |", "  |", " _|", " _|", "|_|", "|_ ", "|_ ", "  |", "|_|", "|_|"},
	{"|_|", "  |", "|_ ", " _|", "  |", " _|", "|_|", "  |", "|_|", " _|"},
}

func sevenSegment() (args []string, out string) {
	digits := shuffle([]byte("00112233445566778899"))

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
