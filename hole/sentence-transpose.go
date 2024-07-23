package hole

import (
	"math/rand"
	"strings"
)

const minWordCount = 4
const wordCountVariance = 4
const minCharCount = 4
const charCountVariance = 4
const chars = "abcdefghijklmnopqrstuvwxyz"

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

func transposeSentence() []Run {
	count := 100 + rand.Intn(30)

	tests := make([]test, count)

	for i := 0; i < count; i++ {
		sbr := strings.Builder{}
		nWord := minWordCount - 1 + rand.Intn(wordCountVariance)

		for j := 0; j < nWord; j++ {
			sbr.WriteString(randString(minCharCount + rand.Intn(charCountVariance)))
			sbr.WriteString(" ")
		}
		sbr.WriteString(randString(minCharCount + rand.Intn(charCountVariance)))

		sentence := sbr.String()

		tests[i] = test{sentence, transpose(sentence)}
	}

	return outputTests(tests)
}

func transpose(s string) string {
	srcWords := strings.Fields(s)
	maxLen := 0

	for _, w := range srcWords {
		if l := len(w); l > maxLen {
			maxLen = l
		}
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

	transposed := strings.Builder{}

	for _, sbr := range sbrs {
		transposed.WriteString(sbr.String())
		transposed.WriteString(" ")
	}

	return strings.TrimSpace(transposed.String())
}
