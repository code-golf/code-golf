package hole

import (
	"slices"
	"strings"
)

var _ = answerFunc("suit-up", func() []Answer {
	answers := make([]Answer, 3)
	for i := range answers {
		answers[i] = suitUp()
	}
	return answers
})

func suitUp() Answer {
	var deck []string
	for _, suit := range "SHDC" {
		for _, rank := range "A23456789TJQK" {
			deck = append(deck, string(rank)+string(suit))
		}
	}
	shuffle(deck)

	args := slices.Clone(deck)
	piles := map[byte][]string{}
	var revealed []string

	for len(deck) > 0 {
		card := deck[0]
		deck = deck[1:]
		rank, suit := card[0], card[1]
		pile := piles[suit]

		if len(pile) == 0 || '6' <= rank && rank <= '9' {
			// Empty pile or 6-9: place the card on its suit pile.
			piles[suit] = append(pile, card)
			continue
		}

		// Remove the top card of the pile: discarded unseen by 2-5,
		// revealed by A, T, J, Q, K. The drawn card goes to the bottom.
		top := pile[len(pile)-1]
		piles[suit] = pile[:len(pile)-1]
		if rank < '2' || rank > '5' {
			revealed = append(revealed, top)
		}
		deck = append(deck, card)
	}

	return Answer{Args: args, Answer: strings.Join(revealed, "\n")}
}
