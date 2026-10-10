package hole

import (
	"fmt"
	"math/rand/v2"
	"slices"
)

var _ = answerFunc("box-plot", func() []Answer {
	tests := make([]test, 50)

	for i := 0; i < len(tests); i++ {
		datasetlength := rand.IntN(16) + 5
		datasetmin := rand.IntN(7)
		datasetmax := rand.IntN(11-datasetmin) + datasetmin

		var dataset []int
		dataset = append(dataset, datasetmin, datasetmax)

		for range datasetlength - 2 {
			dataset = append(dataset, datasetmin+rand.IntN(datasetmax+1-datasetmin))
		}

		slices.Sort(dataset)

		inp := ""
		for _, val := range dataset {
			inp = fmt.Sprint(inp, " ", val)
		}
		inp = inp[1:]

		var q1pos, medpos, q3pos int
		medpos = datasetlength - 1
		if medpos%2 == 0 {
			q1pos = (medpos - 2) / 2
			q3pos = (medpos + 2 + datasetlength*2 - 2) / 2
		} else {
			q1pos = (medpos - 1) / 2
			q3pos = (medpos + 1 + datasetlength*2 - 2) / 2
		}

		var minval, q1, median, q3, maxval int
		minval = datasetmin * 2
		maxval = datasetmax * 2
		if q1pos%2 == 0 {
			q1 = dataset[q1pos/2] * 2
		} else {
			q1 = dataset[q1pos/2] + dataset[(1+q1pos)/2]
		}
		if medpos%2 == 0 {
			median = dataset[medpos/2] * 2
		} else {
			median = dataset[medpos/2] + dataset[(1+medpos)/2]
		}
		if q3pos%2 == 0 {
			q3 = dataset[q3pos/2] * 2
		} else {
			q3 = dataset[q3pos/2] + dataset[(1+q3pos)/2]
		}

		outp := ""
		for j := range 21 {
			x := "  "
			if q1 < j && j < q3 {
				x = "──"
			}
			if j == median {
				x = "┬─"
			}
			if j == q1 {
				x = "┌─"
			}
			if j == q3 {
				x = "┐ "
			}
			if j == q1 && j == q3 {
				x = "╷ "
			}
			outp += x
		}
		outp += "\n"
		for j := range 21 {
			x := "  "
			if minval < j && j < q1 || q3 < j && j < maxval {
				x = "──"
			}
			if j == minval {
				x = "├─"
			}
			if j == q1 {
				x = "┤ "
			}
			if j == median {
				x = "│ "
			}
			if j == q3 {
				x = "├─"
			}
			if j == maxval {
				x = "┤ "
			}
			if j == minval && j == q1 {
				x = "│ "
			}
			if j == q1 && j == median {
				x = "┤ "
			}
			if j == median && j == q3 {
				x = "├─"
			}
			if j == q3 && j == maxval {
				x = "┤ "
			}
			if j == minval && j == median {
				x = "│ "
			}
			if j == q1 && j == q3 {
				x = "┼ "
			}
			if j == median && j == maxval {
				x = "│ "
			}
			if j == minval && j == q3 {
				x = "├─"
			}
			if j == q1 && j == maxval {
				x = "┤ "
			}
			if j == minval && j == maxval {
				x = "│ "
			}
			outp += x
		}
		outp += "\n"
		for j := range 21 {
			x := "  "
			if q1 < j && j < q3 {
				x = "──"
			}
			if j == median {
				x = "┴─"
			}
			if j == q1 {
				x = "└─"
			}
			if j == q3 {
				x = "┘ "
			}
			if j == q1 && j == q3 {
				x = "╵ "
			}
			outp += x
		}
		outp += "\n└─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┘\n0   1   2   3   4   5   6   7   8   9  10"

		tests[i] = test{inp, outp}
	}

	return outputTests(shuffle(tests))
})
