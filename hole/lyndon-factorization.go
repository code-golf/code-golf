package hole

import (
	"math/rand/v2"
	"strings"
)

var fixedInputs = []string{
	"lyndon", "factorization", "codegolf",
	"abcdefghijklmnopqrstuvwxyzyxwvutsrqponmlkjihgfedcbabcdefghijklmnopqrstuvwxyzyxwvutsrqponmlkjihgfedcba",
	"vdsxqiytqmptjqintoblcromtgjpujvhrwjoydwgbwezpfxwlvye",
	"rxexxecsudldbuosnzkvzvhlyofxczphartivkiehefdaffazlle",
	"dntkwgzapobnhyduatjqmkcfidgsdcnastyosbasaolycziwfphq",
	"aaaaaaaaaa",
}

var _ = answerFunc("lyndon-factorization", func() []Answer {
	const alphabet = "qwertzuiopasdfghjklyxcvbnm"

	tests := make([]test, 100)

	for i, input := range fixedInputs {
		tests[i] = lyndonFactorizationTest(input)
	}

	for i := len(fixedInputs); i < len(tests); i++ {
		testLength := 1 + rand.IntN(30)
		input := ""
		for range testLength {
			j := randInt(0, 25)
			input += alphabet[j : j+1]
		}
		tests[i] = lyndonFactorizationTest(input)
	}
	return outputTests(shuffle(tests))
})

func lyndonFactorizationTest(input string) test {
	in := input

	var out strings.Builder
	i := 0

	for i < len(input) {
		j := i + 1
		k := i

		for j < len(input) && input[k] <= input[j] {
			if input[k] < input[j] {
				k = i
			} else {
				k++
			}
			j++
		}

		for i <= k {
			out.WriteString(input[i : i+j-k])
			out.WriteByte(' ')
			i += j - k
		}

	}

	return test{in, strings.TrimSpace(out.String())}
}
