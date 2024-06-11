package hole

import (
	"math/rand/v2"
	"strings"
)

const (
	north  = 1
	south  = 2
	west   = 4
	east   = 8
	width  = 25
	height = 25
)

var (
	dj       = map[int]int{east: 1, west: -1, north: 0, south: 0}
	di       = map[int]int{east: 0, west: 0, north: -1, south: 1}
	opposite = map[int]int{east: west, west: east, north: south, south: north}
)

// http://weblog.jamisbuck.org/2010/12/27/maze-generation-recursive-backtracking
func dig(i, j int, grid, dist [height][width]int) ([height][width]int, [height][width]int) {
	directions := shuffle([]int{north, south, west, east})

	for _, d := range directions {
		newi := i + di[d]
		newj := j + dj[d]

		if newj >= 0 && newj < width && newi >= 0 && newi < height && grid[newi][newj] == 0 {
			grid[i][j] |= d
			dist[newi][newj] = dist[i][j] + 1
			grid[newi][newj] |= opposite[d]
			grid, dist = dig(newi, newj, grid, dist)
		}
	}
	return grid, dist
}

func findExit(dist [height][width]int) (ei, ej int) {
	maxd := -1
	for i := range height {
		for j := range width {
			if dist[i][j] > maxd {
				maxd = dist[i][j]
				ei, ej = i, j
				// distance roulette
				if rand.Float64() < 0.1 {
					return
				}
			}
		}
	}
	return
}

func tracePath(dist [height][width]int, ei, ej int) (path [height][width]int) {
	directions := []int{north, south, west, east}
	d := dist[ei][ej]
	path[ei][ej] = 1
	i, j := ei, ej
	for d > 0 {
		for _, dir := range directions {
			newi, newj := i+di[dir], j+dj[dir]
			if newi >= 0 && newi < height && newj >= 0 && newj < width {
				if dist[newi][newj] == d-1 {
					d--
					path[newi][newj] = 1
					i, j = newi, newj
					break
				}
			}
		}
	}
	return
}

func draw(grid [height][width]int, si, sj, ei, ej int, path [height][width]int, drawpath bool) string {
	const wall = "#"

	var track, top, bottom, cell, eastboundary, southboundary string

	if drawpath {
		track = "."
	} else {
		track = " "
	}

	mazestr := strings.Repeat(wall, 2*width+1) + "\n"

	for i := range height {
		top = wall
		bottom = wall
		for j := range width {
			if i == si && j == sj {
				cell = "S"
			} else if i == ei && j == ej {
				cell = "E"
			} else {
				if path[i][j] == 0 {
					cell = " "
				} else {
					cell = track
				}
			}
			if grid[i][j]&east != 0 {
				if path[i][j+1] != 0 && path[i][j] != 0 {
					eastboundary = track
				} else {
					eastboundary = " "
				}
			} else {
				eastboundary = wall
			}
			if grid[i][j]&south != 0 {
				if path[i+1][j] != 0 && path[i][j] != 0 {
					southboundary = track
				} else {
					southboundary = " "
				}
			} else {
				southboundary = wall
			}
			top += cell + eastboundary
			bottom += southboundary + wall
		}
		mazestr += top + "\n" + bottom + "\n"
	}

	// Trim trailing newline.
	return mazestr[:len(mazestr)-1]
}

func maze() []Run {
	runs := make([][]test, 5)

	for i := range runs {
		var grid, dist [height][width]int

		sj := rand.IntN(width)
		si := rand.IntN(height)

		grid, dist = dig(si, sj, grid, dist)
		ei, ej := findExit(dist)
		path := tracePath(dist, ei, ej)

		runs[i] = []test{{
			in:  draw(grid, si, sj, ei, ej, path, false),
			out: draw(grid, si, sj, ei, ej, path, true),
		}}
	}

	return outputTests(runs...)
}
