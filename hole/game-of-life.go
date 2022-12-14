package hole

import (
	"math/rand"
	"strings"
)

const lifeSize = 32

var lifeTemplate = [lifeSize]string{
	"................................",
	"..##.#####..###......#.##??????.",
	".#......#.###..#.###..#.#??????.",
	".....#..######.#..##.##.#??????.",
	".#......#.###..#.###..#.#??????.",
	".####.##......#.##..#..#.??????.",
	".#......#.###..#.###..#.#??????.",
	".....#..######.#..##.##.#??????.",
	"....##...###.#.##.#####..??????.",
	".##...#.#..##.##.#.#....#?????#.",
	"..######.#...##.#...##.#.?????#.",
	".####.##......#.##..#..#.?????#.",
	".##.#.#....#..####.##....??????.",
	".#.#...##.#.#....##...###??????.",
	"..######.#...##.#...##.#.??????.",
	".####.##......#.##..#..#.??????.",
	".#.##..#..#.....###.#.##.??????.",
	"..##.#####..###......#.##??????.",
	".#......#.###..#.###..#.#??????.",
	".##..#..#.#.#.##...######??????.",
	"..###..#.###.#...##.#...#??????.",
	"..#.#....##..#.#.####..##??????.",
	"..###..#.###...#..####.##??????.",
	".#.#.####..#######.#..##.??????.",
	"..###..#.###.#...##.#...#??????.",
	"..#.#....##..#.#.####..##??????.",
	"....#.#..#...#.#####...#.??????.",
	"..##...########..#..#.#.#??????.",
	".#...##.#...#.###..#.###.??????.",
	".??????????????????????????????.",
	".??????????????????????###?????.",
	"................................",
}

type grid []string

func randGrid() grid {
	grid := make(grid, lifeSize)
	var array [lifeSize]byte
	for i := 0; i < lifeSize; i++ {
		for j := 0; j < lifeSize; j++ {
			c := lifeTemplate[i][j]
			if c == '?' {
				c = ".#"[rand.Int31()&1]
			}
			array[j] = c
		}
		grid[i] = string(array[:])
	}
	return grid
}

func (grid grid) get(i, j int) int {
	if i < 0 || i >= lifeSize || j < 0 || j >= lifeSize || grid[i][j] == '.' {
		return 0
	}
	return 1
}

func (g grid) step() grid {
	next := make(grid, lifeSize)
	var array [lifeSize]byte
	for i := 0; i < lifeSize; i++ {
		for j := 0; j < lifeSize; j++ {
			neighbours := g.get(i+1, j+1) + g.get(i+1, j) + g.get(i+1, j-1) +
				g.get(i, j+1) + g.get(i, j-1) +
				g.get(i-1, j+1) + g.get(i-1, j) + g.get(i-1, j-1)
			if neighbours == 3 || neighbours == 2 && g.get(i, j) == 1 {
				array[j] = '#'
			} else {
				array[j] = '.'
			}
		}
		next[i] = string(array[:])
	}
	return next
}

func (grid grid) toString() string {
	return strings.Join(grid, "\n")
}

func gameOfLife() []Scorecard {
	grid := randGrid()
	return []Scorecard{{
		Args:   []string{grid.toString()},
		Answer: grid.step().toString(),
	}}
}
