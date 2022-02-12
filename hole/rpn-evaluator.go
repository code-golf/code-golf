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

	var left_val int
	var right_val int

	switch op {
	case '+':
		left_val = randInt(0, val)
		right_val = val - left_val
	case '-':
		left_val = randInt(val, MaxInt)
		right_val = left_val - val
	case '*':
		if val == 0 {
			left_val = randInt(0, MaxInt)
			right_val = 0
		} else {
			factors := []int{1}
			for i := 2; i*i <= val; i++ {
				if val%i == 0 {
					factors = append(factors, i)
				}
			}
			left_val = factors[rand.Intn(len(factors))]
			right_val = val / left_val
		}
		if rand.Intn(2) == 1 {
			left_val, right_val = right_val, left_val
		}
	case '/':
		if val == 0 {
			left_val = 0
			right_val = randInt(0, MaxInt)
		} else {
			right_val = randInt(1, MaxInt/val)
			left_val = val * right_val
		}
	}

	node.op = op
	node.left = asNode(left_val)
	node.right = asNode(right_val)
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
		nodeId := rand.Intn(nodesCount)
		node := valueNodes[nodeId]
		expand(node)
		valueNodes[nodeId] = node.left
		valueNodes = append(valueNodes, node.right)
	}
}

func printNode(sb *strings.Builder, node *Node) {
	if node.op == '=' {
		if sb.Len() > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(node.value))
	} else {
		printNode(sb, node.left)
		printNode(sb, node.right)
		sb.WriteByte(' ')
		sb.WriteByte(node.op)
	}
}

func genExpr(initVal int, expander func(*Node, int), expandCount int) *Node {
	init := asNode(initVal)
	expander(init, expandCount)
	return init
}

func rpnEvaluator() ([]string, string) {
	const TestsCount = 20

	exprs := [TestsCount]*Node{
		genExpr(randInt(1, MaxInt), expandLeft, randInt(16, 31)),
		genExpr(randInt(1, MaxInt), expandRight, randInt(16, 31)),
		genExpr(randInt(1, MaxInt), expandRight, 0),
		genExpr(0, expandRand, randInt(16, 31)),
	}

	for i := 4; i < TestsCount; i++ {
		exprs[i] = genExpr(randInt(1, MaxInt), expandRand, randInt(1, 31))
	}

	rand.Shuffle(len(exprs), func(i, j int) {
		exprs[i], exprs[j] = exprs[j], exprs[i]
	})

	args := make([]string, TestsCount)
	var answer strings.Builder

	for k, expr := range exprs {
		var arg strings.Builder
		printNode(&arg, expr)
		args[k] = arg.String()
		if k > 0 {
			answer.WriteByte('\n')
		}
		answer.WriteString(strconv.Itoa(expr.value))
	}

	return args, answer.String()
}
