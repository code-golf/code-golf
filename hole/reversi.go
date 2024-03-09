package hole

import (
	"math/rand/v2"
	"strings"
)

const (
	reversiGridSize = 8
)

type ReversiTile int64

const (
	Empty ReversiTile = iota
	Black
	White
	PotentialSpot
)

func otherTeam(team ReversiTile) ReversiTile {
	if team == Black {
		return White
	} else if team == White {
		return Black
	} else {
		return team
	}
}

func inRange(x, y int) bool {
	return x >= 0 && y >= 0 && x < reversiGridSize && y < reversiGridSize
}

type ReversiSpot struct {
	Pos   [2]int
	Tiles [][2]int
}

func getPotentialSpots(team ReversiTile, grid [reversiGridSize][reversiGridSize]ReversiTile) []ReversiSpot {
	out := []ReversiSpot{}

	directions := [8][2]int{
		{-1, -1},
		{-1, 0},
		{-1, 1},
		{0, 1},
		{1, 1},
		{1, 0},
		{1, -1},
		{0, -1},
	}

	for i := range reversiGridSize {
		for j := range reversiGridSize {
			flippableTiles := [][2]int{}

			if grid[i][j] != Empty {
				continue
			}

			for _, direction := range directions {
				nextTileI := i + direction[0]
				nextTileJ := j + direction[1]

				if !inRange(nextTileJ, nextTileI) {
					continue
				}

				if grid[nextTileI][nextTileJ] != otherTeam(team) {
					continue
				}

				flippableTilesInDirection := [][2]int{}

				for inRange(nextTileJ, nextTileI) && grid[nextTileI][nextTileJ] == otherTeam(team) {
					flippableTilesInDirection = append(flippableTilesInDirection, [2]int{nextTileI, nextTileJ})

					nextTileI += direction[0]
					nextTileJ += direction[1]
				}

				if inRange(nextTileJ, nextTileI) && grid[nextTileI][nextTileJ] == team {
					flippableTiles = append(flippableTiles, flippableTilesInDirection...)
				}
			}

			if len(flippableTiles) > 0 {
				out = append(out, ReversiSpot{
					Pos:   [2]int{i, j},
					Tiles: flippableTiles,
				})
			}
		}
	}

	return out
}

func getReversiInitialState() [reversiGridSize][reversiGridSize]ReversiTile {
	var out [reversiGridSize][reversiGridSize]ReversiTile

	out[3][4] = Black
	out[4][3] = Black
	out[3][3] = White
	out[4][4] = White
	return out
}

func genReversiBoard(steps int) [reversiGridSize][reversiGridSize]ReversiTile {
	out := getReversiInitialState()

	for i := range steps {
		team := []ReversiTile{Black, White}[i%2]

		spots := getPotentialSpots(team, out)

		spot := spots[rand.IntN(len(spots))]

		out[spot.Pos[0]][spot.Pos[1]] = team
		for _, reversedSpot := range spot.Tiles {
			out[reversedSpot[0]][reversedSpot[1]] = team
		}

	}

	return out
}

func drawReversiBoard(board [reversiGridSize][reversiGridSize]ReversiTile) string {
	const blackChar string = "x"
	const whiteChar string = "o"
	const emptyChar string = "."
	const potentialSpotChar string = "!"

	reversiString := ""

	for i := 0; i < reversiGridSize; i++ {
		for j := 0; j < reversiGridSize; j++ {
			if board[i][j] == Empty {
				reversiString += emptyChar
			} else if board[i][j] == Black {
				reversiString += blackChar
			} else if board[i][j] == White {
				reversiString += whiteChar
			} else if board[i][j] == PotentialSpot {
				reversiString += potentialSpotChar
			}
		}

		if reversiGridSize-1 != i {
			reversiString += "\n"
		}
	}

	return reversiString
}

func highlightCorrectAnswersReversiBoard(board [reversiGridSize][reversiGridSize]ReversiTile) [reversiGridSize][reversiGridSize]ReversiTile {
	for _, spot := range getPotentialSpots(White, board) {
		board[spot.Pos[0]][spot.Pos[1]] = PotentialSpot
	}
	return board
}

func reversi() []Run {
	const runs = 20

	args := []string{}
	answer := []string{}

	for run := range runs {
		grid := genReversiBoard((run+1)/2*2 + 1)

		args = append(args, drawReversiBoard(grid))
		answer = append(answer, drawReversiBoard(highlightCorrectAnswersReversiBoard(grid)))
	}

	return []Run{
		{
			Args:   args,
			Answer: strings.Join(answer, "\n\n"),
		},
	}
}
