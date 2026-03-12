package hole

import (
	"math/rand/v2"
	"strings"
)

const (
	rectSize = 10
	rectSquare = rectSize * rectSize
)

var _ = answerFunc("rectangle-fill", func() []Answer {
	answers := make([]Answer, 3)

	for i := range answers {
		var rect [rectSize][rectSize]byte

		initRect(&rect)

		argument, expected := []string{drawRect(rect)}, fillRect(rect)

		answers[i] = Answer{Args: argument, Answer: expected}
	}

	return answers
})

func initRect(rect *[rectSize][rectSize]byte) {
	for i, row := range rect {
		for j := range row {
			rect[i][j] = '0'
		}
	}

	for i := 1; i < rectSize; {
		var top, bot = rand.IntN(rectSquare), rand.IntN(rectSquare)

		if top == bot {
			continue
		}

		tr, tc := top/rectSize, top%rectSize
		br, bc := bot/rectSize, bot%rectSize

		if rect[tr][tc] != '0' || rect[br][bc] != '0' {
			continue
		}

		rect[tr][tc] += byte(i)
		rect[br][bc] += byte(i)

		i++
	}
}

func drawRect(rect [rectSize][rectSize]byte) string {
	var s strings.Builder

	for _, row := range rect {
		for _, ch := range row {
			s.WriteByte(ch)
		}

		s.WriteByte('\n')
	}

	return strings.Trim(s.String(), "\n")
}

func fillRect(rect [rectSize][rectSize]byte) string {
	for i := 1; i < rectSize; i++ {
		var digit = byte('0' + i)

		tr, tc := rectSize, rectSize
		br, bc := 0, 0

		for j, row := range rect {
			for k, ch := range row {
				if ch == digit {
					if j < tr {
						tr = j
					}

					if j > br {
						br = j
					}

					if k < tc {
						tc = k
					}

					if k > bc {
						bc = k
					}
				}
			}
		}

		if br == 0 {
			continue
		}

		for j := tr; j <= br; j++ {
			for k := tc; k <= bc; k++ {
				if digit > rect[j][k] {
					rect[j][k] = digit
				}
			}
		}
	}

	return drawRect(rect)
}
