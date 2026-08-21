package hole

import (
	"fmt"
	"math/rand/v2"
)

var hardCodedRules = []string{
	"LLLLR", "LRLRLLRR", "RRRRL", "LLLRRR", "RRRLLL",
}

func generateRuleset(length int) string {
	ruleset := ""
	for range length {
		ruleset += randChoice([]string{"L", "R"})
	}
	return ruleset
}

func runAnt(ruleset string) string {
	rulelength := len(ruleset)

	var grid [33][33]int

	posX, posY, velX, velY := 16, 16, 0, -1

	for range 1000 {
		dir := ruleset[grid[posY][posX]]
		if dir == 'L' {
			velX, velY = velY, -velX
		} else {
			velX, velY = -velY, velX
		}
		grid[posY][posX] = (grid[posY][posX] + 1) % rulelength
		posX += velX
		posY += velY

		if posX < 0 || posX > 32 || posY < 0 || posY > 32 {
			break
		}
	}

	outp := ""

	for _, row := range grid {
		if outp != "" {
			outp += "\n"
		}
		for _, cell := range row {
			outp += fmt.Sprint(cell)
		}
	}

	return outp + "\n"
}

var _ = answerFunc("langtons-ant", func() []Answer {
	tests := make([]test, 40)
	for i, rule := range hardCodedRules {
		grid := runAnt(rule)
		tests[i] = test{rule, grid}
	}
	for i := range 9 {
		ruleset := generateRuleset(i + 2)
		grid := runAnt(ruleset)
		tests[i+5] = test{ruleset, grid}
	}
	for i := len(hardCodedRules) + 9; i < 40; i++ {
		ruleset := generateRuleset(rand.IntN(8) + 2)
		grid := runAnt(ruleset)
		tests[i] = test{ruleset, grid}
	}
	shuffle(tests)
	return outputTests(tests[:20], tests[20:])
})
