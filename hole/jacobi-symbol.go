package hole

import (
	"fmt"
	"math/big"
	"math/rand"
)

func randomNatural(max int) int {
	if max <= 0 {
		return 0
	}
	return rand.Intn(max)
}

func randomInClass(max, mod, class int) int {
	return (randomNatural(max-mod) &^ (mod - 1)) | class
}

func randomOdd(max int) int {
	return randomInClass(max, 2, 1)
}

func randomPrime(max int) int {
	k := randomInClass(max, 2, 1)
	for !big.NewInt(int64(k)).ProbablyPrime(20) {
		k = randomOdd(max)
	}
	return k
}

func jacobiSymbol() []Scorecard {
	const mult = 9
	tests := []test{{"0 1", "1"}, {"4622568476421908 4170463869060991", "1"}}
	addTest := func(a, n int) {
		tests = append(tests, test{
			fmt.Sprint(a, n),
			fmt.Sprint(big.Jacobi(big.NewInt(int64(a)), big.NewInt(int64(n)))),
		})
	}

	var a, n int

	for i := 0; i < mult; i++ {
		// Random numbers, varying size, a<n
		n = randomOdd(1 << (53 - 2*i))
		a = randomNatural(n)
		addTest(a, n)
		// Random numbers, varying size, a>n
		a = randomNatural(1 << (53 - 2*i))
		n = randomOdd(a)
		addTest(a, n)
		// Semiprime n
		a = randomNatural(1 << 53)
		n = randomPrime(1<<26) * randomPrime(1<<26)
		addTest(a, n)
		// Prime n
		a = randomNatural(1 << 53)
		n = randomPrime(1 << 53)
		addTest(a, n)
		// Common multiple
		common := randomOdd(1 << 26)
		a = randomNatural(1 << 27)
		n = randomOdd(1 << 27)
		addTest(a*common, n*common)
		// a multiple of n
		a = randomNatural(1 << 26)
		n = randomOdd(1 << 27)
		addTest(a*n, n)
		// n multiple of a
		a = randomOdd(1 << 26)
		n = randomOdd(1 << 27)
		addTest(a, n*a)
		// a is small
		n = randomOdd(1 << 53)
		addTest(i, n)
		// n is small
		a = randomNatural(1 << 53)
		addTest(a, 2*i+1)

	}
	// Different residue classes
	for i := 0; i < 4; i++ {
		for j := 1; j < 8; j += 2 {
			if i == 3 && j == 7 { break; }
			a = randomInClass(1<<53, 4, i)
			n = randomInClass(a, 8, j)
			addTest(a, n)
		}
	}
	return outputTests(shuffle(tests))
}
