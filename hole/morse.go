package hole

import "strings"

var morseMap = map[rune]string{
	'A': "▄ ▄▄▄",
	'B': "▄▄▄ ▄ ▄ ▄",
	'C': "▄▄▄ ▄ ▄▄▄ ▄",
	'D': "▄▄▄ ▄ ▄",
	'E': "▄",
	'F': "▄ ▄ ▄▄▄ ▄",
	'G': "▄▄▄ ▄▄▄ ▄",
	'H': "▄ ▄ ▄ ▄",
	'I': "▄ ▄",
	'J': "▄ ▄▄▄ ▄▄▄ ▄▄▄",
	'K': "▄▄▄ ▄ ▄▄▄",
	'L': "▄ ▄▄▄ ▄ ▄",
	'M': "▄▄▄ ▄▄▄",
	'N': "▄▄▄ ▄",
	'O': "▄▄▄ ▄▄▄ ▄▄▄",
	'P': "▄ ▄▄▄ ▄▄▄ ▄",
	'Q': "▄▄▄ ▄▄▄ ▄ ▄▄▄",
	'R': "▄ ▄▄▄ ▄",
	'S': "▄ ▄ ▄",
	'T': "▄▄▄",
	'U': "▄ ▄ ▄▄▄",
	'V': "▄ ▄ ▄ ▄▄▄",
	'W': "▄ ▄▄▄ ▄▄▄",
	'X': "▄▄▄ ▄ ▄ ▄▄▄",
	'Y': "▄▄▄ ▄ ▄▄▄ ▄▄▄",
	'Z': "▄▄▄ ▄▄▄ ▄ ▄",
	'1': "▄ ▄▄▄ ▄▄▄ ▄▄▄ ▄▄▄",
	'2': "▄ ▄ ▄▄▄ ▄▄▄ ▄▄▄",
	'3': "▄ ▄ ▄ ▄▄▄ ▄▄▄",
	'4': "▄ ▄ ▄ ▄ ▄▄▄",
	'5': "▄ ▄ ▄ ▄ ▄",
	'6': "▄▄▄ ▄ ▄ ▄ ▄",
	'7': "▄▄▄ ▄▄▄ ▄ ▄ ▄",
	'8': "▄▄▄ ▄▄▄ ▄▄▄ ▄ ▄",
	'9': "▄▄▄ ▄▄▄ ▄▄▄ ▄▄▄ ▄",
	'0': "▄▄▄ ▄▄▄ ▄▄▄ ▄▄▄ ▄▄▄",
	' ': "    ",
}

func morse(reverse bool) []Scorecard {
	runs := 1
	if reverse {
		runs = 3
	}

	scorecards := make([]Scorecard, runs)

	for i := 0; i < runs; i++ {
		words := shuffle([]string{
			"BUD",
			"FOR",
			"JIGS",
			"NYMPH",
			"QUICK",
			"VEX",
			"WALTZ",
			string(shuffle([]byte("0123456789"))),
			string(shuffle([]byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ"))),
		})

		arg := strings.Join(words, " ")
		out := ""

		// Randomly remove a few characters for Morse Decoder.
		if reverse {
			arg = string(shuffle([]byte(arg)))[:len(arg)-4]
		}

		for _, char := range arg {
			out += morseMap[char] + "   "
		}

		// Knock off the trailing three spaces.
		out = out[:len(out)-3]

		if reverse {
			arg, out = out, arg
		}
		scorecards[i] = Scorecard{Args: []string{arg}, Answer: out}
	}

	return scorecards
}
