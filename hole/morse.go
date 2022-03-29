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

func morse(reverse bool) (args []string, out string) {
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

	args = []string{strings.Join(words, " ")}

	for _, char := range args[0] {
		out += morseMap[char] + "   "
	}

	// Knock off the trailing three spaces.
	out = out[:len(out)-3]

	if reverse {
		args[0], out = out, args[0]
	}

	return
}
