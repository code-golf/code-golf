package hole

func rockPaperScissorsSpockLizard() []Scorecard {
	tests := shuffle([]test{
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
	})

	mid := len(tests) / 2
	return outputTests(tests, tests[:mid], tests[mid:])
}
