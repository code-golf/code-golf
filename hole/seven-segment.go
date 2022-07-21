package hole

import "strings"

var segments = [][]string{
	{" _ ", "   ", " _ ", " _ ", "   ", " _ ", " _ ", " _ ", " _ ", " _ "},
	{"| |", "  |", " _|", " _|", "|_|", "|_ ", "|_ ", "  |", "|_|", "|_|"},
	{"|_|", "  |", "|_ ", " _|", "  |", " _|", "|_|", "  |", "|_|", " _|"},
}

func sevenSegment() []Scorecard {
	digits := shuffle([]byte("00112233445566778899"))
	score := Scorecard{Args: []string{string(digits)}}

	for row := 0; row < 3; row++ {
		for _, digit := range digits {
			score.Answer += segments[row][digit-'0']
		}

		score.Answer = strings.TrimRight(score.Answer, " ")
		if row < 2 {
			score.Answer += "\n"
		}
	}

	return []Scorecard{score}
}
