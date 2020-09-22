package hole

import (
	"math/rand"
	"strings"
)

func emojify() ([]string, string) {
	args := []string{
		":-D", ":-)", ":-|", ":-(", `:-\`, ":-*", ":-O", ":-#", "':-D",
		"':-(", ":'-)", ":'-(", ":-P", ";-P", "X-P", "X-)", "O:-)", ";-)",
		":-$", ":-", "B-)", ":-J", "}:-)", "}:-(", ":-@",
	}

	outs := []string{
		"😀", "🙂", "😐", "🙁", "😕", "😗", "😮", "🤐", "😅", "😓", "😂", "😢",
		"😛", "😜", "😝", "😆", "😇", "😉", "😳", "😶", "😎", "😏", "😈", "👿",
		"😡",
	}

	rand.Shuffle(len(args), func(i, j int) {
		args[i], args[j] = args[j], args[i]
		outs[i], outs[j] = outs[j], outs[i]
	})

	return args, strings.Join(outs, "\n")
}
