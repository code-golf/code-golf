package hole

import (
	"fmt"
	"maps"
	"strings"
)

const gridSize = 4

var _ = answerFunc("boggle", func() []Answer {
	answers := make([]Answer, gridSize+1)

	for i := 0; i < len(answers); {
		var grid [gridSize][gridSize]byte

		dice := scramble(&grid)

		// Force two runs involving "Qu".
		if i%2 == dice['q'] {
			continue
		}

		var args, answer strings.Builder

		// Build "grid" argument.
		for _, row := range grid {
			for c, letter := range row {
				if args.WriteByte(letter); c < gridSize-1 {
					args.WriteByte(' ')
				}
			}

			fmt.Fprintln(&args)
		}

		dictionary := shuffle(words)

		var words []string

	outer: // Identify all valid and nearly valid words.
		for _, word := range dictionary {
			if letters, unused := maps.Clone(dice), validate(grid, word); unused <= 2 {
				// Make sure every unused letter appears in the grid.
				for _, letter := range word[len(word)-unused:] {
					if letters[byte(letter)]--; letters[byte(letter)] <= 0 {
						continue outer
					}
				}

				if unused != 0 {
					// Force a run with no valid words.
					if words = append(words, word); i == 1 {
						continue
					}
				} else if i != 1 {
					// A perfectly valid word.
					fmt.Fprintln(&answer, word)
					words = append(words, word)
				}
			}
		}

		answers[i].Args = []string{strings.ToUpper(args.String())}

		if args.Reset(); answer.String() != "" {
			// Add a copy of any valid word.
			words = append(words, randChoice(strings.Fields(answer.String())))
		} else {
			// No valid words.
			answer.WriteByte('-')
		}

		const alphabet = "abcdefghijklmnopqrstuvwxyz"

		// Add any letter.
		words = append(words, string(randChoice([]byte(alphabet))))

		// Add a combination of any two unique letters.
		words = append(words, string(shuffle([]byte(alphabet))[:2]))

		// Add four more random words.
		if i != 1 {
			words = append(words, randWord(), randWord(), randWord(), randWord())
		}

		// Build "words" argument.
		args.WriteString(strings.Join(shuffle(words), " "))

		answers[i].Args = append(answers[i].Args, args.String())
		answers[i].Answer = answer.String()

		i++
	}

	return shuffle(answers)
})

func scramble(grid *[gridSize][gridSize]byte) map[byte]int {
	dice, letters := shuffle([]string{
		"aaciot", "ahmors", "egkluy", "abilty",
		"acdemp", "egintv", "gilruw", "elpstu",
		"denosw", "acelrs", "abjmoq", "eefhiy",
		"ehinps", "dknotu", "adenvz", "biforx",
	}), make(map[byte]int)

	for r, row := range grid {
		for c := range row {
			letter := randChoice([]byte(dice[r*gridSize+c]))

			if letters[letter]++; letter == 'q' {
				letters['u']++
			}

			grid[r][c] = letter
		}
	}

	return letters
}

func validate(grid [gridSize][gridSize]byte, word string) int {
	var used int

	if len(word) >= gridSize-1 {
		var uses [gridSize][gridSize]bool

		for r, row := range grid {
			for c := range row {
				if n := dfs(&grid, &uses, word, 0, r, c); n > used {
					used = n
				}
			}
		}
	}

	return len(word) - used
}

func dfs(grid *[gridSize][gridSize]byte, uses *[gridSize][gridSize]bool, word string, index, r, c int) int {
	switch true {
	case index >= len(word), r < 0, r >= gridSize, c < 0, c >= gridSize, uses[r][c], grid[r][c] != word[index]:
		return 0
	}

	var used, offset int

	if uses[r][c] = true; word[index] == 'q' {
		offset++
	}

	for _, direction := range [...][2]int{
		{-1, -1}, {-1, +0}, {-1, +1},
		{+0, -1}, {+0, +1},
		{+1, -1}, {+1, +0}, {+1, +1},
	} {
		if n := dfs(grid, uses, word, index+1+offset, r+direction[0], c+direction[1]) + offset; n > used {
			used = n
		}
	}

	uses[r][c] = false

	return used + 1
}
