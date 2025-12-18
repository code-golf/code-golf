package hole

import "strings"

const alphabet = "abcdefghijklmnopqrstuvwxyz"

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

		// Separate the letters and words as arguments.
		args := make([]strings.Builder, 2)

		shuffle(dice)

		for j, row := range board {
			for k := range row {
				die := dice[k+j*len(row)]

				board[j][k] = die[randInt(0, 5)]

				args[0].WriteByte(board[j][k])

				if k < len(row)-1 {
					args[0].WriteByte(' ')
				}
			}

			if j < len(row)-1 {
				args[0].WriteByte('\n')
			}
		}

		var wordList []string

		// Determine all the valid words.
		for _, word := range shuffle(words) {
			if validate(board, strings.ToUpper(word)) {
				wordList = append(wordList, word)
			}
		}

		// Duplicate a number of valid words.
		for range 10 {
			wordList = append(wordList, randChoice(wordList))
		}

		// Generate a number of words that are too short.
		for j := range 10 {
			wordList = append(wordList, strings.Repeat(string(randChoice([]byte(alphabet))), j%2+1))
		}

		// Append a number of dictionary words. *** This needs more shizzle ***
		for j := 0; j < 25; {
			word := randWord()

			if len(word) < len(board)-1 || len(word) > len(board)*2-1 {
				continue
			}

			wordList = append(wordList, word)
			j++
		}

		var answer strings.Builder

		unique := make(map[string]bool)

		// Build the argument and the answer.
		for _, word := range shuffle(wordList) {
			if args[1].String() != "" {
				args[1].WriteByte(' ')
			}

			if !unique[word] && validate(board, strings.ToUpper(word)) {
				answer.WriteString(word)
				answer.WriteByte('\n')

				unique[word] = true
			}

			args[1].WriteString(word)
		}

		// Enforce no valid words for one of the runs if it didn't already generate.
		// If more than one such run occurs, regenerate the run so there's only one. But, unlikely.
		//
		// Implement case 'Q' and explain in description.

		// No valid words.
		if answer.String() == "" {
			answer.WriteByte('-')
		}

		answers[i] = Answer{Args: []string{args[0].String(), args[1].String()}, Answer: answer.String()}
	}

	return answers
})

func validate(b [4][4]byte, s string) bool {
	if l := len(b); len(s) < l-1 || len(s) > l*l {
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
		{+0, -1}, {+0, +1},
		{+1, -1}, {+1, +0}, {+1, +1},
	} {
		if dfs(b, used, s, i+1, j+direction[0], k+direction[1]) {
			return true
		}
	}

	used[j][k] = false

	return false
}
