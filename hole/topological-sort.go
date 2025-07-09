package hole

import (
	"fmt"
	"math/rand/v2"
	"slices"
	"strings"
)

// does not support isolated nodes
// cycles are insurmountable obstructions
func tsort(in string) string {
	children, in_degrees := make(map[string][]string), make(map[string]int)
	for pair := range slices.Chunk(strings.Fields(in), 2) {
		children[pair[0]] = append(children[pair[0]], pair[1])
		in_degrees[pair[0]] = in_degrees[pair[0]] // include zero-in-degree nodes
		if pair[0] != pair[1] {
			in_degrees[pair[1]]++
		}
	}
	var stack []string
	for node, deg := range in_degrees {
		if deg == 0 {
			stack = append(stack, node)
		}
	}
	var result []string
	for len(stack) > 0 {
		next := stack[len(stack)-1]
		stack = slices.Delete(stack, len(stack)-1, len(stack))
		result = append(result, next)
		for _, child := range children[next] {
			in_degrees[child]--
			if in_degrees[child] == 0 {
				stack = append(stack, child)
			}
		}
	}
	return strings.Join(result, " ")
}

// guaranteed Hamiltonian path, so the output is unique
func generateTsortTest(length, additional int) string {
	init := rand.Perm(length)
	var (
		seq  []int
		rows int
	)
	arc_choices := rand.Perm((length - 2) * (length - 1) / 2)
	for i, x := range init {
		for j, y := range init[i+1:] {
			if j == 0 || arc_choices[len(arc_choices)-1] < additional {
				seq = append(seq, x, y)
				rows++
			}
			if j != 0 {
				arc_choices = slices.Delete(arc_choices, len(arc_choices)-1, len(arc_choices))
			}
		}
	}
	var test string
	for i, x := range rand.Perm(rows) {
		if i > 0 {
			test += "\n"
		}
		test += fmt.Sprint(seq[2*x], seq[2*x+1])
	}
	return test
}

// single-digit numbers
var _ = answerFunc("topological-sort", func() []Answer {
	var tests []test
	for n := range 10 - 2 + 1 {
		n += 2
		for a := range (n-2)*(n-1)/2 + 1 {
			for range 3 {
				tcase := generateTsortTest(n, a)
				tests = append(tests, test{tcase, tsort(tcase)})
			}
		}
	}
	return outputTests(shuffle(tests))
})
