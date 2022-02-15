package hole

import (
	"math/big"
	"math/rand"
	"strconv"
	"strings"
)

func randomNatural(bits int) *big.Int {
	if bits == 64 {
		return new(big.Int).SetUint64(rand.Uint64())
	}
	if bits == 63 {
		return new(big.Int).SetInt64(rand.Int63())
	}
	return new(big.Int).SetInt64(rand.Int63n((1 << bits) - 1))
}

func randomOdd(bits int) *big.Int {
	k := randomNatural(bits)
	return k.SetBit(k, 0, 1)
}

func randomPrime(bits int) *big.Int {
	k := randomOdd(bits)
	for !k.ProbablyPrime(20) {
		k = randomOdd(bits)
	}
	return k
}

func jacobiSymbol() ([]string, string) {
	const tests = 20
	inputs := make([]*big.Int, 3*tests)

	for i := 0; i < tests; i++ {
		inputs[i] = randomPrime(64)
	}
	for i := tests; i < 2*tests; i++ {
		inputs[i] = new(big.Int).Mul(randomPrime(32), randomPrime(32))
	}
	for i := 2 * tests; i < 3*tests; i++ {
		inputs[i] = randomNatural(64)
	}

	rand.Shuffle(len(inputs), func(i, j int) {
		inputs[i], inputs[j] = inputs[j], inputs[i]
	})
	var answer strings.Builder
	args := make([]string, 3*tests)

	for i, n := range inputs {
		k := randomOdd(64)
		args[i] = n.String() + " " + k.String()
		if i > 0 {
			answer.WriteByte('\n')
		}
		answer.WriteString(strconv.Itoa(big.Jacobi(n, k)))
	}

	return args, answer.String()
}
