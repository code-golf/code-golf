package hole

import "strings"

var dice = []string{
	"AACIOT", "ABILTY", "ABJMOQ", "ACDEMP",
	"ACELRS", "ADENVZ", "AHMORS", "BIFORX",
	"DENOSW", "DKNOTU", "EEFHIY", "EGINTV",
	"EGKLUY", "EHINPS", "ELPSTU", "GILRUW",
}

var _ = answerFunc("boggle", func() []Answer {
	answers := make([]Answer, 5)

	for i := range answers {
		var board [4][4]byte

		args := make([]strings.Builder, 2)

		shuffle(dice)

		for j, row := range board {
			for k := range row {
				die := dice[k + j * len(row)]

				board[j][k] = die[randInt(0, 5)]

				args[0].WriteByte(board[j][k])

				if k < len(row) - 1 {
					args[0].WriteByte(' ')
				}
			}

			if j < len(row) - 1 {
				args[0].WriteByte('\n')
			}
		}

		var answer strings.Builder

		for _, word := range shuffle(words) {
			if validate(board, strings.ToUpper(word)) {
				if answer.String() != "" {
					args[1].WriteByte(' ')
				}

				answer.WriteString(word)
				answer.WriteByte('\n')

				args[1].WriteString(word)
			}
		}

		// Expand args[1] with random:
		// - invalid words less than 3 characters long
		// - invalid words because of all the rest
		// - valid but duplicate words (because no duplicate prints)
		//
		// Enforce no valid words for one of the runs if it didn't already generate.
		// If more than one such run occurs, regenerate the run so there's only one. But, unlikely.
		//
		// Implement the 'Q' thing.
		//
		// Update more description.

		if answer.String() == "" {
			answer.WriteByte('-')
		}

		answers[i] = Answer{Args: []string{args[0].String(), args[1].String()}, Answer: answer.String()}
	}

	return answers
})

func validate(b [4][4]byte, s string) bool {
	if l := len(b); len(s) < l - 1 || len(s) > l * l {
		return false
	}

	used := [4][4]bool{}

	for i, row := range b {
		for j, letter := range row {
			if dfs(b, used, s, 0, i, j) && s[0] == letter {
				return true
			}
		}
	}

	return false
}

func dfs(b [4][4]byte, used [4][4]bool, s string, i, j, k int) bool {
	if len(s) == i {
		return true
	}

	if l := len(b) - 1; j < 0 || j > l || k < 0 || k > l || used[j][k] || s[i] != b[j][k] {
		return false
	}

	used[j][k] = true

	for _, direction := range [][2]int{
		{-1, -1}, {-1, +0}, {-1, +1},
		{+0, -1},           {+0, +1},
		{+1, -1}, {+1, +0}, {+1, +1},
	} {
		if dfs(b, used, s, i+1, j+direction[0], k+direction[1]) {
			return true
		}
	}

	used[j][k] = false

	return false
}
