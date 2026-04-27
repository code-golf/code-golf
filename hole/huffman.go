package hole

import (
	"container/heap"
	"math/rand/v2"
	"strings"
)

func randString(n int) string {
	alphabet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = alphabet[rand.IntN(len(alphabet))]
	}
	return string(b)
}

func randLength() int {
	return rand.IntN(22) + 10 // [10, 31]
}

type node struct {
	ch    rune
	freq  int
	left  *node
	right *node
}

type item struct {
	freq int
	id   int
	n    *node
}

type pq []item

func (p pq) Len() int { return len(p) }
func (p pq) Less(i, j int) bool {
	if p[i].freq == p[j].freq {
		return p[i].id < p[j].id
	}
	return p[i].freq < p[j].freq
}
func (p pq) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func (p *pq) Push(x any) { *p = append(*p, x.(item)) }
func (p *pq) Pop() any {
	old := *p
	n := len(old)
	x := old[n-1]
	*p = old[:n-1]
	return x
}

func huffmanEncode(s string) string {
	if len(s) == 0 {
		return ""
	}

	// frequency + first occurrence index (like Python s.index)
	freq := map[rune]int{}
	first := map[rune]int{}

	for i, c := range s {
		freq[c]++
		if _, ok := first[c]; !ok {
			first[c] = i
		}
	}

	// sort keys by (freq, first occurrence)
	type pair struct {
		ch   rune
		freq int
		idx  int
	}

	arr := []pair{}
	for c := range freq {
		arr = append(arr, pair{c, freq[c], first[c]})
	}

	// manual stable selection sort
	for i := 0; i < len(arr); i++ {
		for j := i + 1; j < len(arr); j++ {
			if arr[j].freq < arr[i].freq ||
				(arr[j].freq == arr[i].freq && arr[j].idx < arr[i].idx) {
				arr[i], arr[j] = arr[j], arr[i]
			}
		}
	}

	h := &pq{}
	heap.Init(h)

	id := 0

	// push in EXACT sorted order
	for _, p := range arr {
		heap.Push(h, item{
			freq: p.freq,
			id:   id,
			n: &node{
				ch:   p.ch,
				freq: p.freq,
			},
		})
		id++
	}

	// build tree
	for h.Len() > 1 {
		a := heap.Pop(h).(item)
		b := heap.Pop(h).(item)

		merged := &node{
			freq:  a.freq + b.freq,
			left:  a.n,
			right: b.n,
		}

		heap.Push(h, item{
			freq: merged.freq,
			id:   id,
			n:    merged,
		})
		id++
	}

	root := heap.Pop(h).(item).n

	codes := map[rune]string{}

	var build func(*node, string)
	build = func(n *node, path string) {
		if n.left == nil && n.right == nil {
			codes[n.ch] = path
			return
		}
		if n.left != nil {
			build(n.left, path+"0")
		}
		if n.right != nil {
			build(n.right, path+"1")
		}
	}

	build(root, "")

	var sb strings.Builder
	for _, c := range s {
		sb.WriteString(codes[c])
	}

	return sb.String()
}

var _ = answerFunc("huffman-encoder", func() []Answer {
	tests := []test{
		{"baabddddde", "111101011100000110"},
		{"AAAB", "1110"},
		{"ILOVECHICKEN123", "00110001001101001001110110010111100010110111101111000"},
	}

	for i := 0; i < 10; i++ {
		// generate test case of random alphanumeric of length 10 to 30
		length := randLength()
		testInput := randString(length)
		testOutput := huffmanEncode(testInput)
		tests = append(tests, test{testInput, testOutput})
	}

	shuffle(tests)

	return outputTests(tests)
})
