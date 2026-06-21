package hole

import (
	"slices"
	"strconv"
)

var _ = answerFunc("gray-code-decoder", func() []Answer { return grayCode(true) })
var _ = answerFunc("gray-code-encoder", func() []Answer { return grayCode(false) })

func grayCode(reverse bool) []Answer {
	const count = 2000

	// Start with the 3 edge cases.
	numbers := make([]int, count)
	numbers[0] = 0
	numbers[1] = 1
	numbers[2] = 4095

	// Append another random cases.
	for i := 3; i < count; i++ {
		numbers[i] = randInt(1, 4095)
	}

	tests := make([]test, count)
	unpermuted := make([]test, count)
	numbersCopy := slices.Clone(numbers)
	for i, n := range shuffle(numbersCopy) {
		dec := strconv.Itoa(n)
		rbc := strconv.FormatInt(int64(n^n>>1), 2)

		if reverse {
			tests[i] = test{rbc, dec}
		} else {
			tests[i] = test{dec, rbc}
			unpermuted[i] = test{strconv.Itoa(numbers[i]), strconv.FormatInt(int64(numbers[i]^numbers[i]>>1), 2)}
		}
	}

	if reverse {
		return outputTests(tests)
	}
	return outputTests(tests, unpermuted)
}
