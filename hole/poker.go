package hole

import (
	"math/rand"
	"reflect"
	"sort"
)

func cardRune(number, suit int) rune {
	if number > 10 {
		// Skip over the unused Knight face card
		number++
	}

	return 0x1f0a1 + 16*rune(suit) + rune(number)
}

func straightCheck(numbers []int) bool {
	sort.Ints(numbers)
	if reflect.DeepEqual(numbers, []int{0, 9, 10, 11, 12}) {
		return true
	}

	return numbers[1]-numbers[0] == 1 &&
		numbers[2]-numbers[1] == 1 &&
		numbers[3]-numbers[2] == 1 &&
		numbers[4]-numbers[3] == 1
}

func poker() (args []string, out string) {
	type Hand struct {
		Type  string
		Cards []rune
	}

	var hands []Hand
	const handCount = 3

	// High card
	for i := 0; i < handCount; i++ {
		cards := rand.Perm(13)
		if straightCheck(cards[:5]) {
			i--
			continue
		}
		suits := rand.Perm(4)
		hands = append(hands, Hand{"High Card", []rune{
			cardRune(cards[0], suits[0]),
			cardRune(cards[1], suits[1]), // Avoid flush
			cardRune(cards[2], rand.Intn(4)),
			cardRune(cards[3], rand.Intn(4)),
			cardRune(cards[4], rand.Intn(4)),
		}})
	}

	// Pair
	for i := 0; i < handCount; i++ {
		cards := rand.Perm(13)
		suits := rand.Perm(4)
		hands = append(hands, Hand{"Pair", []rune{
			cardRune(cards[0], suits[0]),
			cardRune(cards[0], suits[1]),
			cardRune(cards[1], rand.Intn(4)),
			cardRune(cards[2], rand.Intn(4)),
			cardRune(cards[3], rand.Intn(4)),
		}})
	}

	// Two Pair
	for i := 0; i < handCount; i++ {
		cards := rand.Perm(13)
		suit1 := rand.Perm(4)
		suit2 := rand.Perm(4)
		hands = append(hands, Hand{"Two Pair", []rune{
			cardRune(cards[0], suit1[0]),
			cardRune(cards[0], suit1[1]),
			cardRune(cards[1], suit2[0]),
			cardRune(cards[1], suit2[1]),
			cardRune(cards[2], rand.Intn(4)),
		}})
	}

	// Three of a Kind
	for i := 0; i < handCount; i++ {
		cards := rand.Perm(13)
		suits := rand.Perm(4)
		hands = append(hands, Hand{"Three of a Kind", []rune{
			cardRune(cards[0], suits[0]),
			cardRune(cards[0], suits[1]),
			cardRune(cards[0], suits[2]),
			cardRune(cards[1], rand.Intn(4)),
			cardRune(cards[2], rand.Intn(4)),
		}})
	}

	// Four of a Kind
	for i := 0; i < handCount; i++ {
		cards := rand.Perm(13)
		hands = append(hands, Hand{"Four of a Kind", []rune{
			cardRune(cards[0], 0),
			cardRune(cards[0], 1),
			cardRune(cards[0], 2),
			cardRune(cards[0], 3),
			cardRune(cards[1], rand.Intn(4)),
		}})
	}

	// Full House
	for i := 0; i < handCount; i++ {
		cards := rand.Perm(13)
		suit1 := rand.Perm(4)
		suit2 := rand.Perm(4)
		hands = append(hands, Hand{"Full House", []rune{
			cardRune(cards[0], suit1[0]),
			cardRune(cards[0], suit1[1]),
			cardRune(cards[0], suit1[2]),
			cardRune(cards[1], suit2[0]),
			cardRune(cards[1], suit2[1]),
		}})
	}

	// Flush
	for i := 0; i < handCount; i++ {
		cards := rand.Perm(13)
		if straightCheck(cards[:5]) {
			i--
			continue
		}
		suit := rand.Intn(4)
		hands = append(hands, Hand{"Flush", []rune{
			cardRune(cards[0], suit),
			cardRune(cards[1], suit),
			cardRune(cards[2], suit),
			cardRune(cards[3], suit),
			cardRune(cards[4], suit),
		}})
	}

	// Straight
	lowCards := rand.Perm(9)
	lowCards[0] = 0 // Ensure at least one low ace
	lowCards[1] = 9 // Ensure at least one high ace
	for _, lowCard := range lowCards[:handCount] {
		suits := rand.Perm(4)
		hands = append(hands, Hand{"Straight", []rune{
			cardRune(lowCard, suits[0]),
			cardRune(lowCard+1, suits[1]), // Avoid flush
			cardRune(lowCard+2, rand.Intn(4)),
			cardRune(lowCard+3, rand.Intn(4)),
			cardRune((lowCard+4)%13, rand.Intn(4)),
		}})
	}

	// Straight Flush
	lowCards = rand.Perm(9)
	lowCards[0] = 0 // Ensure at least one low ace
	lowCards[1] = 8 // Ensure at least one 9 through king, because it could be mistaken for a royal flush.
	for _, lowCard := range lowCards[:handCount] {
		suit := rand.Intn(4)
		var hand []rune
		for card := lowCard; card < lowCard+5; card++ {
			hand = append(hand, cardRune(card%13, suit))
		}
		hands = append(hands, Hand{"Straight Flush", hand})
	}

	// Royal Flush
	for suit := 0; suit < 4; suit++ {
		var hand []rune
		for card := 9; card < 14; card++ {
			hand = append(hand, cardRune(card%13, suit))
		}
		hands = append(hands, Hand{"Royal Flush", hand})
	}

	// For a flush, the highest minus lowest codepoint is at most 13, but this is not sufficient
	// for detecting a flush. Generate hands that meet this criteria that aren't flushes.
	for suit := 0; suit < 3; suit++ {
		// Start near the top of the range on one of the lower three suits.
		// The hand will have cards with two different suits and five different face values.
		// The start card is at least 10 to avoid a straight.
		startCard := 12 - rand.Intn(3)
		// The offset of the end card from the start card must be at least 4.
		// It should be at most 11 to ensure that the range in codepoints doesn't exceed 13.
		endOffset := 4 + rand.Intn(7)
		offsets := rand.Perm(endOffset)
		hand := []rune{cardRune(startCard, suit)}
		for _, offset := range offsets[:4] {
			card := startCard + offset + 1
			hand = append(hand, cardRune(card%13, suit+card/13))
		}
		hands = append(hands, Hand{"High Card", hand})
	}

	// High Card, but could be mistaken for a straight.
	for suit := 0; suit < 3; suit++ {
		hands = append(hands, Hand{"High Card", []rune{
			cardRune(12, suit),
			cardRune(0, suit+1),
			cardRune(1, suit+1),
			cardRune(2, suit+1),
			cardRune(3, suit+1),
		}})
	}

	for _, hand := range shuffle(hands) {
		args = append(args, string(shuffle(hand.Cards)))

		out += hand.Type + "\n"
	}

	// Drop the trailing newline.
	out = out[:len(out)-1]

	return
}
