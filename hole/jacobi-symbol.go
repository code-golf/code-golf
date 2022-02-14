package hole

import (
	"crypto/rand"
	"math/big"
	mrand "math/rand"
	"strconv"
	"strings"
)

func randomOdd() *big.Int {
	k, _ := rand.Int(rand.Reader, new(big.Int).SetUint64(^uint64(0)))
	return k.SetBit(k, 0, 1)
}

func jacobiSymbol() ([]string, string) {
	const tests = 20
	inputs := make([]*big.Int, 3*tests)

	for i := 0; i < tests; i++ {
		inputs[i], _ = rand.Prime(rand.Reader, 64)
	}
	for i := tests; i < 2*tests; i++ {
		p1, _ := rand.Prime(rand.Reader, 32)
		p2, _ := rand.Prime(rand.Reader, 32)
		inputs[i] = new(big.Int).Mul(p1, p2)
	}
	for i := 2*tests; i < 3*tests; i++ {
		inputs[i], _ = rand.Int(rand.Reader, new(big.Int).SetUint64(^uint64(0)))
	}

	mrand.Shuffle(len(inputs), func(i, j int) {
		inputs[i], inputs[j] = inputs[j], inputs[i]
	})
	var answer strings.Builder
	args := make([]string, 3*tests)

	for i, n := range inputs {
		k := randomOdd()
		args[i] = n.String() + " " + k.String()
		if i > 0 {
			answer.WriteByte('\n')
		}
		answer.WriteString(strconv.Itoa(big.Jacobi(n, k)))
	}

	return args, answer.String()
}