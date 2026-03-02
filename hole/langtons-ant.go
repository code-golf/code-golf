package hole

import (
	"math/rand/v2"
	"fmt"
)

func generateRuleset(length int) string {
  ruleset := ""
  for range length {
    ruleset += []string{"L", "R"}[rand.IntN(2)]
  }
  return ruleset
}


func computeGrid(ruleset string) string {
	var directions = [][]int{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}
	
	numStates := len(ruleset)
	outp := ""

	var g [33][33]int

	xPos, yPos, d := 16, 16, 0

	for x := range 33 {
		for y := range 33 {
			g[y][x] = 0
		}
	}

	for range 1000 {
		switch ruleset[g[yPos][xPos]] {
		case 'L':
			d = (d + 3) % 4
		case 'R':
			d = (d + 1) % 4
		}
		g[yPos][xPos] = (g[yPos][xPos] + 1) % numStates
		xPos += directions[d][0]
		yPos += directions[d][1]
	}

	for x := range 33 {
		for y := range 33 {
			outp += fmt.Sprint(g[y][x])
		}
		outp += "\n"
	}

	return outp

}

var _ = answerFunc("langtons-ant", func() []Answer {
	tests := make([]test, 25)

	for i := range tests {
		ruleset := generateRuleset(rand.IntN(8)+2)

		tests[i] = test{ruleset, computeGrid(ruleset)}
	}

	return outputTests(shuffle(tests))
})

