package hole

import (
	"math/rand/v2"
	"strings"
)

func boggle() []Run {
	const argc = 100 // Preserve original argc

	tests := make([]test, 0, argc)

	for range 2 * argc {
		var argument, expected strings.Builder

		board, count, words := scramble(), randInt(2 * argc, 4 * argc), make(map[string]struct{})

		for i := 0; i < count; {
			word := randWord()

			if _, dupe := words[word]; dupe {
				continue
			}

			words[word] = struct{}{}
			i++

			argument.WriteString(strings.ToLower(word))

			if i < count {
				argument.WriteRune(' ')
			}
		}

		for _, word := range strings.Fields(argument.String()) {
			if validate(board, strings.ToUpper(word)) {
				expected.WriteString(word)
				expected.WriteRune(' ')
			}
		}

		if len(expected.String()) == 0 {
			expected.WriteRune('-')
		}

		argument.WriteString("\n" + stringify(board))

		tests = append(tests, test{
			argument.String(),
			expected.String(),
		})
	}

	tests = shuffle(tests)

	return outputTests(tests[:argc], tests[len(tests) - argc:])
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
	var dice = [16][]rune{
		[]rune("RIFOBX"), []rune("IFEHEY"), []rune("DENOWS"), []rune("UTOKND"),
		[]rune("HMSRAO"), []rune("LUPETS"), []rune("ACITOA"), []rune("YLGKUE"),
		[]rune("QBMJOA"), []rune("EHISPN"), []rune("VETIGN"), []rune("BALIYT"),
		[]rune("EZAVND"), []rune("RALESC"), []rune("UWILRG"), []rune("PACEMD"),
	}

	for _, die := range dice {
		// Scramble positions.
		rand.Shuffle(len(dice), func(i, j int) {
			dice[i], dice[j] = dice[j], dice[i]
		})

		// Scramble letters.
		rand.Shuffle(len(die), func(i, j int) {
			die[i], die[j] = die[j], die[i]
		})
	}

	var board [4][4]rune

	for i := range 4 {
		for j := range 4 {
			board[i][j] = dice[j+i*4][0] // Index 0 represents the side facing up after scrambling.
		}
	}

	return board
}

func validate(board [4][4]rune, word string) bool {
	if len(word) < 3 || len(word) > 16 {
		return false
	}

	letters, used := []rune(word), [4][4]bool{}

	for i := range 4 {
		for j := range 4 {
			if board[i][j] == letters[0] && dfs(board, used, letters, 0, i, j) {
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

	if i < 0 || i > 3 || j < 0 || j > 3 || used[i][j] {
		return false
	}

	if letter := board[i][j]; letter == 'q' {
		if index + 1 < len(word) && word[index] == 'q' && word[index + 1] == 'u' {
			used[i][j] = true

			for _, direction := range [][]int{
				{-1, -1}, {-1, +0}, {-1, +1},
				{+0, -1}, {+0, +1},
				{+1, -1}, {+1, +0}, {+1, +1},
			} {
				if dfs(board, used, word, index + 2, i + direction[0], j + direction[1]) {
					return true
				}
			}

			used[i][j] = false

			return false
		}
	} else if letter != word[index] {
		return false
	}

	used[i][j] = true

	for _, direction := range [][]int{
		{-1, -1}, {-1, +0}, {-1, +1},
		{+0, -1}, {+0, +1},
		{+1, -1}, {+1, +0}, {+1, +1},
	} {
		x, y := i+direction[0], j+direction[1]

		if dfs(board, used, word, index+1, x, y) {
			return true
		}
	}

	used[i][j] = false

	return false
}
