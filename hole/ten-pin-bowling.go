package hole

import (
	"math/rand"
	"strconv"
	"strings"
)

var games = [...]struct {
	frames []rune
	score  string
}{
	{[]rune(" X  X  X  X  X  X  X  X  X XXX"), "300"},
	{[]rune(" X 17 36 63 4-  X 61 7- 6- -- "), "85"},
	{[]rune(" X 7/ 9-  X -8 8/ -6  X  X X81"), "167"},
	{[]rune(" X 7/ 9-  X -8 8/ -6  X  X X8/"), "168"},
	{[]rune(" X 8- 51  X 35 36 7- 9-  X 8- "), "109"},
	{[]rune("-- -- -- -- -- -- -- -- -- -- "), "0"},
	{[]rune("-9 5- 35 31 61 43 6- 63 6- 71 "), "69"},
	{[]rune("32 3/  X  X  X  X 43 33 33 3/6"), "161"},
	{[]rune("32 3/  X  X  X  X 43 33 33 36 "), "154"},
	{[]rune("43 44 54 45  X  X  X  X 43 23 "), "146"},
	{[]rune("53 33 34  X  X  X 53 3/  X X43"), "163"},
	{[]rune("7/ 4- 36 81 8- 54 44 53 31 8- "), "81"},
	{[]rune("71 33 45 45  X  X  X  X 5/ 23 "), "154"},
	{[]rune("71 7- 72 8- 81 51 8-  X 6- 81 "), "86"},
	{[]rune("72 9- 81  X 9- 8/  X  X  X 9- "), "162"},
	{[]rune("81 16 8/ 33  X 7- -7 9- 8- -- "), "83"},
	{[]rune("9- -2 35  X  X  X  X 62 22 62 "), "143"},
	{[]rune("9/ 5F 5- F2  X 5- 81  X F7 9/3"), "93"},
}

// Randomly create splits and fouls to make it more interesting.
func randReplacements(gFrames []rune) []rune {
	frames := make([]rune, len(gFrames))
	copy(frames, gFrames)

	for j, char := range frames {
		var replacement rune

		switch char {
		case '0':
			frames[j] = '-'
			replacement = 'F'
		case '-':
			replacement = 'F'
		case '5', '6', '7', '8':
			// Only split on the first ball of the frame.
			if j%3 == 0 {
				replacement = '①' - '1' + char
			} else {
				continue
			}
		default:
			continue
		}

		if rand.Intn(2) == 0 {
			frames[j] = replacement
		}
	}
	return frames
}

func tenPinBowling() ([]string, string) {
	extraCases := 22
	args := make([]string, len(games)+extraCases)
	outs := make([]string, len(games)+extraCases)

	for i, game := range games {
		frames := randReplacements(game.frames)
		args[i] = string(frames)
		outs[i] = game.score
	}

	for i := 0; i < extraCases; i++ {
		rolls := make([]int, 24)

		// Generate some random rolls
		for rollNum := 0; rollNum < 23; rollNum++ {
			roll := 0
			maxRoll := 10
			if rollNum%2 == 1 {
				maxRoll = 10 - rolls[rollNum-1]
			}
			// Roll with bias towards strikes/spares and misses
			roll = rand.Intn(maxRoll+2) - 1
			if roll > maxRoll {
				roll = maxRoll
			} else if roll < 0 {
				roll = 0
			}
			rolls[rollNum] = roll
		}

		if rolls[18] != 10 {
			rolls[21] = 0
		}

		// Now let's score the rolls and format the input
		arg := ""
		score := 0
		for frame := 0; frame < 12; frame++ {
			// Skip the bonus rolls if necessary
			if frame > 9 {
				if rolls[18]+rolls[19] != 10 || (frame == 11 && (rolls[18] != 10 || rolls[20] != 10)) {
					arg += " "
					break
				}
			} else if frame > 0 {
				arg += " "
				// Add the bonus scores if the last frame was a spare or strike
				if rolls[frame*2-2]+rolls[frame*2-1] == 10 {
					score += rolls[frame*2]

					// If it was a strike specifically
					if rolls[frame*2-2] == 10 {
						score += rolls[frame*2+1]
						// And if the one before that was also a strike
						if frame > 1 && rolls[frame*2-4] == 10 {
							score += rolls[frame*2]
						}
					}
				}
			}
			score += rolls[frame*2] + rolls[frame*2+1]

			if rolls[frame*2] == 10 {
				if frame < 9 {
					arg += " "
				}
				arg += "X"
			} else {
				arg += string('0' + rolls[frame*2])
				if frame < 10 || (frame == 10 && rolls[18] == 10) {
					if rolls[frame*2]+rolls[frame*2+1] == 10 {
						arg += "/"
					} else {
						arg += string('0' + rolls[frame*2+1])
					}
				}
			}
		}
		r := []rune(arg)
		frames := randReplacements(r)
		args[i+len(games)] = string(frames)
		outs[i+len(games)] = strconv.Itoa(score)
	}

	rand.Shuffle(len(args), func(i, j int) {
		args[i], args[j] = args[j], args[i]
		outs[i], outs[j] = outs[j], outs[i]
	})

	return args, strings.Join(outs, "\n")
}
