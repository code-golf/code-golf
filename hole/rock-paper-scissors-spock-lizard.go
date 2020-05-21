package hole

import (
	"math/rand"
	"strings"
)

var rockPaperScissorsSpockLizardGames = [...][2]string{
	{"💎💎", "Tie"},
	{"💎📄", "📄 covers 💎"},
	{"💎✂", "💎 crushes ✂"},
	{"💎🖖", "🖖 vaporizes 💎"},
	{"💎🦎", "💎 crushes 🦎"},
	{"📄💎", "📄 covers 💎"},
	{"📄📄", "Tie"},
	{"📄✂", "✂ cuts 📄"},
	{"📄🖖", "📄 disproves 🖖"},
	{"📄🦎", "🦎 eats 📄"},
	{"✂💎", "💎 crushes ✂"},
	{"✂📄", "✂ cuts 📄"},
	{"✂✂", "Tie"},
	{"✂🖖", "🖖 smashes ✂"},
	{"✂🦎", "✂ decapitates 🦎"},
	{"🖖💎", "🖖 vaporizes 💎"},
	{"🖖📄", "📄 disproves 🖖"},
	{"🖖✂", "🖖 smashes ✂"},
	{"🖖🖖", "Tie"},
	{"🖖🦎", "🦎 poisons 🖖"},
	{"🦎💎", "💎 crushes 🦎"},
	{"🦎📄", "🦎 eats 📄"},
	{"🦎✂", "✂ decapitates 🦎"},
	{"🦎🖖", "🦎 poisons 🖖"},
	{"🦎🦎", "Tie"},
}

func rockPaperScissorsSpockLizard() ([]string, string) {
	args := make([]string, len(rockPaperScissorsSpockLizardGames))
	outs := make([]string, len(rockPaperScissorsSpockLizardGames))

	for i, game := range rockPaperScissorsSpockLizardGames {
		args[i] = game[0]
		outs[i] = game[1]
	}

	rand.Shuffle(len(args), func(i, j int) {
		args[i], args[j] = args[j], args[i]
		outs[i], outs[j] = outs[j], outs[i]
	})

	return args, strings.Join(outs, "\n")
}
