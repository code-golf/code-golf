func grayCode(reverse bool) []Scorecard {
	const count = 2000

	// Start with the 3 edge cases.
	numbers := make([]int, count)
	numbers[0] = 0
	numbers[1] = 1
	numbers[2] = 4096

	// Append another random cases.
	for i := 3; i < count; i++ {
		numbers[i] = randInt(1, 4096)
	}

	tests := make([]test, count)
	for i, n := range shuffle(numbers) {
		dec := strconv.Itoa(n)
		rbc := strconv.FormatInt(n ^ n>>1, 2)

		if reverse {
			tests[i] = test{rbc, dec}
		} else {
			tests[i] = test{dec, rbc}
		}
	}

	return outputTests(tests)
}