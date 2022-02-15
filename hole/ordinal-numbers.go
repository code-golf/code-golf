package hole

import (
	"math/rand"
	"strconv"
	"strings"

	"github.com/code-golf/code-golf/pretty"
)

func ordinalNumbers() ([]string, string) {
	const count = 1000

	args := make([]string, count)
	var out strings.Builder

	// The strings "0th" to "999th", newline delimited, are this len.
	out.Grow(5889)

	for i, n := range rand.Perm(count) {
		args[i] = strconv.Itoa(n)

		if i > 0 {
			out.WriteByte('\n')
		}
		out.WriteString(args[i])
		out.WriteString(pretty.Ordinal(n))
	}

	return args, out.String()
}
