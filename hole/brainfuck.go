package hole

import (
	"fmt"
	"strings"
)

var _ = answerFunc("brainfuck", func() []Answer {
	tests := fixedTests("brainfuck")
	tests = append(tests, randomBFCase(5, 3, 8, -3, 2))
	tests = append(tests, randomBFCase(10, 5, 1, 1, -3))

	for range 3 {
		jumpSize := randInt(5, 12)
		buckets := randInt(2, 8)
		initialBucketSize := randInt(1, 7)
		bucketSizeChange := randInt(-1, 2)
		charShift := randInt(-3, 3)
		tests = append(tests, randomBFCase(jumpSize, buckets, initialBucketSize, bucketSizeChange, charShift))
	}

	shuffle(tests)
	const argc = 12 // Preserve original argc
	return outputTests(tests[:argc], tests[len(tests)-argc:])
})

func intToBFString(n int) string {
	if n < 1 {
		return strings.Repeat("-", -n)
	}
	return strings.Repeat("+", n)
}

func randomBFCase(jumpSize, buckets, initialBucketSize, bucketSizeChange, charShift int) test {
	// Only generates ASCII values in a range that was already used by fixed cases before random cases were added, in order to not break legacy solutions.
	const ASCIIMin = 50
	const ASCIIMax = 122

	bucketString := ""
	out := ""

	// Ensure all bucket sizes are strictly positive
	if bucketSizeChange*buckets+initialBucketSize < 1 {
		bucketSizeChange = 1
	}

	// Populate buckets
	bucketSize := initialBucketSize
	for range buckets {
		// Choose a base for this bucket that ensures all chars remain in the printable ASCII range
		totalShift := charShift * (bucketSize - 1)
		minChar := ASCIIMin
		maxChar := ASCIIMax
		if totalShift > 0 {
			maxChar -= totalShift
		} else {
			minChar -= totalShift
		}
		base := randInt(minChar/jumpSize+1, maxChar/jumpSize)
		bucketString += ">" + intToBFString(base)

		// Add all the chars in the bucket to the solution
		char := base * jumpSize
		for range bucketSize {
			out += fmt.Sprintf("%c", char)
			char += charShift
		}

		bucketSize += bucketSizeChange
	}

	in := fmt.Sprintf(">%s[%s[<]>-]<%s>>[<<[>+>.%s<<-]>%s>[-]>[<]<<[>++++++++++.>]>>>]",
		intToBFString(jumpSize),
		bucketString,
		intToBFString(initialBucketSize),
		intToBFString(charShift),
		intToBFString(bucketSizeChange))

	return test{in, out}
}
