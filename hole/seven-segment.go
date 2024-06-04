package hole

import "strings"

var segments = [][]string{
	{" _ ", "   ", " _ ", " _ ", "   ", " _ ", " _ ", " _ ", " _ ", " _ "},
	{"| |", "  |", " _|", " _|", "|_|", "|_ ", "|_ ", "  |", "|_|", "|_|"},
	{"|_|", "  |", "|_ ", " _|", "  |", " _|", "|_|", "  |", "|_|", " _|"},
}

func sevenSegment() []Run {
	digits := shuffle([]byte("00112233445566778899"))
	run := Run{Args: []string{string(digits)}}

	for row, segment := range segments {
		for _, digit := range digits {
			run.Answer += segment[digit-'0']
		}

		run.Answer = strings.TrimRight(run.Answer, " ")
		if row < 2 {
			run.Answer += "\n"
		}
	}

	return []Run{run}
}
