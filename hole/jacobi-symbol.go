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
	const mult = 20
	type input struct {
		n, k int
	}
	inputs := []input{{221, 7594623}}
	tests := []test{}

	for i := 0; i < mult; i++ { // Random numbers, varying size
		k := randomOdd(1 << (53 - 2*i))
		n := randomNatural(k)
		inputs = append(inputs, input{n, k})
	}
	for i := 0; i < mult; i++ { // Semiprime
		k := randomOdd(1 << 53)
		n1 := randomPrime(1 << 26)
		n2 := randomPrime(k / n1)
		inputs = append(inputs, input{n1 * n2, k})
	}
	for i := 0; i < mult; i++ { // Prime
		k := randomOdd(1 << 53)
		n := randomPrime(k)
		inputs = append(inputs, input{n, k})
	}
	for i := 0; i < mult; i++ { // Common multiple
		common := randomOdd(1 << 26)
		k := randomOdd(1 << 27)
		n := randomNatural(k)
		inputs = append(inputs, input{n * common, k * common})
	}
	for i := 0; i < 4; i++ { // Different residue classes
		for j := 1; j < 8; j += 2 {
			k := randomInClass(1<<53, 8, j)
			n := randomInClass(k, 4, i)
			inputs = append(inputs, input{n, k})
		}
	}

	for _, inp := range inputs {
		tests = append(tests, test{
			strconv.Itoa(inp.n) + " " + strconv.Itoa(inp.k),
			strconv.Itoa(big.Jacobi(big.NewInt(int64(inp.n)), big.NewInt(int64(inp.k)))),
		})
	}
	return outputTests(shuffle(tests))
}
