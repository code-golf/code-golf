package hole

import (
	"math/rand"
	"strconv"
	"strings"
)

func turtle() (args []string, out string) {
	// set argVec might add the bigger numbers if you want
	moveAllVec := []string{
		"0 N", "0 E", "0 S", "0 W", "1 N", "1 E", "1 S", "1 W", "2 N", "2 E", "2 S", "2 W", "3 N", "3 E", "3 S", "3 W", "4 N", "4 E", "4 S", "4 W", "5 N", "5 E", "5 S", "5 W", "6 N", "6 E", "6 S", "6 W", "7 N", "7 E", "7 S", "7 W", "8 N", "8 E", "8 S", "8 W", "9 N", "9 E", "9 S", "9 W", "10 N", "10 E", "10 S", "10 W",
		"11 N", "11 E", "11 S", "11 W", "12 N", "12 E", "12 S", "12 W", "13 N", "13 E", "13 S", "13 W", "14 N", "14 E", "14 S", "14 W", "15 N", "15 E", "15 S", "15 W", "16 N", "16 E", "16 S", "16 W", "17 N", "17 E", "17 S", "17 W", "18 N", "18 E", "18 S", "18 W", "19 N", "19 E", "19 S", "19 W", "20 N", "20 E", "20 S", "20 W",
		// "21 N", "21 E","21 S", "21 W", "22 N", "22 E", "22 S", "22 W", "23 N", "23 E", "23 S","23 W", "24 N", "24 E", "24 S", "24 W", "25 N", "25 E", "25 S", "25 W","26 N", "26 E", "26 S", "26 W", "27 N", "27 E", "27 S", "27 W", "28 N","28 E", "28 S", "28 W", "29 N", "29 E", "29 S", "29 W", "30 N", "30 E","30 S", "30 W",
		/* "31 N", "31 E", "31 S", "31 W", "32 N", "32 E", "32 S","32 W", "33 N", "33 E", "33 S", "33 W", "34 N", "34 E", "34 S", "34 W","35 N", "35 E", "35 S", "35 W", "36 N", "36 E", "36 S", "36 W", "37 N","37 E", "37 S", "37 W", "38 N", "38 E", "38 S", "38 W", "39 N", "39 E","39 S", "39 W", "40 N", "40 E", "40 S", "40 W",
		   "41 N", "41 E", "41 S","41 W", "42 N", "42 E", "42 S", "42 W", "43 N", "43 E", "43 S", "43 W","44 N", "44 E", "44 S", "44 W", "45 N", "45 E", "45 S", "45 W", "46 N","46 E", "46 S", "46 W", "47 N", "47 E", "47 S", "47 W", "48 N", "48 E","48 S", "48 W", "49 N", "49 E", "49 S", "49 W", "50 N", "50 E", "50 S","50 W",
		   "51 N", "51 E", "51 S", "51 W", "52 N", "52 E", "52 S", "52 W","53 N", "53 E", "53 S", "53 W", "54 N", "54 E", "54 S", "54 W", "55 N","55 E", "55 S", "55 W", "56 N", "56 E", "56 S", "56 W", "57 N", "57 E","57 S", "57 W", "58 N", "58 E", "58 S", "58 W", "59 N", "59 E", "59 S","59 W", "60 N", "60 E", "60 S", "60 W",
		   "61 N", "61 E", "61 S", "61 W","62 N", "62 E", "62 S", "62 W", "63 N", "63 E", "63 S", "63 W", "64 N","64 E", "64 S", "64 W", "65 N", "65 E", "65 S", "65 W", "66 N", "66 E","66 S", "66 W", "67 N", "67 E", "67 S", "67 W", "68 N", "68 E", "68 S","68 W", "69 N", "69 E", "69 S", "69 W", "70 N", "70 E", "70 S", "70 W",
		   "71 N", "71 E", "71 S", "71 W", "72 N", "72 E", "72 S", "72 W", "73 N","73 E", "73 S", "73 W", "74 N", "74 E", "74 S", "74 W", "75 N", "75 E","75 S", "75 W", "76 N", "76 E", "76 S", "76 W", "77 N", "77 E", "77 S","77 W", "78 N", "78 E", "78 S", "78 W", "79 N", "79 E", "79 S", "79 W","80 N", "80 E", "80 S", "80 W",
		   "81 N", "81 E", "81 S", "81 W", "82 N","82 E", "82 S", "82 W", "83 N", "83 E", "83 S", "83 W", "84 N", "84 E","84 S", "84 W", "85 N", "85 E", "85 S", "85 W", "86 N", "86 E", "86 S","86 W", "87 N", "87 E", "87 S", "87 W", "88 N", "88 E", "88 S", "88 W","89 N", "89 E", "89 S", "89 W", "90 N", "90 E", "90 S", "90 W",
		   "91 N","91 E", "91 S", "91 W", "92 N", "92 E", "92 S", "92 W", "93 N", "93 E","93 S", "93 W", "94 N", "94 E", "94 S", "94 W", "95 N", "95 E", "95 S","95 W", "96 N", "96 E", "96 S", "96 W", "97 N", "97 E", "97 S", "97 W","98 N", "98 E", "98 S", "98 W", "99 N", "99 E", "99 S", "99 W", "100 N", "100 E", "100 S", "100 W", */
	}

	rand.Shuffle(len(moveAllVec), func(i, j int) {
		moveAllVec[i], moveAllVec[j] = moveAllVec[j], moveAllVec[i]
	})

	var argVec []string

	for i := 0; i < 69; i++ {
		argVec = append(argVec, moveAllVec[i])
	}

	// set out
	var currI, currJ, minI, minJ, maxI, maxJ, startI, startJ, endI, endJ int
	posHash := make(map[[2]int]int)

	for _, arg := range argVec {

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

	args = append(args, strings.Join(argVec[:], "\n"))
	out = strings.Join(outVec[:], "\n")

	return
}
