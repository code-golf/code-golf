package hole

import "strings"

var _ = answerFunc("transpose-sentence", func() []Answer {
	tests := make([]test, randInt(100, 130))

	for i := range tests {
		var sbr strings.Builder

		// Four to seven random words.
		for j := range randInt(4, 7) {
			if j != 0 {
				sbr.WriteByte(' ')
			}

			sbr.WriteString(randWord())
		}

		sentence := sbr.String()
		tests[i] = test{sentence, transpose(sentence)}
	}

	// Ensure a delta of at least 10.
	sentence := "concentration is essential for success"
	tests = append(tests, test{sentence, transpose(sentence)})

	return outputTests(shuffle(tests))
})

func transpose(s string) string {
	srcWords := strings.Fields(s)
	maxLen := 0

	for _, w := range srcWords {
		maxLen = max(maxLen, len(w))
	}

	i, l, c := 0, len(srcWords), true
	sbrs := make([]strings.Builder, maxLen)

	for c {
		c = false
		for wi := 0; wi < l; wi++ {
			if w := srcWords[wi]; i < len(w) {
				sbrs[i].WriteByte(w[i])
				c = true
			}
		}
		i++
	}

	var transposed strings.Builder

	for _, sbr := range sbrs {
		transposed.WriteString(sbr.String())
		transposed.WriteByte(' ')
	}

	return strings.TrimSpace(transposed.String())
}
