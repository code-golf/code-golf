package hole

import "strings"

var segments = [][]string{
	{" _ ", "   ", " _ ", " _ ", "   ", " _ ", " _ ", " _ ", " _ ", " _ "},
	{"| |", "  |", " _|", " _|", "|_|", "|_ ", "|_ ", "  |", "|_|", "|_|"},
	{"|_|", "  |", "|_ ", " _|", "  |", " _|", "|_|", "  |", "|_|", " _|"},
}

func sevenSegment() []Run {
	digits := shuffle([]byte("00112233445566778899"))
	t := test{in: string(digits)}

	for row, segment := range segments {
		for _, digit := range digits {
			t.out += segment[digit-'0']
		}

		t.out = strings.TrimRight(t.out, " ")
		if row < 2 {
			t.out += "\n"
		}
	}

	return outputTests([]test{t})
}
