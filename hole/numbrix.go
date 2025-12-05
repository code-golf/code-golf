package hole

import (
	"fmt"
	_ "math"
	_ "math/rand"
	"strings"
)

// This implementation uses the constant variables defined in the Sudoku holes.

var numbrixPatterns = shuffle([]string{
	"000000000000101000001010100010000010001000100010000010001010100000101000000000000",
	"000000000001101100010000010010000010000000000010000010010000010001101100000000000",
	"000000000010101010000000000010000010000010000010000010000000000010101010000000000",
	"000000000010101010001010100010000010001000100010000010001010100010101010000000000",
	"000000000011010110010000010000000000010000010000000000010000010011010110000000000",
	"000010000001000100011000110000000000100000001000000000011000110001000100000010000",
	"001010100000010000101010101000101000111000111000101000101010101000010000001010100",
	"100000001000101000000101000011000110000000000011000110000101000000101000100000001",
	"100010001010010010000000000000000000110000011000000000000000000010010010100010001",
	"101010101000000000100000001000000000100000001000000000100000001000000000101010101",
})

var _ = answerFunc("numbrix", func() []Answer {
	answers := make([]Answer, blockSize)

	for i := range answers {
		var puzzle [boardSize][boardSize]int

		solveNumbrix(&puzzle)

		expected := printNumbrix(puzzle)

		for j, ch := range randChoice(numbrixPatterns) {
			if ch == '0' {
				puzzle[j/boardSize][j%boardSize] = 0
			}
		}

		argument := []string{printNumbrix(puzzle)}

		answers[i] = Answer{Args: argument, Answer: expected}
	}

	return answers
})

func printNumbrix(puzzle [boardSize][boardSize]int) string {
	var s strings.Builder

	for i, row := range puzzle {
		if i == 0 {
			s.WriteString("┏━━━━┯━━━━┯━━━━┯━━━━┯━━━━┯━━━━┯━━━━┯━━━━┯━━━━┓\n")
		} else {
			s.WriteString("┠────┼────┼────┼────┼────┼────┼────┼────┼────┨\n")
		}

		for j, number := range row {
			if j == 0 {
				s.WriteRune('┃')
			} else {
				s.WriteRune('│')
			}

			if number == 0 {
				s.WriteString("    ")
			} else {
				s.WriteString(fmt.Sprintf(" %2d ", number))
			}
		}

		s.WriteString("┃\n")
	}

	s.WriteString("┗━━━━┷━━━━┷━━━━┷━━━━┷━━━━┷━━━━┷━━━━┷━━━━┷━━━━┛")

	return s.String()
}

func solveNumbrix(puzzle *[boardSize][boardSize]int) bool {
	// Tis generating, but needs optimizing.

	return false
}
