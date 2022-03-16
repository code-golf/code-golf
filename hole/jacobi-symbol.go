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
	type Test struct {
		n, k *big.Int
	}
	inputs := make([]Test, 4*tests)

	for i := 0; i < tests; i++ {
		inputs[i] = Test{randomPrime(53), randomOdd(53)}
	}
	for i := tests; i < 2*tests; i++ {
		inputs[i] = Test{new(big.Int).Mul(randomPrime(26), randomPrime(27)), randomOdd(53)}
	}
	for i := 2 * tests; i < 3*tests; i++ {
		inputs[i] = Test{randomNatural(53), randomOdd(53)}
	}
	for i := 3 * tests; i < 4*tests; i++ {
		common := randomOdd(26)
		n := new(big.Int).Mul(randomOdd(27), common)
		k := new(big.Int).Mul(randomOdd(27), common)
		inputs[i] = Test{n, k}
	}

	rand.Shuffle(len(inputs), func(i, j int) {
		inputs[i], inputs[j] = inputs[j], inputs[i]
	})
	var answer strings.Builder
	args := make([]string, 4*tests)

	for i, test := range inputs {
		args[i] = test.n.String() + " " + test.k.String()
		if i > 0 {
			answer.WriteByte('\n')
		}
		answer.WriteString(strconv.Itoa(big.Jacobi(test.n, test.k)))
	}

	return args, answer.String()
}
