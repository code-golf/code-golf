package hole

import (
	"fmt"
	"math/rand/v2"
	"strings"
)

func ticTacToe() []Run {
	boards, tests := make(map[string]struct{}), make([]test, 0, 500)

	for i := 0; i < 500; {
		var argument strings.Builder

		board := generateBoard()

		for j := 0; j < 9; j += 3 {
			argument.WriteString(fmt.Sprintln(string(board[j : j+3])))
		}

		if _, dupe := boards[argument.String()]; dupe {
			continue
		}

		boards[argument.String()] = struct{}{}
		i++

		tests = append(tests, test{
			argument.String(),
			announceWinner(board),
		})
	}

	return outputTests(shuffle(tests))
}

func generateBoard() []rune {
	board, characters, X, O := []rune("........."), []rune{'X', 'O'}, 0, 0

	rand.Shuffle(len(characters), func(i, j int) {
		characters[i], characters[j] = characters[j], characters[i]
	})

	for i := range 9 {
		if isGameOver(board, characters[0]) || isGameOver(board, characters[1]) {
			break
		}

		dots := []int{}

		for j, ch := range board {
			if ch == '.' {
				dots = append(dots, j)
			}
		}

		if len(dots) == 0 {
			break
		}

		index := dots[rand.IntN(len(dots))]

		if i%2 == 0 {
			board[index] = 'X'
			X++
		} else {
			board[index] = 'O'
			O++
		}
	}

	return board
}

func isGameOver(board []rune, ch rune) bool {
	for _, pattern := range [][]int{
		{0, 1, 2}, {3, 4, 5}, {6, 7, 8},
		{0, 3, 6}, {1, 4, 7}, {2, 5, 8},
		{0, 4, 8}, {2, 4, 6},
	} {
		if ch == board[pattern[0]] && ch == board[pattern[1]] && ch == board[pattern[2]] {
			return true
		}
	}

	return false
}

func announceWinner(board []rune) string {
	if isGameOver(board, 'X') {
		return "X"
	} else if isGameOver(board, 'O') {
		return "O"
	}

	return "-"
}
