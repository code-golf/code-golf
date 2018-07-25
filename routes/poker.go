package routes

import "math/rand"

var handTypes = []string{
	"High Card",
	"Pair",
	"Two Pair",
	"Three of a Kind",
	"Straight",
	"Flush",
	"Full House",
	"Four of a Kind",
	"Straight Flush",
	"Royal Flush",
}

func poker() (args []string, out string) {
	hands := []struct {
		Type  uint8
		Cards []rune
	}{
		// All the royal flushes.
		{9, []rune{'🂪', '🂫', '🂭', '🂮', '🂡'}},
		{9, []rune{'🂺', '🂻', '🂽', '🂾', '🂱'}},
		{9, []rune{'🃊', '🃋', '🃍', '🃎', '🃁'}},
		{9, []rune{'🃚', '🃛', '🃝', '🃞', '🃑'}},
		// TODO Needs more random.
		{8, []rune{'🃘', '🃗', '🃖', '🃕', '🃔'}},
		{7, []rune{'🂻', '🃋', '🂫', '🃛', '🃇'}},
		{6, []rune{'🂺', '🃊', '🂪', '🃙', '🃉'}},
		{5, []rune{'🂤', '🂫', '🂨', '🂢', '🂩'}},
		{4, []rune{'🃙', '🃈', '🂧', '🃆', '🂵'}},
		{3, []rune{'🃗', '🃇', '🂧', '🃞', '🃃'}},
		{2, []rune{'🃔', '🂤', '🃓', '🃃', '🃝'}},
		{1, []rune{'🂱', '🃁', '🃘', '🂤', '🂷'}},
		{0, []rune{'🃃', '🃛', '🂨', '🂴', '🂢'}},
	}

	// Shuffle the hands.
	for i := range hands {
		j := rand.Intn(i + 1)
		hands[i], hands[j] = hands[j], hands[i]
	}

	for _, hand := range hands {
		// Shuffle the cards in the hand.
		for i := range hand.Cards {
			j := rand.Intn(i + 1)
			hand.Cards[i], hand.Cards[j] = hand.Cards[j], hand.Cards[i]
		}

		args = append(args, string(hand.Cards))

		out += handTypes[hand.Type] + "\n"
	}

	// Drop the trailing newline.
	out = out[:len(out)-1]

	return
}
