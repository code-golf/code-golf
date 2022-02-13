package hole

import (
	"math/rand"
	"strconv"
	"strings"
)

func randInt(a, b int) int {
	return rand.Intn(b-a+1) + a
}

const MaxInt = 1<<15 - 1

type Node struct {
	op    byte
	value int
	left  *Node
	right *Node
}

func asNode(val int) *Node {
	return &Node{
		op:    '=',
		value: val,
		left:  nil,
		right: nil,
	}
}

func expand(node *Node) {
	val := node.value
	op := "+-*/"[rand.Intn(4)]

	var leftVal, rightVal int

	switch op {
	case '+':
		leftVal = randInt(0, val)
		rightVal = val - leftVal
	case '-':
		leftVal = randInt(val, MaxInt)
		rightVal = leftVal - val
	case '*':
		if val == 0 {
			leftVal = randInt(0, MaxInt)
			rightVal = 0
		} else {
			factors := []int{1}
			for i := 2; i*i <= val; i++ {
				if val%i == 0 {
					factors = append(factors, i)
				}
			}
			leftVal = factors[rand.Intn(len(factors))]
			rightVal = val / leftVal
		}
		if rand.Intn(2) == 1 {
			leftVal, rightVal = rightVal, leftVal
		}
	case '/':
		if val == 0 {
			leftVal = 0
			rightVal = randInt(1, MaxInt)
		} else {
			rightVal = randInt(1, MaxInt/val)
			leftVal = val * rightVal
		}
	}

	node.op = op
	node.left = asNode(leftVal)
	node.right = asNode(rightVal)
}

func expandLeft(init *Node, count int) {
	for ; count > 0; count-- {
		expand(init)
		init = init.left
	}
}

func expandRight(init *Node, count int) {
	for ; count > 0; count-- {
		expand(init)
		init = init.right
	}
}

func expandRand(init *Node, count int) {
	valueNodes := []*Node{init}
	for nodesCount := 1; nodesCount <= count; nodesCount++ {
		nodeIdx := rand.Intn(nodesCount)
		node := valueNodes[nodeIdx]
		expand(node)
		valueNodes[nodeIdx] = node.left
		valueNodes = append(valueNodes, node.right)
	}
}

func writeNode(sb *strings.Builder, node *Node) {
	if node.op == '=' {
		if sb.Len() > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(node.value))
	} else {
		writeNode(sb, node.left)
		writeNode(sb, node.right)
		sb.WriteByte(' ')
		sb.WriteByte(node.op)
	}
}

func genExpr(initVal int, expander func(*Node, int), expandCount int) *Node {
	init := asNode(initVal)
	expander(init, expandCount)
	return init
}

func reversePolishNotation() ([]string, string) {
	const tests = 20

	exprs := [tests]*Node{
		genExpr(randInt(1, MaxInt), expandLeft, randInt(16, 31)),
		genExpr(randInt(1, MaxInt), expandRight, randInt(16, 31)),
		genExpr(randInt(1, MaxInt), expandRight, 0),
		genExpr(0, expandRand, randInt(16, 31)),
	}

	for i := 4; i < tests; i++ {
		exprs[i] = genExpr(randInt(1, MaxInt), expandRand, randInt(1, 31))
	}

	rand.Shuffle(len(exprs), func(i, j int) {
		exprs[i], exprs[j] = exprs[j], exprs[i]
	})

	args := make([]string, tests)
	var answer strings.Builder

	for i, expr := range exprs {
		var arg strings.Builder
		writeNode(&arg, expr)
		args[i] = arg.String()
		if i > 0 {
			answer.WriteByte('\n')
		}
		answer.WriteString(strconv.Itoa(expr.value))
	}

	return args, answer.String()
}
