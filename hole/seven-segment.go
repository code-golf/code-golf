package hole

import "strings"

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

		out = strings.TrimRight(out, " ")
		if row < 2 {
			out += "\n"
		}
	}

	return
}
