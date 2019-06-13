package routes

import (
	"math/rand"
	"strings"
)

var games = [...]struct {
	frames []rune
	score  string
}{
	{[]rune(" X  X  X  X  X  X  X  X  X XXX"), "300"},
	{[]rune(" X 17 36 63 4-  X 61 7- 6- -- "), "85"},
	{[]rune(" X 7/ 9-  X -8 8/ -6  X  X X81"), "167"},
	{[]rune(" X 8- 51  X 35 36 7- 9-  X 8- "), "109"},
	{[]rune("-- -- -- -- -- -- -- -- -- -- "), "0"},
	{[]rune("-9 5- 35 31 61 43 6- 63 6- 71 "), "69"},
	{[]rune("32 3/  X  X  X  X 43 33 33 36 "), "154"},
	{[]rune("43 44 54 45  X  X  X  X 43 23 "), "146"},
	{[]rune("53 33 34  X  X  X 53 3/  X X43"), "163"},
	{[]rune("7/ 4- 36 81 8- 54 44 53 31 8- "), "81"},
	{[]rune("71 33 45 45  X  X  X  X 5/ 23 "), "154"},
	{[]rune("71 7- 72 8- 81 51 8-  X 6- 81 "), "86"},
	{[]rune("72 9- 81  X 9- 8/  X  X  X 9- "), "162"},
	{[]rune("81 16 8/ 33  X 7- -7 9- 8- -- "), "83"},
	{[]rune("9- -2 35  X  X  X  X 62 22 62 "), "143"},
    {[]rune("32 3/  X  X  X  X 43 33 33 3/6"), "161"},
}

func tenPinBowling() ([]string, string) {
	args := make([]string, len(games))
	outs := make([]string, len(games))

	for i, game := range games {
		frames := make([]rune, len(game.frames))
		copy(frames, game.frames)

		// Randomly create splits and fouls to make it more interesting.
		for j, char := range frames {
			var replacement rune

			switch char {
			case '-':
				replacement = 'F'
			case '5', '6', '7', '8':
				if (j%3==0) { // Only split on the first ball of the frame
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

		args[i] = string(frames)
		outs[i] = game.score
	}

	// Shuffle
	for i := range args {
		j := rand.Intn(i + 1)
		args[i], args[j] = args[j], args[i]
		outs[i], outs[j] = outs[j], outs[i]
	}

	return args, strings.Join(outs, "\n")
}
