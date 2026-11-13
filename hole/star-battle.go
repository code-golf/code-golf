package hole

import (
	"math"
	"math/bits"
	"math/rand/v2"
	"strings"
)

// Modelled on https://github.com/code-golf/code-golf/issues/2781: a 9 × 9
// Star Battle board with 9 regions, printed in the style of Sudoku (Fill-in).
// The argument is the drawn board, the expected output is the same board with
// the 9 stars filled in. Regions are shown as digits 1-9; heavy lines mark
// borders between regions, dots separate cells within a region.
//
// Every board has exactly one solution: pick a random valid placement of 9
// non-touching stars, grow 9 connected regions outward from the stars, then
// while a second solution exists move one of its cells into a neighbouring
// region, keeping regions connected, until the solution is unique.

// All 47622 placements of 9 stars with one star per row, one per column and
// no two stars touching, even diagonally.
var starPlacements = func() [][9]int {
	var out [][9]int

	var rec func(p [9]int, row, used int)

	rec = func(p [9]int, row, used int) {
		if row == 9 {
			out = append(out, p)
			return
		}

		for col := 0; col < 9; col++ {
			if used&(1<<col) != 0 {
				continue
			}

			if row > 0 && (col == p[row-1] || col == p[row-1]+1 || col == p[row-1]-1) {
				continue
			}

			p[row] = col
			rec(p, row+1, used|1<<col)
		}
	}

	rec([9]int{}, 0, 0)

	return out
}()

// Heavy box drawing runes for the border arms meeting at a junction, packed as
// up<<3|down<<2|left<<1|right. Interior junctions have an even number of arms
// (region borders are closed loops); along the frame a third arm can join.
var starJunction = [16]rune{
	1:  '╺', // right
	2:  '╸', // left
	3:  '━', // left + right
	4:  '╻', // down
	5:  '┏', // down + right
	6:  '┓', // down + left
	7:  '┳', // down + left + right
	8:  '╹', // up
	9:  '┗', // up + right
	10: '┛', // up + left
	11: '┻', // up + left + right
	12: '┃', // up + down
	13: '┣', // up + down + right
	14: '┫', // up + down + left
	15: '╋', // all four
}

type starCell struct{ row, col int }

// Grow 9 connected regions outward from the 9 stars; each star cell starts
// its own region (1-9 in row order). Every empty cell adjacent to a region
// can grow into it, chosen with weighted randomness.
func starBattleGrow(sol [9]int) [9][9]byte {
	var g [9][9]byte
	var size [10]int

	for row := range 9 {
		g[row][sol[row]] = byte(row + 1)
		size[row+1] = 1
	}

	// Big regions may grow more or less aggressively than small ones.
	bias := 0.6 + rand.Float64()*1.4
	left := 81 - 9

	type frontierCell struct {
		cell   starCell
		region byte
		weight float64
	}

	var frontier []frontierCell

	for left > 0 {
		frontier = frontier[:0]
		total := 0.

		for row := range 9 {
			for col := range 9 {
				if g[row][col] != 0 {
					continue
				}

				for _, n := range [4][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
					r, c := row+n[0], col+n[1]
					if r < 0 || r > 8 || c < 0 || c > 8 || g[r][c] == 0 {
						continue
					}

					weight := math.Pow(float64(size[g[r][c]]), -bias)
					frontier = append(frontier, frontierCell{starCell{row, col}, g[r][c], weight})
					total += weight
				}
			}
		}

		pick, i := rand.Float64()*total, 0
		for pick > frontier[i].weight && i < len(frontier)-1 {
			pick -= frontier[i].weight
			i++
		}

		region, cell := frontier[i].region, frontier[i].cell
		g[cell.row][cell.col] = region
		size[region]++
		left--
	}

	return g
}

// Check that the given region, with the cell itself removed, is still
// orthogonally connected.
func starBattleConnected(g *[9][9]byte, region byte, skip starCell) bool {
	var start starCell
	var seen [81]bool
	total := 0

	for row := range 9 {
		for col := range 9 {
			if g[row][col] != region || (row == skip.row && col == skip.col) {
				continue
			}

			total++

			if total == 1 {
				start = starCell{row, col}
			}
		}
	}

	if total == 0 {
		return false
	}

	stack := []starCell{start}
	seen[start.row*9+start.col] = true
	reached := 0

	for len(stack) > 0 {
		cur := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		reached++

		for _, d := range [4][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
			r, c := cur.row+d[0], cur.col+d[1]
			if r < 0 || r > 8 || c < 0 || c > 8 || seen[r*9+c] ||
				g[r][c] != region || (r == skip.row && c == skip.col) {
				continue
			}

			seen[r*9+c] = true
			stack = append(stack, starCell{r, c})
		}
	}

	return reached == total
}

// Return up to limit solutions of the board, as star column per row.
func starBattleSolutions(g *[9][9]byte, limit int) [][9]int {
	var out [][9]int

	for _, p := range starPlacements {
		var regions uint

		for row := range 9 {
			regions |= 1 << (g[row][p[row]] - 1)
		}

		if bits.OnesCount(regions) == 9 {
			out = append(out, p)

			if len(out) == limit {
				break
			}
		}
	}

	return out
}

func starBattleBoard() (board [9][9]byte, sol [9]int) {
	for {
		p := starPlacements[rand.IntN(len(starPlacements))]
		g := starBattleGrow(p)

		for range 400 {
			solutions := starBattleSolutions(&g, 2)

			if len(solutions) == 1 {
				// Rename the regions so the digits aren't tied to rows.
				perm := shuffle([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9})

				for row := range 9 {
					for col := range 9 {
						board[row][col] = perm[g[row][col]-1] + '0'
					}

					sol[row] = solutions[0][row]
				}

				return
			}

			// Unique-solution fix: move one cell of the alternative
			// solution into a neighbouring region.
			var alt [9]int
			for _, s := range solutions {
				if s != p {
					alt = s
					break
				}
			}

			type starMove struct {
				row, col int
				region   byte
			}

			var moves []starMove

			for row := range 9 {
				col := alt[row]

				// Skip cells of the board's own solution and cells whose
				// region would fall apart without them.
				if col == p[row] || !starBattleConnected(&g, g[row][col], starCell{row, col}) {
					continue
				}

				// The distinct region values orthogonally adjacent.
				adjacent := uint(0)
				for _, n := range [4][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
					if r, c := row+n[0], col+n[1]; r >= 0 && r < 9 && c >= 0 && c < 9 {
						adjacent |= 1 << g[r][c]
					}
				}

				for r := byte(1); r <= 9; r++ {
					if adjacent&(1<<r) != 0 && r != g[row][col] {
						moves = append(moves, starMove{row, col, r})
					}
				}
			}

			if len(moves) == 0 {
				break
			}

			mov := randChoice(moves)
			g[mov.row][mov.col] = mov.region
		}
	}
}

// Print the board as a 9 × 9 grid. Heavy lines mark the outer frame and the
// borders between regions, while cells inside a region are separated by dots.
// Cells hold their region digit, or a star where stars[row] == col.
func starBattlePrint(g [9][9]byte, stars *[9]int) string {
	// vert[r][c] is the edge left of cell (r, c), horz[r][c] the edge above
	// it. Both are true along the outer frame.
	var vert, horz [10][10]bool

	for row := range 9 {
		for col := range 9 {
			vert[row][col] = col == 0 || g[row][col-1] != g[row][col]
			horz[row][col] = row == 0 || g[row-1][col] != g[row][col]
		}

		vert[row][9], horz[9][row] = true, true
	}

	var b strings.Builder
	b.Grow(19 * 112)

	for row := range 10 {
		for col := range 10 {
			arms := 0

			if row > 0 && vert[row-1][col] {
				arms |= 8
			}

			if row < 9 && vert[row][col] {
				arms |= 4
			}

			if col > 0 && horz[row][col-1] {
				arms |= 2
			}

			if col < 9 && horz[row][col] {
				arms |= 1
			}

			if arms == 0 {
				b.WriteRune('∙')
			} else {
				b.WriteRune(starJunction[arms])
			}

			if col < 9 {
				if horz[row][col] {
					b.WriteString("━━━")
				} else {
					b.WriteString("∙∙∙")
				}
			}
		}

		b.WriteByte('\n')

		if row < 9 {
			for col := range 9 {
				if vert[row][col] {
					b.WriteRune('┃')
				} else {
					b.WriteRune('∙')
				}

				b.WriteByte(' ')

				if stars != nil && stars[row] == col {
					b.WriteByte('*')
				} else {
					b.WriteByte(g[row][col])
				}

				b.WriteString(" ")
			}

			b.WriteString("┃\n")
		}
	}

	return b.String()[:b.Len()-1]
}

var _ = answerFunc("star-battle", func() []Answer {
	answers := make([]Answer, 3)

	for i := range answers {
		board, stars := starBattleBoard()

		answers[i] = Answer{
			Args:   []string{starBattlePrint(board, nil)},
			Answer: starBattlePrint(board, &stars),
		}
	}

	return answers
})
