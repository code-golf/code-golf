package routes

import (
	"math/rand"
	"strconv"
)

func ordinalNumbers() (args []string, out string) {
	for _, i := range rand.Perm(200) {
		s := strconv.Itoa(i)

		args = append(args, s)
		out += s + ord(i) + "\n"
	}

	// Drop the trailing newline.
	out = out[:len(out)-1]

	return
}
