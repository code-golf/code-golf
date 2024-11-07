package hole

import (
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

func wordify(i int) string {
	var out strings.Builder
	if i == 1000 {
		out.WriteString("one thousand")
	} else if i < 20 {
		out.Write(teens[i])
	} else if i < 100 {
		out.Write(tens[i/10])

		if j := i % 10; j > 0 {
			out.WriteByte('-')
			out.Write(teens[j])
		}
	} else {
		out.Write(teens[i/100])
		out.WriteString(" hundred")

		if j := i % 100; j > 0 {
			out.WriteString(" and ")
			out.WriteString(wordify(j))
		}
	}
	return out.String()
}

func spellingNumbers() []Run {
	tests := make([]test, 1001)

	for i := range tests {
		tests[i] = test{strconv.Itoa(i), wordify(i)}
	}

	return outputTests(shuffle(tests))
}
