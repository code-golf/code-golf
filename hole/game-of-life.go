package hole

import "strings"

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
	for i := range lifeSize {
		for j := range lifeSize {
			c := lifeTemplate[i][j]
			if c == '?' {
				c = randChoice([]byte(".#"))
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
	for i := range lifeSize {
		for j := range lifeSize {
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

func gameOfLife() []Run {
	grid := randGrid()
	return []Run{{
		Args:   []string{grid.toString()},
		Answer: grid.step().toString(),
	}}
}
