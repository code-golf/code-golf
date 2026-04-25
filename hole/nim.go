package hole

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

var _ = answerFunc("nim", func() []Answer {
	tests := fixedTests("nim")

	findBestMove := func(piles []int) (int, int, int) {
		nimSum := 0
		for _, p := range piles {
			nimSum ^= p
		}

		count := 0
		bestPile := -1
		bestAmount := -1

		for i, p := range piles {
			target := nimSum ^ p
			if target < p {
				count++
				bestPile = i
				bestAmount = p - target
			}
		}

		return bestPile, bestAmount, count
	}

	addTest := func(sortedPiles []int) bool {
		shuffled := make([]int, len(sortedPiles))
		copy(shuffled, sortedPiles)
		rand.Shuffle(len(shuffled), func(i, j int) {
			shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
		})

		pileIdx, amount, count := findBestMove(shuffled)
		// check if there is only one best move
		if count != 1 {
			return false
		}

		parts := make([]string, len(shuffled))
		for i, p := range shuffled {
			parts[i] = strconv.Itoa(p)
		}
		input := strings.Join(parts, " ")
		output := fmt.Sprintf("%d %d", pileIdx, amount)

		tests = append(tests, test{input, output})
		return true
	}

    // sorted ordering to ensure no duplicates
	for a := 1; a <= 5; a++ {
		addTest([]int{a})
		for b := 1; b <= a; b++ {
			addTest([]int{a, b})
			for c := 1; c <= b; c++ {
				addTest([]int{a, b, c})
				for d := 1; d <= c; d++ {
					addTest([]int{a, b, c, d})
					for e := 1; e <= d; e++ {
						addTest([]int{a, b, c, d, e})
					}
				}
			}
		}
	}

	return outputTests(shuffle(tests))
})
