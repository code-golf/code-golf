package hole

import (
	"math/rand/v2"
	"strconv"
)

// G2 tap pairs for PRN 1-32, IS-GPS-200 Table 3-Ia.
var gpsTaps = [32][2]int{
	{2, 6}, {3, 7}, {4, 8}, {5, 9}, {1, 9}, {2, 10}, {1, 8}, {2, 9},
	{3, 10}, {2, 3}, {3, 4}, {5, 6}, {6, 7}, {7, 8}, {8, 9}, {9, 10},
	{1, 4}, {2, 5}, {3, 6}, {4, 7}, {5, 8}, {6, 9}, {1, 3}, {4, 6},
	{5, 7}, {6, 8}, {7, 9}, {8, 10}, {1, 6}, {2, 7}, {3, 8}, {4, 9},
}

// Generate the 1023-chip C/A code of the given PRN as a string of 0s and 1s.
func gpsCode(prn int) []byte {
	var g1, g2 [11]byte // Stages 1-10, index 0 unused.
	for i := 1; i <= 10; i++ {
		g1[i], g2[i] = 1, 1
	}

	a, b := gpsTaps[prn-1][0], gpsTaps[prn-1][1]
	code := make([]byte, 1023)
	for i := range code {
		code[i] = '0' + (g1[10] ^ g2[a] ^ g2[b])

		f1 := g1[3] ^ g1[10]
		f2 := g2[2] ^ g2[3] ^ g2[6] ^ g2[8] ^ g2[9] ^ g2[10]
		copy(g1[2:], g1[1:10])
		copy(g2[2:], g2[1:10])
		g1[1], g2[1] = f1, f2
	}

	return code
}

// Rotate the code left by shift chips and flip the chips at the given positions.
func gpsReceive(prn, shift int, flips ...int) test {
	code := gpsCode(prn)
	rx := append(code[shift:], code[:shift]...)
	for _, p := range flips {
		rx[p] ^= 1
	}
	return test{string(rx), strconv.Itoa(prn)}
}

const gpsMaxFlips = 8

var _ = answerFunc("gps-satellite-identification", func() []Answer {
	randFlips := func() []int {
		return rand.Perm(1023)[:rand.IntN(gpsMaxFlips+1)]
	}

	// Every PRN at least once, with random shifts and 0-8 flipped chips.
	var tests []test
	for prn := 1; prn <= 32; prn++ {
		tests = append(tests, gpsReceive(prn, rand.IntN(1023), randFlips()...))
	}

	// Evenly spread flips leave no long error-free stretch, so looking at
	// only part of the recording can't work.
	offset := rand.IntN(1023 / gpsMaxFlips)
	spread := make([]int, gpsMaxFlips)
	for i := range spread {
		spread[i] = offset + i*1023/gpsMaxFlips
	}

	// A burst of adjacent flips, wrapping around the end of the recording.
	start := rand.IntN(1023)
	burst := make([]int, gpsMaxFlips)
	for i := range burst {
		burst[i] = (start + i) % 1023
	}

	// Edge cases: no shift, maximal shift, no errors, errors at both ends,
	// errors at the very start, spread and burst errors.
	tests = append(tests,
		gpsReceive(randInt(1, 32), 0),
		gpsReceive(randInt(1, 32), 0, 0, 1022),
		gpsReceive(randInt(1, 32), 1022, randFlips()...),
		gpsReceive(randInt(1, 32), rand.IntN(1023), 0, 1, 2, 3, 4, 5, 6, 7),
		gpsReceive(randInt(1, 32), rand.IntN(1023), spread...),
		gpsReceive(randInt(1, 32), rand.IntN(1023), burst...),
	)

	for len(tests) < 40 {
		tests = append(tests, gpsReceive(randInt(1, 32), rand.IntN(1023), randFlips()...))
	}

	return outputMultirunTests(tests)
})
