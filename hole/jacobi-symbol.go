package hole

import (
	"math/big"
	"math/rand"
	"strconv"
)

func randomNatural(max int) int {
	return rand.Intn(max - 1)
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

func jacobiSymbol() ([]string, string) {
	const mult = 9
	tests := []test{}
	addTest := func(a, n int) {
		tests = append(tests, test{
			strconv.Itoa(a) + " " + strconv.Itoa(n),
			strconv.Itoa(big.Jacobi(big.NewInt(int64(a)), big.NewInt(int64(n)))),
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
			n = randomInClass(1<<53, 8, j)
			a = randomInClass(1<<53, 4, i)
			addTest(a, n)
		}
	}
	return outputTests(shuffle(tests))
}
