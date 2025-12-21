package hole

import (
	"maps"
	"slices"
	"strings"
)

const gridSize = 4

var dice = []string{
	"AACIOT", "ABILTY", "ABJMOQ", "ACDEMP",
	"ACELRS", "ADENVZ", "AHMORS", "BIFORX",
	"DENOSW", "DKNOTU", "EEFHIY", "EGINTV",
	"EGKLUY", "EHINPS", "ELPSTU", "GILRUW",
}

var _ = answerFunc("boggle", func() []Answer {
	answers := make([]Answer, gridSize)

	for i := range answers {
		var args, answer strings.Builder

		shuffle(dice)

		dictionary, letters := slices.Clone(words), make(map[byte]int)

		var grid [gridSize][gridSize]byte

		// Build "grid" argument.
		for r, row := range grid {
			for c := range row {
				letter := randChoice([]byte(dice[r*gridSize+c]))

				grid[r][c] = letter
				letters[letter]++
				args.WriteByte(letter)

				if c < gridSize-1 {
					args.WriteByte(' ')
				}
			}

			if r < gridSize-1 {
				args.WriteByte('\n')
			}
		}

		answers[i].Args = []string{args.String()}

		var words []string

	outer: // Identify all valid and nearly valid words.
		for _, word := range shuffle(dictionary) {
			if used, uses := validate(grid, strings.ToUpper(word)), maps.Clone(letters); len(word)-used <= 2 {
				// Before adding the word, first check if the unused letters appear in the grid.
				for _, letter := range strings.ToUpper(word)[used:] {
					if uses[byte(letter)]--; uses[byte(letter)] <= 0 {
						continue outer
					}
				}

				words = append(words, word)

				// A perfectly valid word.
				if len(word) == used {
					answer.WriteString(word)
					answer.WriteByte('\n')
				}
			}
		}

		// No valid words.
		if answer.String() == "" {
			answer.WriteByte('-')
		}

		answers[i].Answer = answer.String()

		const alphabet = "abcdefghijklmnopqrstuvwxyz"

		// Add any letter of the alphabet.
		words = append(words, string(randChoice([]byte(alphabet))))

		// Add a combination of any two unique letters from the alphabet.
		words = append(words, string(shuffle([]byte(alphabet))[:2]))

		// Add a copy of any valid word.
		words = append(words, randChoice(strings.Fields(answers[i].Answer)))

		// Add any word from the dictionary.
		words = append(words, randWord())

		args.Reset()

		// Build "words" argument.
		for _, word := range shuffle(words) {
			if args.String() != "" {
				args.WriteByte(' ')
			}

			args.WriteString(word)
		}

		answers[i].Args = append(answers[i].Args, args.String())
	}

	return answers
})

func validate(grid [gridSize][gridSize]byte, word string) int {
	var used int

	if len(word) >= gridSize-1 && len(word) < gridSize*gridSize {
		var uses [gridSize][gridSize]bool

		for r, row := range grid {
			for c := range row {
				if n := dfs(grid, uses, word, 0, r, c); n > used {
					used = n
				}
			}
		}
	}

	return used
}

func dfs(grid [gridSize][gridSize]byte, uses [gridSize][gridSize]bool, word string, index, row, col int) int {
	switch true {
	case index == len(word), row < 0, row > gridSize-1, col < 0, col > gridSize-1, uses[row][col], grid[row][col] != word[index]:
		return 0
	}

	var used int

	uses[row][col] = true

	for _, direction := range [...][2]int{
		{-1, -1}, {-1, +0}, {-1, +1},
		{+0, -1}, {+0, +1},
		{+1, -1}, {+1, +0}, {+1, +1},
	} {
		if n := dfs(grid, uses, word, index+1, row+direction[0], col+direction[1]); n > used {
			used = n
		}
	}

	uses[row][col] = false

	return used + 1
}
