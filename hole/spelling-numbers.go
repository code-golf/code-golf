package hole

import (
	"math/rand"
	"strconv"
	"strings"
)

var teens = [][]byte{
	[]byte("zero"),
	[]byte("one"),
	[]byte("two"),
	[]byte("three"),
	[]byte("four"),
	[]byte("five"),
	[]byte("six"),
	[]byte("seven"),
	[]byte("eight"),
	[]byte("nine"),
	[]byte("ten"),
	[]byte("eleven"),
	[]byte("twelve"),
	[]byte("thirteen"),
	[]byte("fourteen"),
	[]byte("fifteen"),
	[]byte("sixteen"),
	[]byte("seventeen"),
	[]byte("eighteen"),
	[]byte("nineteen"),
}

var tens = [][]byte{
	nil,
	nil,
	[]byte("twenty"),
	[]byte("thirty"),
	[]byte("forty"),
	[]byte("fifty"),
	[]byte("sixty"),
	[]byte("seventy"),
	[]byte("eighty"),
	[]byte("ninety"),
}

func wordify(out *strings.Builder, i int) {
	if i == 1000 {
		out.WriteString("one thousand")
	} else if i < 20 {
		out.Write(teens[i])
	} else if i < 100 {
		out.Write(tens[i/10])

		if j := i % 10; j > 0 {
			out.WriteRune('-')
			out.Write(teens[j])
		}
	} else {
		out.Write(teens[i/100])
		out.WriteString(" hundred")

		if j := i % 100; j > 0 {
			out.WriteString(" and ")
			wordify(out, j)
		}
	}
}

func spellingNumbers() ([]string, string) {
	const max = 1000

	args := make([]string, max+1)
	var out strings.Builder

	// The strings "zero" to "one thousand", newline delimited, are this len.
	out.Grow(25531)

	for i, n := range rand.Perm(max + 1) {
		args[i] = strconv.Itoa(n)

		wordify(&out, n)

		if i < max {
			out.WriteRune('\n')
		}
	}

	return args, out.String()
}
