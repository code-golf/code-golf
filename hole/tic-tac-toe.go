package hole

import (
	"fmt"
	"math/rand/v2"
	"strings"
)

var lines = [][]int{
	// horizontal
	{0, 1, 2}, {3, 4, 5}, {6, 7, 8},
	// vertical
	{0, 3, 6}, {1, 4, 7}, {2, 5, 8},
	// diagonal
	{0, 4, 8}, {2, 4, 6},
}

func ticTacToe() []Run {
	tests := make([]test, 0, 100)

	for range 100 {
		var argument strings.Builder

		board := generateBoard()

		for i := range 3 {
			argument.WriteString(fmt.Sprintln(string(board[i*3 : i*3+3])))
		}

		tests = append(tests, test{
			argument.String()[:10],
			determineWinner(board),
		})
	}

	return outputTests(shuffle(tests))
}

func determineWinner(board []rune) string {
	if isWonBy(board, 'X') {
		return "X"
	} else if isWonBy(board, 'O') {
		return "O"
	}

	return "-"
}

func generateBoard() []rune {
	board, players, X, O := []rune("         "), []rune{'X', 'O'}, 0, 0

	rand.Shuffle(len(players), func(i, j int) {
		players[i], players[j] = players[j], players[i]
	})

	for i := range 9 {
		if isWonBy(board, players[0]) || isWonBy(board, players[1]) {
			break
		}

		squares := []int{}

		for j, sq := range board {
			if sq == ' ' {
				squares = append(squares, j)
			}
		}

		if len(squares) == 0 {
			break
		}

		square := squares[rand.IntN(len(squares))]

		if i%2 == 0 {
			board[square] = 'X'
			X++
		} else {
			board[square] = 'O'
			O++
		}
	}

	if O > X {
		return generateBoard()
	}

	return board
}

func isWonBy(board []rune, ch rune) bool {
	for _, line := range lines {
		if ch == board[line[0]] && ch == board[line[1]] && ch == board[line[2]] {
			return true
		}
	}

	return false
}
