package hole

import (
	"math/rand"
	"strings"
)

const (
	boardSize = 9
	blockSize = 3
)

func solve(board [boardSize][boardSize]int, cell int, count *int) bool {
	var i, j int

	for {
		if cell == boardSize*boardSize {
			*count++
			return *count == 2
		}

		i = cell / boardSize
		j = cell % boardSize

		if board[i][j] == 0 {
			break
		}

		cell++
	}

	// Origin of block.
	i0 := i - i%blockSize
	j0 := j - j%blockSize

	// 1 - 9 in random order.
	var numbers [boardSize]int

	for i := range numbers {
		j := rand.Intn(i + 1)
		numbers[i] = numbers[j]
		numbers[j] = i + 1
	}

Numbers:
	for _, number := range numbers {
		// number is already in the row.
		for _, numberInRow := range board[i] {
			if number == numberInRow {
				continue Numbers
			}
		}

		// number is already in the column.
		for _, row := range board {
			if number == row[j] {
				continue Numbers
			}
		}

		// number is already in the block.
		for _, row := range board[i0 : i0+blockSize] {
			for _, numberInRow := range row[j0 : j0+blockSize] {
				if number == numberInRow {
					continue Numbers
				}
			}
		}

		board[i][j] = number

		if solve(board, cell+1, count) {
			return true
		}
	}

	// No valid number for this cell, let's backtrack.
	board[i][j] = 0

	return false
}

func sudoku() (args []string, out string) {
	var board [boardSize][boardSize]int

	var generate func(int) bool
	generate = func(cell int) bool {
		i := cell / boardSize
		j := cell % boardSize

		// Origin of block.
		i0 := i - i%blockSize
		j0 := j - j%blockSize

		// 1 - 9 in random order.
		var numbers [boardSize]int

		for i := range numbers {
			j := rand.Intn(i + 1)
			numbers[i] = numbers[j]
			numbers[j] = i + 1
		}

	Numbers:
		for _, number := range numbers {
			// number is already in the row.
			for _, numberInRow := range board[i] {
				if number == numberInRow {
					continue Numbers
				}
			}

			// number is already in the column.
			for _, row := range board {
				if number == row[j] {
					continue Numbers
				}
			}

			// number is already in the block.
			for _, row := range board[i0:i] {
				for _, numberInRow := range row[j0 : j0+blockSize] {
					if number == numberInRow {
						continue Numbers
					}
				}
			}

			board[i][j] = number

			if cell+1 == boardSize*boardSize || generate(cell+1) {
				return true
			}
		}

		// No valid number for this cell, let's backtrack.
		board[i][j] = 0

		return false
	}

	generate(0)

	var b strings.Builder

	for i, row := range board {
		if i == 0 {
			b.WriteString("┏━━━┯━━━┯━━━┳━━━┯━━━┯━━━┳━━━┯━━━┯━━━┓\n")
		} else if i%blockSize == 0 {
			b.WriteString("┣━━━┿━━━┿━━━╋━━━┿━━━┿━━━╋━━━┿━━━┿━━━┫\n")
		} else {
			b.WriteString("┠───┼───┼───╂───┼───┼───╂───┼───┼───┨\n")
		}

		for j, number := range row {
			if j%blockSize == 0 {
				b.WriteString("┃")
			} else {
				b.WriteString("│")
			}

			b.WriteRune(' ')
			if number == 0 {
				b.WriteRune(' ')
			} else {
				b.WriteRune(rune('0' + number))
			}
			b.WriteRune(' ')
		}

		b.WriteString("┃\n")
	}

	b.WriteString("┗━━━┷━━━┷━━━┻━━━┷━━━┷━━━┻━━━┷━━━┷━━━┛")

	out = b.String()

	for k, cell := range rand.Perm(boardSize * boardSize) {
		i := cell / boardSize
		j := cell % boardSize

		orig := board[i][j]
		board[i][j] = 0

		var count int
		solve(board, 0, &count)

		// Removing this cell creates too many solutions, put it back.
		if count == 2 {
			board[i][j] = orig
		}

		// TODO Switch to dancing links to remove this grotesque hack!
		// http://garethrees.org/2007/06/10/zendoku-generation/#section-4
		if k == 50 {
			break
		}
	}

	for _, row := range board {
		var b strings.Builder

		for _, number := range row {
			if number == 0 {
				b.WriteRune('_')
			} else {
				b.WriteRune(rune('0' + number))
			}
		}

		args = append(args, b.String())
	}

	return
}
