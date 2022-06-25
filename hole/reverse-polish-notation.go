package hole

import (
	"math"
	"math/rand"
	"strconv"
	"strings"
)

type Node struct {
	op          byte
	value       int
	left, right *Node
}

func asNode(val int) *Node { return &Node{op: '=', value: val} }

func expand(node *Node) {
	val := node.value
	var left, right int

	switch node.op = "+-*/"[rand.Intn(4)]; node.op {
	case '+':
		left = randInt(0, val)
		right = val - left
	case '-':
		left = randInt(val, math.MaxInt16)
		right = left - val
	case '*':
		if val == 0 {
			left = randInt(0, math.MaxInt16)
			right = 0
		} else {
			factors := []int{1}
			for i := 2; i*i <= val; i++ {
				if val%i == 0 {
					factors = append(factors, i)
				}
			}
			left = randChoice(factors)
			right = val / left
		}
		if rand.Intn(2) == 1 {
			left, right = right, left
		}
	case '/':
		if val == 0 {
			left = 0
			right = randInt(1, math.MaxInt16)
		} else {
			right = randInt(1, math.MaxInt16/val)
			left = val * right
		}
	}

	node.left = asNode(left)
	node.right = asNode(right)
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

func genExpr(init int, expander func(*Node, int), expandCount int) *Node {
	node := asNode(init)
	expander(node, expandCount)
	return node
}

func reversePolishNotation() ([]string, string) {
	const tests = 20

	exprs := [tests]*Node{
		genExpr(randInt(1, math.MaxInt16), expandLeft, randInt(16, 31)),
		genExpr(randInt(1, math.MaxInt16), expandRight, randInt(16, 31)),
		genExpr(randInt(1, math.MaxInt16), expandRight, 0),
		genExpr(0, expandRand, randInt(16, 31)),
	}

	for i := 4; i < tests; i++ {
		exprs[i] = genExpr(randInt(1, math.MaxInt16), expandRand, randInt(1, 31))
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
