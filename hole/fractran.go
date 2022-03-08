package hole

import (
	"math"
	"math/rand"
	"strconv"
	"strings"
)

type Number = map[int]int
type ABPair struct {
	a, b int
}
type Program = []Number

var primes = []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47}

func Mul(a Number, b Number) Number {
	prod := Number{}
	for p, e := range a {
		prod[p] = e
	}
	for p, e := range b {
		prod[p] = prod[p] + e
		if prod[p] == 0 {
			delete(prod, p)
		}
	}
	return prod
}

func IsInteger(a Number) bool {
	for _, e := range a {
		if e < 0 {
			return false
		}
	}
	return true
}

func ToFraction(a Number) (int, int) {
	n, d := 1, 1
	for p, e := range a {
		for e > 0 {
			e -= 1
			n *= p
		}
		for e < 0 {
			e += 1
			d *= p
		}
	}
	return n, d
}
func ToString(a Number) string {
	n, d := ToFraction(a)
	return strconv.Itoa(n) + "/" + strconv.Itoa(d)
}
func ToStringNumerator(a Number) string {
	n, _ := ToFraction(a)
	return strconv.Itoa(n)
}

func Run(n Number, program Program) Number {
	for {
		halt := true
		for _, x := range program {
			newN := Mul(n, x)
			if IsInteger(newN) {
				n = newN
				halt = false
				//fmt.Println(n)
				break
			}
		}
		if halt {
			return n
		}
	}
}

func randomVars(reserved int, vars ...*int) {
	nums := []int{}
	for i := reserved; i < len(primes); i++ {
		nums = append(nums, primes[i])
	}
	rand.Shuffle(len(nums), func(i, j int) {
		nums[i], nums[j] = nums[j], nums[i]
	})
	for i := 0; i < len(vars); i++ {
		*vars[i] = nums[i]
	}
}

func SwapProgram() Program {
	//Given 2^a × 3^b × 5 returns 2^b × 3^a
	var S1, T0, T1, A int
	S0 := 5
	randomVars(3, &S1, &T0, &T1, &A)
	return Program{
		{S1: 1, S0: -1, A: 1, 3: -1}, //State S0,S1: "3"--;A++
		{S0: 1, S1: -1},
		{T0: 1, S0: -1},              //Switch to state T0,T1
		{T1: 1, T0: -1, 3: 1, 2: -1}, //State T0,T1: "2"--;"3"++
		{T0: 1, T1: -1},
		{T0: -1},
		{2: 1, A: -1}, //State default:A--;"2"++
	}
}

func AddProgram() Program {
	//Given 2^a × 3^b returns 2^(a+b)
	return Program{
		{2: 1, 3: -1},
	}
}

func MinProgram() Program {
	//Given 2^a × 3^b returns 2^min(a,b)
	var S0, S1, A int
	randomVars(2, &S0, &S1, &A)
	return Program{
		{2: -1, 3: -1, A: 1, S0: 1}, //A=min("2","3"),"2"-=A;"3"-=A;goto S
		{2: -1, S0: -1, S1: 1},      //state S:"2"=0
		{S1: -1, S0: 1},
		{S1: -1}, //goto output
		{S0: -1},
		{3: -1},       //"3"=0
		{2: 1, A: -1}, //"2"=A
	}
}

func MaxProgram() Program {
	//Given 2^a × 3^b returns 2^max(a,b)
	var S0, S1, A int
	randomVars(2, &S0, &S1, &A)
	return Program{
		{2: -1, 3: -1, A: 1, S0: 1},  //A=min("2","3"),"2"-=A;"3"-=A;goto S
		{2: -1, S0: -1, S1: 1, A: 1}, //state S:A+="2";"2"=0
		{S1: -1, S0: 1},
		{S0: -1},      //goto output
		{3: -1, A: 1}, //"3"=0
		{2: 1, A: -1}, //"2"=A
	}
}

func FibonacciProgram() Program {
	//Given 3^a returns 2^Fibonacci(a)
	var S0, S1, T0, T1, U0, U1, B, C int
	A := 2
	randomVars(2, &S0, &S1, &T0, &T1, &U0, &U1, &B, &C)
	return Program{
		{C: 1, S1: 1, S0: -1, A: -1}, //state S: C+=A;A=0
		{S0: 1, S1: -1},
		{T0: 1, S0: -1},                    //goto T
		{A: 1, C: 1, T1: 1, T0: -1, B: -1}, //state T: A+=B;C+=B;B=0
		{T0: 1, T1: -1},
		{U0: 1, T0: -1},              //goto U
		{B: 1, U1: 1, C: -1, U0: -1}, //state U: B+=C;C=0
		{U0: 1, U1: -1},
		{U0: -1, B: -1},      //B-=1; goto init
		{S0: 1, B: 1, 3: -1}, //state init "3">0: B++, goto S
		{B: -1},              //state init "3"==0: B=0
	}
}

const MAXINT = int(^uint(0) >> 1)

func RandomABPair(lim int) ABPair {
	pairs := []ABPair{}
	for a := 0; float64(a) < math.Log2(float64(lim)); a++ {
		for b := 0; float64(b) < math.Log(float64(lim)/math.Pow(2.0, float64(a)))/math.Log(3); b++ {
			pairs = append(pairs, ABPair{a, b})
		}
	}
	return pairs[rand.Intn(len(pairs))]
}

func insert(a Program, index int, value Number) Program {
	if len(a) == index { // nil or empty slice or after last element
		return append(a, value)
	}
	a = append(a[:index+1], a[index:]...) // index < len(a)
	a[index] = value
	return a
}

func Obfuscate(program Program) Program {
	//Works by doubling random instructions
	//TODO: replace with random powers of random instructions
	result := Program{}
	for i := 0; i < len(program); i++ {
		result = append(result, program[i])
	}
	for i := 0; i < 12; i++ {
		r := rand.Intn(len(result))
		result = insert(result, r, result[r])
	}
	return result
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func fractran() ([]string, string) {
	type TestInput struct {
		Input   Number
		Program Program
	}
	tests := []TestInput{}
	expected := Program{}
	for _, pair := range []ABPair{RandomABPair(MAXINT / 5), ABPair{12, 8}} {
		tests = append(tests, TestInput{Number{2: pair.a, 3: pair.b, 5: 1}, SwapProgram()})
		expected = append(expected, Number{2: pair.b, 3: pair.a})
	}
	for _, pair := range []ABPair{RandomABPair(MAXINT), ABPair{12, 8}} {
		tests = append(tests, TestInput{Number{2: pair.a, 3: pair.b}, AddProgram()})
		expected = append(expected, Number{2: (pair.a + pair.b)})
	}
	for _, pair := range []ABPair{RandomABPair(MAXINT), ABPair{12, 8}} {
		tests = append(tests, TestInput{Number{2: pair.a, 3: pair.b}, MinProgram()})
		expected = append(expected, Number{2: min(pair.a, pair.b)})
	}
	for _, pair := range []ABPair{RandomABPair(MAXINT), ABPair{12, 8}} {
		tests = append(tests, TestInput{Number{2: pair.a, 3: pair.b}, MaxProgram()})
		expected = append(expected, Number{2: -min(-pair.a, -pair.b)})
	}
	tests = append(tests, TestInput{Number{3: 7}, FibonacciProgram()})
	expected = append(expected, Number{2: 13})
	tests = append(tests, TestInput{Number{3: 10}, FibonacciProgram()})
	expected = append(expected, Number{2: 55})

	for i := 0; i < len(tests); i++ {
		if ToString(Run(tests[i].Input, tests[i].Program)) != ToString(expected[i]) {
			panic("Assertion error." + strconv.Itoa(i))
		}
	}

	rand.Shuffle(len(tests), func(i, j int) {
		tests[i], tests[j] = tests[j], tests[i]
	})

	args := make([]string, len(tests))
	var answer strings.Builder

	for i, input := range tests {
		args[i] = ToStringNumerator(input.Input)
		for _, frac := range input.Program {
			args[i] = args[i] + " " + ToString(frac)
		}
		if i > 0 {
			answer.WriteByte('\n')
		}
		answer.WriteString(ToStringNumerator(Run(input.Input, input.Program)))
	}

	return args, answer.String()
}
