package hole

import (
	"math/rand"
	"strings"
)

var (
	notes = [12][2]string{
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
	triadTypes = [4]string{
		"°", "m", "", "+",
	}
	triadSteps = [4][2]int{
		{3, 3},
		{3, 4},
		{4, 3},
		{4, 4},
	}
	orderings = [6][3]int{
		{0, 1, 2},
		{0, 2, 1},
		{1, 0, 2},
		{1, 2, 0},
		{2, 0, 1},
		{2, 1, 0},
	}
)

func letterVal(note string) int {
	return int(note[0]) - 65
}

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

func musicalChords() (args []string, out string) {
	var outs []string

	// Skip a random combination for anti-cheese
	skipNum := rand.Intn(61)
	combNum := 0

	for rootIdx, rootNames := range notes {

		// Loop once for each unique name the note has
		uniqueNames := 2
		if rootNames[0] == rootNames[1] {
			uniqueNames = 1
		}
		for _, rootNote := range rootNames[:uniqueNames] {
			for triadIdx := 0; triadIdx < 4; triadIdx++ {
				triad := triadTypes[triadIdx]
				steps := triadSteps[triadIdx]
				chordNotes := genNotes(rootIdx, rootNote, steps)
				if len(chordNotes) > 0 {
					if skipNum != combNum {
						chord := rootNote + triad
						for _, ordering := range orderings {
							rearrangedNotes := []string{chordNotes[ordering[0]], chordNotes[ordering[1]], chordNotes[ordering[2]]}
							args = append(args, strings.Join(rearrangedNotes, " "))
							outs = append(outs, chord)
						}
					}
					combNum = combNum + 1
				}

			}
		}
	}
	// shuffle args and outputs in the same way
	rand.Shuffle(len(args), func(i, j int) {
		args[i], args[j] = args[j], args[i]
		outs[i], outs[j] = outs[j], outs[i]
	})

	// Cut 3 args
	args = args[3:]
	outs = outs[3:]

	out = strings.Join(outs, "\n")
	return
}
