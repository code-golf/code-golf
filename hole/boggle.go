package hole

import (
	"math/rand/v2"
	"strings"
)

// TODO Figure out dictionary whats and whatnots
// TODO 'Qu' implementation

func boggle() []Run {
	tests := make([]test, 0, 100)

	for range 100 {
		var argument, expected strings.Builder

		board, words := scramble(), make(map[string]struct{})

		for i := 0; i < 225; {
			word := strings.ToUpper(randWord())

			if _, dupe := words[word]; dupe {
				continue
			}

			words[word] = struct{}{}
			i++

			argument.WriteString(word)

			if i < 225 {
				argument.WriteRune(' ')
			}
		}

		for _, word := range strings.Fields(argument.String()) {
			if validate(board, word) {
				expected.WriteString(word + " ")
			}
		}

		if expected.String() == "" {
			expected.WriteRune('-')
		}

		argument.WriteString("\n\n" + stringify(board))

		tests = append(tests, test{
			argument.String(),
			expected.String(),
		})
	}

	return outputTests(shuffle(tests))
}

func stringify(board [4][4]rune) string {
	var grid strings.Builder

	for i := range 4 {
		for j := range 4 {
			grid.WriteRune(board[i][j])

			if j < 3 {
				grid.WriteRune(' ')
			} else {
				grid.WriteRune('\n')
			}
		}
	}

	return grid.String()
}

func scramble() [4][4]rune {
	var cubes = [16][]rune{
		[]rune("RIFOBX"), []rune("IFEHEY"), []rune("DENOWS"), []rune("UTOKND"),
		[]rune("HMSRAO"), []rune("LUPETS"), []rune("ACITOA"), []rune("YLGKUE"),
		[]rune("QBMJOA"), []rune("EHISPN"), []rune("VETIGN"), []rune("BALIYT"),
		[]rune("EZAVND"), []rune("RALESC"), []rune("UWILRG"), []rune("PACEMD"),
	}

	for _, cube := range cubes {
		// Scramble positions.
		rand.Shuffle(len(cubes), func(i, j int) {
			cubes[i], cubes[j] = cubes[j], cubes[i]
		})

		// Scramble letters.
		rand.Shuffle(len(cube), func(i, j int) {
			cube[i], cube[j] = cube[j], cube[i]
		})
	}

	var board [4][4]rune

	for i := range 4 {
		for j := range 4 {
			board[i][j] = cubes[j + i * 4][0] // Index 0 represents the side facing up after scrambling.
		}
	}

	return board
}

func validate(board [4][4]rune, word string) bool {
	if len(word) < 3 {
		return false
	}

	crossed, letters := [4][4]bool{}, []rune(word)

	for i := range 4 {
		for j := range 4 {
			if board[i][j] == letters[0] && dfs(board, crossed, letters, 0, i, j) {
				return true
			}
		}
	}

	return false
}

func dfs(board [4][4]rune, used [4][4]bool, word []rune, index, i, j int) bool {
	if index == len(word) {
		return true
	}

	if i < 0 || i > 3 || j < 0 || j > 3 || used[i][j] || board[i][j] != word[index] {
		return false
	}

	used[i][j] = true

	for _, direction := range [][]int{
		{-1, -1}, {-1, +0}, {-1, +1},
		{+0, -1}, {+0, +0}, {+0, +1},
		{+1, -1}, {+1, +0}, {+1, +1},
	} {
		x, y := i + direction[0], j + direction[1]

		if dfs(board, used, word, index + 1, x, y) {
			return true
		}
	}

	return false
}
