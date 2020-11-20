package hole

import (
	"math/rand"
	"strings"
)

const north = 1
const south = 2
const west = 4
const east = 8

const width = 50
const height = 50

var dj = map[int]int{east: 1, west: -1, north: 0, south: 0}
var di = map[int]int{east: 0, west: 0, north: -1, south: 1}
var opposite = map[int]int{east: west, west: east, north: south, south: north}

// shorturl.at/mQT19
func dig(i, j int, grid, dist [height][width]int) ([height][width]int, [height][width]int) {
	directions := []int{north, south, west, east}
	rand.Shuffle(len(directions), func(m, n int) {
		directions[m], directions[n] = directions[n], directions[m]
	})

	for _, d := range directions {
		new_i := i + di[d]
		new_j := j + dj[d]

		if new_j >= 0 && new_j < width && new_i >= 0 && new_i < height && grid[new_i][new_j] == 0 {
			grid[i][j] |= d
			dist[new_i][new_j] = dist[i][j] + 1
			grid[new_i][new_j] |= opposite[d]
			grid, dist = dig(new_i, new_j, grid, dist)
		}
	}
	return grid, dist
}

func find_exit(dist [height][width]int) (ei, ej int) {
	maxd := -1
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			if dist[i][j] > maxd {
				maxd = dist[i][j]
				ei, ej = i, j
			}
		}
	}
	return ei, ej
}

func trace_path(dist [height][width]int, ei, ej int) (path [height][width]int) {
	directions := []int{north, south, west, east}
	d := dist[ei][ej]
	path[ei][ej] = 1
	i, j := ei, ej
	for d > 0 {
		for _, dir := range directions {
			i_, j_ := i+di[dir], j+dj[dir]
			if i_ >= 0 && i_ < height && j_ >= 0 && j_ < width {
				if dist[i_][j_] == d-1 {
					d -= 1
					path[i_][j_] = 1
					i, j = i_, j_
					break
				}
			}
		}
	}
	return
}

func draw(grid [height][width]int, si, sj, ei, ej int,
	path [height][width]int, draw_path bool) (maze_s string) {

	wall, track := "█", ""
	if draw_path {
		track = "."
	} else {
		track = " "
	}
	top, bottom, cell := "", "", ""
	east_boundary, south_boundary := "", ""

	maze_s = wall + strings.Repeat(strings.Repeat(wall, 2), width) + "\n"
	for i := 0; i < height; i++ {
		top = wall
		bottom = wall
		for j := 0; j < width; j++ {
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
					east_boundary = track
				} else {
					east_boundary = " "
				}
			} else {
				east_boundary = wall
			}
			if grid[i][j]&south != 0 {
				if path[i+1][j] != 0 && path[i][j] != 0 {
					south_boundary = track
				} else {
					south_boundary = " "
				}
			} else {
				south_boundary = wall
			}
			top += cell + east_boundary
			bottom += south_boundary + wall
		}
		maze_s += top + "\n" + bottom + "\n"
	}
	return
}

func maze() (args []string, out string) {

	var grid [height][width]int
	var dist [height][width]int

	sj := rand.Intn(width)
	si := rand.Intn(height)

	grid, dist = dig(si, sj, grid, dist)
	ei, ej := find_exit(dist)
	path := trace_path(dist, ei, ej)
	maze_input := draw(grid, si, sj, ei, ej, path, false)
	maze_solved := draw(grid, si, sj, ei, ej, path, true)

	args = append(args, maze_input)
	out = maze_solved
	return
}
