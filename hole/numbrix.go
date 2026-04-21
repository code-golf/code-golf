package hole

import (
	"fmt"
	"strconv"
	"strings"
)

func init() {
	judge("numbrix", numbrixJudge)
}

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

func numbrixJudge(run Run) string {
	parseGrid := func(s string) ([boardSize][boardSize]int, bool) {
		var grid [boardSize][boardSize]int
		row := 0
		for _, line := range strings.Split(s, "\n") {
			if !strings.ContainsRune(line, '┃') {
				continue
			}
			line = strings.ReplaceAll(line, "┃", "│")
			cells := strings.Split(line, "│")
			if len(cells) < boardSize+2 {
				return grid, false
			}
			for col := range boardSize {
				cell := strings.TrimSpace(cells[col+1])
				if cell == "" {
					grid[row][col] = 0
				} else {
					n, err := strconv.Atoi(cell)
					if err != nil {
						return grid, false
					}
					grid[row][col] = n
				}
			}
			row++
		}
		return grid, row == boardSize
	}

	userGrid, ok := parseGrid(run.Stdout)
	if !ok {
		return run.Answer
	}
	inputGrid, _ := parseGrid(run.Args[0])

	// Validate: all 1–81 present, clues respected, consecutive numbers adjacent.
	valid := true
	var pos [cellCount + 1][2]int
	var seen [cellCount + 1]bool
validate:
	for i := range boardSize {
		for j := range boardSize {
			v := userGrid[i][j]
			if v < 1 || v > cellCount || seen[v] {
				valid = false
				break validate
			}
			seen[v] = true
			pos[v] = [2]int{i, j}
			if inputGrid[i][j] != 0 && inputGrid[i][j] != v {
				valid = false
				break validate
			}
		}
	}
	if valid {
		for n := 1; n < cellCount; n++ {
			dr := pos[n][0] - pos[n+1][0]
			dc := pos[n][1] - pos[n+1][1]
			if dr < 0 {
				dr = -dr
			}
			if dc < 0 {
				dc = -dc
			}
			if dr+dc != 1 {
				valid = false
				break
			}
		}
	}
	if valid {
		return run.Stdout
	}

	// Try to find a valid solution close to the user's output, so the diff
	// highlights the user's actual mistakes rather than comparing against
	// an unrelated solution.
	if match, ok := closestNumbrixMatch(inputGrid, userGrid); ok {
		return printNumbrix(match)
	}

	return run.Answer
}

// closestNumbrixMatch tries to build a valid Numbrix solution that overlaps
// with userGrid by preferring the user's chosen positions for each number
// (without violating the clues in input). Returns the first valid solution
// found, or false if none exists.
func closestNumbrixMatch(input, userGrid [boardSize][boardSize]int) ([boardSize][boardSize]int, bool) {
	// Preferred position per number: clues always win over user's guess.
	var preferred [cellCount + 1][2]int
	var hasPref [cellCount + 1]bool
	for i := range boardSize {
		for j := range boardSize {
			if v := input[i][j]; v > 0 {
				preferred[v] = [2]int{i, j}
				hasPref[v] = true
			}
		}
	}
	for i := range boardSize {
		for j := range boardSize {
			v := userGrid[i][j]
			if v >= 1 && v <= cellCount && !hasPref[v] && input[i][j] == 0 {
				preferred[v] = [2]int{i, j}
				hasPref[v] = true
			}
		}
	}

	// Bound the search so a pathological input can't stall the judge.
	const maxNodes = 200_000
	nodes := 0

	var result [boardSize][boardSize]int
	var dfs func(r, c, n int) bool
	dfs = func(r, c, n int) bool {
		nodes++
		if nodes > maxNodes {
			return false
		}
		if r < 0 || r >= boardSize || c < 0 || c >= boardSize {
			return false
		}
		if result[r][c] != 0 {
			return false
		}
		if input[r][c] != 0 && input[r][c] != n {
			return false
		}

		result[r][c] = n
		if n == cellCount {
			return true
		}

		dirs := [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
		if hasPref[n+1] {
			tr, tc := preferred[n+1][0], preferred[n+1][1]
			for i, d := range dirs {
				if r+d[0] == tr && c+d[1] == tc {
					dirs[0], dirs[i] = dirs[i], dirs[0]
					break
				}
			}
		}
		for _, d := range dirs {
			if dfs(r+d[0], c+d[1], n+1) {
				return true
			}
		}
		result[r][c] = 0
		return false
	}

	if hasPref[1] && dfs(preferred[1][0], preferred[1][1], 1) {
		return result, true
	}
	for i := range boardSize {
		for j := range boardSize {
			if hasPref[1] && i == preferred[1][0] && j == preferred[1][1] {
				continue
			}
			if dfs(i, j, 1) {
				return result, true
			}
		}
	}
	return result, false
}

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
				fmt.Fprintf(&s, " %2d ", number)
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

	return dfs(i, j, 2)
}
