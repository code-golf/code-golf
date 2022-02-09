package hole

import (
	"math/rand"
	"strconv"

	"github.com/code-golf/code-golf/pretty"
)

func ordinalNumbers() (args []string, out string) {
	for _, i := range rand.Perm(200) {
		if i > 100 {
			i += rand.Intn(20) * 100
		}
		s := strconv.Itoa(i)

		args = append(args, s)
		out += s + pretty.Ordinal(i) + "\n"
	}

	// Drop the trailing newline.
	out = out[:len(out)-1]

	return
}
