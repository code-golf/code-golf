package hole

import (
	"strconv"
	"strings"
)

func turtle() []Run {
	var argVec []string
	for i := range 21 {
		for _, dir := range []string{"N", "E", "S", "W"} {
			argVec = append(argVec, strconv.Itoa(i)+" "+dir)
		}
	}

	// set out
	var currI, currJ, minI, minJ, maxI, maxJ, startI, startJ, endI, endJ int
	posHash := make(map[[2]int]int)

	for _, arg := range shuffle(argVec) {
		moveVec := strings.Split(arg, " ")
		move, _ := strconv.Atoi(moveVec[0])
		dir := moveVec[1]

		if dir == "N" {
			for i := currI; i >= currI-move; i-- {
				posHash[[2]int{i, currJ}] = 1
			}

			currI -= move
		} else if dir == "S" {
			for i := currI; i <= currI+move; i++ {
				posHash[[2]int{i, currJ}] = 1
			}

			currI += move
		} else if dir == "E" {
			for j := currJ; j <= currJ+move; j++ {
				posHash[[2]int{currI, j}] = 1
			}

			currJ += move
		} else if dir == "W" {
			for j := currJ; j >= currJ-move; j-- {
				posHash[[2]int{currI, j}] = 1
			}

			currJ -= move
		}

		minI = min(currI, minI)
		minJ = min(currJ, minJ)
	}

	posHashNew := make(map[[2]int]int)

	for pos := range posHash {
		i := pos[0]
		j := pos[1]

		iNew := max(i, i-minI)
		jNew := max(j, j-minJ)

		if i == 0 && j == 0 {
			startI = iNew
			startJ = jNew
		}

		if i == currI && j == currJ {
			endI = iNew
			endJ = jNew
		}

		maxI = max(iNew, maxI)
		maxJ = max(jNew, maxJ)

		posHashNew[[2]int{iNew, jNew}] = 1
	}

	var outVec []string

	for i := 0; i <= maxI; i++ {
		line := ""
		for j := 0; j <= maxJ; j++ {
			if _, ok := posHashNew[[2]int{i, j}]; ok {
				if i == endI && j == endJ {
					line += "🐢"
				} else if i == startI && j == startJ {
					line += "🏁"
				} else {
					line += "⬜"
				}
			} else {
				line += "🟩"
			}
		}

		outVec = append(outVec, line)
	}

	return []Run{{
		Args:   []string{strings.Join(argVec[:], "\n")},
		Answer: strings.Join(outVec[:], "\n"),
	}}
}
