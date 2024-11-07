package hole

import (
	"math/rand/v2"
	"strings"
)

var (
	notes = [...][2]string{
		{"C", "B♯"},
		{"C♯", "D♭"},
		{"D", "D"},
		{"D♯", "E♭"},
		{"E", "F♭"},
		{"F", "E♯"},
		{"F♯", "G♭"},
		{"G", "G"},
		{"G♯", "A♭"},
		{"A", "A"},
		{"A♯", "B♭"},
		{"B", "C♭"},
	}
	triadTypes = [...]string{
		"°", "m", "", "+",
	}
	triadSteps = [...][2]int{
		{3, 3},
		{3, 4},
		{4, 3},
		{4, 4},
	}
	orderings = [...][3]int{
		{0, 1, 2},
		{0, 2, 1},
		{1, 0, 2},
		{1, 2, 0},
		{2, 0, 1},
		{2, 1, 0},
	}
)

func letterVal(note string) byte { return note[0] - 'A' }

func genNotes(rootIdx int, rootNote string, steps [2]int) []string {
	thirdIdx := (rootIdx + steps[0]) % 12
	fifthIdx := (rootIdx + steps[0] + steps[1]) % 12
	thirdNote := notes[thirdIdx][0]
	fifthNote := notes[fifthIdx][0]

	// Enforce strict spelling. The third should be 2 letters
	// above the root, and the fifth should be 4 letters above,
	// wrapping at G
	if (letterVal(rootNote)+2)%7 != letterVal(thirdNote) {
		thirdNote = notes[thirdIdx][1]
	}
	if (letterVal(rootNote)+4)%7 != letterVal(fifthNote) {
		fifthNote = notes[fifthIdx][1]
	}

	// Return empty if strict spelling is impossible
	if (letterVal(rootNote)+2)%7 != letterVal(thirdNote) || (letterVal(rootNote)+4)%7 != letterVal(fifthNote) {
		return []string{}
	}
	return []string{rootNote, thirdNote, fifthNote}
}

func musicalChords() []Run {
	var tests []test

	// Skip a random combination for anti-cheese
	skipNum := rand.IntN(61)
	combNum := 0

	for rootIdx, rootNames := range notes {
		// Loop once for each unique name the note has
		uniqueNames := 2
		if rootNames[0] == rootNames[1] {
			uniqueNames = 1
		}

		for _, rootNote := range rootNames[:uniqueNames] {
			for triadIdx, triad := range triadTypes {
				steps := triadSteps[triadIdx]
				chordNotes := genNotes(rootIdx, rootNote, steps)
				if len(chordNotes) > 0 {
					if skipNum != combNum {
						chord := rootNote + triad
						for _, ordering := range orderings {
							rearrangedNotes := []string{chordNotes[ordering[0]], chordNotes[ordering[1]], chordNotes[ordering[2]]}
							tests = append(tests, test{
								strings.Join(rearrangedNotes, " "), chord,
							})
						}
					}
					combNum++
				}

			}
		}
	}

	// Cut 3 tests.
	return outputTests(shuffle(tests)[3:])
}
