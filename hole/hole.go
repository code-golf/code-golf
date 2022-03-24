package hole

import "strings"

type test struct{ in, out string }

func outputTests(tests []test) ([]string, string) {
	args := make([]string, len(tests))
	var answer strings.Builder

	for i, t := range tests {
		args[i] = t.in

		if i > 0 {
			answer.WriteByte('\n')
		}
		answer.WriteString(t.out)
	}

	return args, answer.String()
}
