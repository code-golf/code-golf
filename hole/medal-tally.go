package hole

import (
	"fmt"
	"math/rand/v2"
	"slices"
)

var _ = answerFunc("medal-tally", func() []Answer {
	var tests []test

	for size := 3; size <= 10; size++ {
		// 3-9 golds
		for j := 3; j <= size-1; j++ {
			for i := 0; i <= size-j; i++ {
				tests = append(tests, medalTallyTest(size, j, 0, 0, i))
			}
		}

		// 2 golds, 1+ bronze
		for i := 1; i <= size-2; i++ {
			tests = append(tests, medalTallyTest(size, 2, 0, i, 0))
		}

		// 1 gold, 2+ silver
		for i := 2; i <= size-1; i++ {
			tests = append(tests, medalTallyTest(size, 1, i, 0, 0))
		}

		// 1 gold, 1 silver, 1+ bronze
		for i := 1; i <= size-2; i++ {
			tests = append(tests, medalTallyTest(size, 1, 1, i, 0))
		}
	}

	// Special cases where size < 3
	tests = append(tests, medalTallyTest(2, 2, 0, 0, 0))
	tests = append(tests, medalTallyTest(2, 1, 1, 0, 0))
	tests = append(tests, medalTallyTest(1, 1, 0, 0, 0))

	return outputTests(shuffle(tests))
})

func medalTallyTest(size, g, s, b, tie int) test {
	uniqueScores := rand.Perm(99)[:size+4]
	slices.Sort(uniqueScores)

	for i := range uniqueScores {
		uniqueScores[i]++
	}

	scores := make([]int, size)
	idx := 0

	for range g {
		scores[idx] = uniqueScores[0]
		idx++
	}

	for range s {
		scores[idx] = uniqueScores[1]
		idx++
	}

	for range b {
		scores[idx] = uniqueScores[2]
		idx++
	}

	for range tie {
		scores[idx] = uniqueScores[3]
		idx++
	}

	for uniqIdx := 4; idx < size; uniqIdx++ {
		scores[idx] = uniqueScores[uniqIdx]
		idx++
	}

	var in, out string

	for i := range size {
		if i != 0 {
			in += " "
		}
		in += fmt.Sprint(scores[i])
	}

	if g == 1 {
		out += "1💎 "
	}
	out += fmt.Sprint(g, "🥇 ")
	if s > 0 {
		out += fmt.Sprint(s, "🥈 ")
	}
	if b > 0 {
		out += fmt.Sprint(b, "🥉 ")
	}

	return test{in, out[:len(out)-1]}
}
