package hole

import (
	"math/bits"
	"math/rand"
	"strings"
)

// de Bruijn torus (16,32;3,3) - contains every 3×3 matrix exactly once
// source: https://en.wikipedia.org/wiki/File:De_bruijn_torus_3x3.stl
var deBruijnTorus = [16]uint32{
	0b00010011111101001101101010000110,
	0b00000010111001011100101110010111,
	0b00000010111001011100101110010111,
	0b11101100000010110010010101111001,
	0b00000010111001011100101110010111,
	0b00010011111101001101101010000110,
	0b00110001110101101111100010100100,
	0b10001010011011010100001100011111,
	0b11111101000110100011010001101000,
	0b11101100000010110010010101111001,
	0b10101000010011110110000100111101,
	0b01000110101000011000111111010011,
	0b11111101000110100011010001101000,
	0b11101100000010110010010101111001,
	0b01100100100000111010110111110001,
	0b11011111001110000001011001001010,
}

// 32x32
type grid []uint32

func randGrid() grid {
	grid := make(grid, 32)
	torusShift := rand.Intn(16)
	dx := rand.Intn(32)
	dy := rand.Intn(32)
	for i := 0; i < 18; i++ {
		grid[(i+dx)&31] = bits.RotateLeft32(deBruijnTorus[(i+torusShift)&15], dy)
	}
	for i := 18; i < 32; i++ {
		grid[(i+dx)&31] = rand.Uint32()
	}
	return grid
}

func (grid grid) get(i, j int) uint32 {
	return (grid[i&31] >> (j & 31)) & 1
}

func (grid grid) set(i, j int) {
	grid[i&31] |= 1 << (j & 31)
}

func (g grid) step() grid {
	next := make(grid, 32)
	for i := 0; i < 32; i++ {
		for j := 0; j < 32; j++ {
			neighbours := g.get(i+1, j+1) + g.get(i+1, j) + g.get(i+1, j-1) +
				g.get(i, j+1) + g.get(i, j-1) +
				g.get(i-1, j+1) + g.get(i-1, j) + g.get(i-1, j-1)
			if neighbours == 3 || neighbours == 2 && g.get(i, j) == 1 {
				next.set(i, j)
			}
		}
	}
	return next
}

func (grid grid) toString() string {
	var buf strings.Builder
	for i := 0; i < 32; i++ {
		if i > 0 {
			buf.WriteByte('\n')
		}
		for j := 0; j < 32; j++ {
			buf.WriteByte(".#"[grid.get(i, j)])
		}
	}
	return buf.String()
}

func gameOfLife() []Scorecard {
	grid := randGrid()
	return []Scorecard{{
		Args:   []string{grid.toString()},
		Answer: grid.step().toString(),
	}}
}
