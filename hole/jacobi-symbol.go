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
	type input struct {
		a, n int
	}
	inputs := []input{}
	tests := []test{}

	var a, n int

	for i := 0; i < mult; i++ {
		// Random numbers, varying size, a<n
		n = randomOdd(1 << (53 - 2*i))
		a = randomNatural(n)
		inputs = append(inputs, input{a, n})
		// Random numbers, varying size, a>n
		a = randomNatural(1 << (53 - 2*i))
		n = randomOdd(a)
		inputs = append(inputs, input{a, n})
		// Semiprime n
		a = randomNatural(1 << 53)
		n1 := randomPrime(1 << 26)
		n2 := randomPrime(1 << 26)
		inputs = append(inputs, input{a, n1 * n2})
		// Prime n
		a = randomNatural(1 << 53)
		n = randomPrime(1 << 53)
		inputs = append(inputs, input{a, n})
		// Common multiple
		common := randomOdd(1 << 26)
		a = randomNatural(1 << 27)
		n = randomOdd(1 << 27)
		inputs = append(inputs, input{a * common, n * common})
		// a multiple of n
		a = randomNatural(1 << 26)
		n = randomOdd(1 << 27)
		inputs = append(inputs, input{a * n, n})
		// n multiple of a
		a = randomOdd(1 << 26)
		n = randomOdd(1 << 27)
		inputs = append(inputs, input{a, n * a})
		// a is small
		n = randomOdd(1 << 53)
		inputs = append(inputs, input{i, n})
		// n is small
		a = randomNatural(1 << 53)
		inputs = append(inputs, input{a, 2*i + 1})

	}
	// Different residue classes
	for i := 0; i < 4; i++ {
		for j := 1; j < 8; j += 2 {
			n = randomInClass(1<<53, 8, j)
			a = randomInClass(1<<53, 4, i)
			inputs = append(inputs, input{a, n})
		}
	}

	for _, inp := range inputs {
		tests = append(tests, test{
			strconv.Itoa(inp.a) + " " + strconv.Itoa(inp.n),
			strconv.Itoa(big.Jacobi(big.NewInt(int64(inp.a)), big.NewInt(int64(inp.n)))),
		})
	}
	return outputTests(shuffle(tests))
}
