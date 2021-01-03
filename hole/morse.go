package hole

import (
	"math/rand"
	"strings"
)

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

func shuffle(str string) string {
	chars := []byte(str)

	rand.Shuffle(len(chars), func(i, j int) {
		chars[i], chars[j] = chars[j], chars[i]
	})

	return string(chars)
}

func morse(reverse bool) (args []string, out string) {
	words := []string{
		"BUD",
		"FOR",
		"JIGS",
		"NYMPH",
		"QUICK",
		"VEX",
		"WALTZ",
		shuffle("0123456789"),
		shuffle("ABCDEFGHIJKLMNOPQRSTUVWXYZ"),
	}

	rand.Shuffle(len(words), func(i, j int) {
		words[i], words[j] = words[j], words[i]
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
