package hole

import (
	"fmt"
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

var numbrixSquares = shuffle([][2]int{
	{0, 0}, {0, 8}, {1, 1}, {1, 7}, {2, 2}, {2, 6}, {3, 3}, {3, 5}, {4, 0},
	{4, 4}, {5, 3}, {5, 5}, {6, 2}, {6, 6}, {7, 1}, {7, 7}, {8, 0}, {8, 8},
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
	sq, connects := randChoice(numbrixSquares), func() bool {
		n, r, c := 0, -1, -1

		for i, row := range puzzle {
			for j, number := range row {
				if number == 0 {
					n++

					if r < 0 {
						r, c = i, j
					}
				}
			}
		}

		if n == 0 {
			return true
		}

		var used [boardSize][boardSize]bool

		squares, t := [][2]int{{r, c}}, 1

		used[r][c] = true

		for i := 0; i < len(squares); i++ {
			r, c := squares[i][0], squares[i][1]

			for _, d := range [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
				j, k := r+d[0], c+d[1]

				if j >= 0 && j < boardSize && k >= 0 && k < boardSize && !used[j][k] && puzzle[j][k] == 0 {
					squares = append(squares, [2]int{j, k})
					t++
					used[j][k] = true
				}
			}
		}

		return n == t
	}

	var dfs func(i, j, k int) bool

	dfs = func(i, j, k int) bool {
		if k > cellCount {
			return true
		}

		type Direction struct{ i, j, k int }

		dirs := make([]Direction, 0, 4)

		add := func(r, c int) {
			if r >= 0 && r < boardSize && c >= 0 && c < boardSize && puzzle[r][c] == 0 {
				d := 0

				if r > 0 && puzzle[r-1][c] == 0 {
					d++
				}

				if r < boardSize-1 && puzzle[r+1][c] == 0 {
					d++
				}

				if c > 0 && puzzle[r][c-1] == 0 {
					d++
				}

				if c < boardSize-1 && puzzle[r][c+1] == 0 {
					d++
				}

				dirs = append(dirs, Direction{r, c, d})
			}
		}

		add(i-1, j)
		add(i+1, j)
		add(i, j-1)
		add(i, j+1)

		if len(dirs) == 0 {
			return false
		}

		for i := 0; i < len(dirs); i++ {
			for j := i; j > 0 && (dirs[j].k < dirs[j-1].k || dirs[j].k == dirs[j-1].k && randBool()); j-- {
				dirs[j-1], dirs[j] = dirs[j], dirs[j-1]
			}
		}

		for _, d := range dirs {
			puzzle[d.i][d.j] = k

			if connects() && dfs(d.i, d.j, k+1) {
				return true
			}

			puzzle[d.i][d.j] = 0
		}

		return false
	}

	i, j := sq[0], sq[1]

	puzzle[i][j] = 1

	if dfs(i, j, 2) {
		return true
	}

	return false
}
