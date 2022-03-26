package hole

import (
	"math/big"
	"math/rand"
	"strconv"
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
	const mult = 20
	type input struct {
		n, k *big.Int
	}
	inputs := make([]input, 4*mult)
	tests := make([]test, 4*mult)

	for i := 0; i < mult; i++ {
		inputs[i] = input{randomPrime(53), randomOdd(53)}
	}
	for i := mult; i < 2*mult; i++ {
		inputs[i] = input{new(big.Int).Mul(randomPrime(26), randomPrime(27)), randomOdd(53)}
	}
	for i := 2 * mult; i < 3*mult; i++ {
		inputs[i] = input{randomNatural(53), randomOdd(53)}
	}
	for i := 3 * mult; i < 4*mult; i++ {
		common := randomOdd(26)
		n := new(big.Int).Mul(randomOdd(27), common)
		k := new(big.Int).Mul(randomOdd(27), common)
		inputs[i] = input{n, k}
	}

	for i, inp := range inputs {
		tests[i] = test{
			inp.n.String() + " " + inp.k.String(),
			strconv.Itoa(big.Jacobi(inp.n, inp.k)),
		}
	}
	return outputTests(shuffle(tests))
}
