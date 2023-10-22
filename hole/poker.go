package hole

import (
	"math/rand"
	"slices"
)

type CardValue int

const (
	ACE   CardValue = 0
	TWO   CardValue = 1
	THREE CardValue = 2
	FOUR  CardValue = 3
	FIVE  CardValue = 4
	SIX   CardValue = 5
	SEVEN CardValue = 6
	EIGHT CardValue = 7
	NINE  CardValue = 8
	TEN   CardValue = 9
	JACK  CardValue = 10
	QUEEN CardValue = 11
	KING  CardValue = 12
)

func cardRune[T int | CardValue](num T, suit int) rune {
	number := int(num)

	if number > 10 {
		// Skip over the unused Knight face card
		number++
	}

	return 0x1f0a1 + 16*rune(suit) + rune(number)
}

func straightCheck(numbers []int) bool {
	slices.Sort(numbers)

	// After sorting we have an Ace-straight or the numbers are sequential.
	return slices.Equal(numbers, []int{0, 9, 10, 11, 12}) ||
		(numbers[1]-numbers[0] == 1 &&
			numbers[2]-numbers[1] == 1 &&
			numbers[3]-numbers[2] == 1 &&
			numbers[4]-numbers[3] == 1)
}

func poker() []Run {
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
	suits := rand.Perm(4)
	hands = append(hands, Hand{"High Card", []rune{
		cardRune(ACE, suits[0]),
		cardRune(TWO, suits[1]), // Avoid flush
		cardRune(THREE, rand.Intn(4)),
		cardRune(FOUR, rand.Intn(4)),
		cardRune(SIX, rand.Intn(4)),
	}})
	suits = rand.Perm(4)
	hands = append(hands, Hand{"High Card", []rune{
		cardRune(TWO, rand.Intn(4)),
		cardRune(THREE, suits[0]),
		cardRune(FOUR, rand.Intn(4)),
		cardRune(SIX, rand.Intn(4)),
		cardRune(EIGHT, suits[1]), // Avoid flush
	}})
	suits = rand.Perm(4)
	hands = append(hands, Hand{"High Card", []rune{
		cardRune(THREE, suits[0]),
		cardRune(FOUR, suits[1]), // Avoid flush
		cardRune(FIVE, rand.Intn(4)),
		cardRune(SIX, rand.Intn(4)),
		cardRune(EIGHT, rand.Intn(4)),
	}})
	suits = rand.Perm(4)
	hands = append(hands, Hand{"High Card", []rune{
		cardRune(FOUR, suits[0]),
		cardRune(FIVE, suits[1]), // Avoid flush
		cardRune(SIX, rand.Intn(4)),
		cardRune(SEVEN, rand.Intn(4)),
		cardRune(NINE, rand.Intn(4)),
	}})
	suits = rand.Perm(4)
	hands = append(hands, Hand{"High Card", []rune{
		cardRune(FOUR, suits[0]),
		cardRune(SIX, suits[1]),
		cardRune(EIGHT, suits[0]),
		cardRune(TEN, suits[1]),
		cardRune(QUEEN, rand.Intn(4)),
	}})
	suits = rand.Perm(4)
	hands = append(hands, Hand{"High Card", []rune{
		cardRune(FIVE, suits[0]),
		cardRune(SIX, suits[1]),
		cardRune(SEVEN, suits[0]),
		cardRune(NINE, suits[1]),
		cardRune(JACK, rand.Intn(4)),
	}})
	suits = rand.Perm(4)
	hands = append(hands, Hand{"High Card", []rune{
		cardRune(SIX, suits[0]),
		cardRune(NINE, suits[1]),
		cardRune(TEN, suits[0]),
		cardRune(JACK, suits[1]),
		cardRune(QUEEN, rand.Intn(4)),
	}})
	suits = rand.Perm(4)
	hands = append(hands, Hand{"High Card", []rune{
		cardRune(SEVEN, suits[0]),
		cardRune(NINE, suits[1]),
		cardRune(JACK, suits[0]),
		cardRune(QUEEN, suits[1]),
		cardRune(KING, rand.Intn(4)),
	}})
	suits = rand.Perm(4)
	hands = append(hands, Hand{"High Card", []rune{
		cardRune(EIGHT, suits[0]),
		cardRune(TEN, suits[0]),
		cardRune(JACK, suits[1]),
		cardRune(QUEEN, rand.Intn(4)),
		cardRune(KING, suits[1]),
	}})

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
	suits = rand.Perm(4)
	hands = append(hands, Hand{"Pair", []rune{
		cardRune(NINE, suits[0]),
		cardRune(NINE, suits[1]),
		cardRune(TEN, rand.Intn(4)),
		cardRune(QUEEN, rand.Intn(4)),
		cardRune(KING, rand.Intn(4)),
	}})
	suits = rand.Perm(4)
	hands = append(hands, Hand{"Pair", []rune{
		cardRune(NINE, suits[0]),
		cardRune(NINE, suits[1]),
		cardRune(JACK, rand.Intn(4)),
		cardRune(QUEEN, rand.Intn(4)),
		cardRune(KING, rand.Intn(4)),
	}})
	suits = rand.Perm(4)
	hands = append(hands, Hand{"Pair", []rune{
		cardRune(KING, suits[0]),
		cardRune(KING, suits[1]),
		cardRune(JACK, rand.Intn(4)),
		cardRune(TEN, rand.Intn(4)),
		cardRune(NINE, rand.Intn(4)),
	}})

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
	suit1 := rand.Perm(4)
	suit2 := rand.Perm(4)
	hands = append(hands, Hand{"Two Pair", []rune{
		cardRune(KING, suit1[0]),
		cardRune(KING, suit1[1]),
		cardRune(NINE, suit2[0]),
		cardRune(NINE, suit2[1]),
		cardRune(TEN, rand.Intn(4)),
	}})
	suit1 = rand.Perm(4)
	suit2 = rand.Perm(4)
	hands = append(hands, Hand{"Two Pair", []rune{
		cardRune(KING, suit1[0]),
		cardRune(KING, suit1[1]),
		cardRune(FOUR, suit2[0]),
		cardRune(FOUR, suit2[1]),
		cardRune(NINE, rand.Intn(4)),
	}})
	suit1 = rand.Perm(4)
	suit2 = rand.Perm(4)
	hands = append(hands, Hand{"Two Pair", []rune{
		cardRune(QUEEN, suit1[0]),
		cardRune(QUEEN, suit1[1]),
		cardRune(TEN, suit2[0]),
		cardRune(TEN, suit2[1]),
		cardRune(JACK, rand.Intn(4)),
	}})

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
	suits = rand.Perm(4)
	card := 1 + rand.Intn(11)
	hands = append(hands, Hand{"Three of a Kind", []rune{
		cardRune(card, suits[0]),
		cardRune(card, suits[1]),
		cardRune(card, suits[2]),
		cardRune(card+1, rand.Intn(4)),
		cardRune(card-1, rand.Intn(4)),
	}})
	suits = rand.Perm(4)
	hands = append(hands, Hand{"Three of a Kind", []rune{
		cardRune(KING, suits[0]),
		cardRune(KING, suits[1]),
		cardRune(KING, suits[2]),
		cardRune(ACE, rand.Intn(4)),
		cardRune(QUEEN, rand.Intn(4)),
	}})
	suits = rand.Perm(4)
	hands = append(hands, Hand{"Three of a Kind", []rune{
		cardRune(ACE, suits[0]),
		cardRune(ACE, suits[1]),
		cardRune(ACE, suits[2]),
		cardRune(KING, rand.Intn(4)),
		cardRune(TWO, rand.Intn(4)),
	}})

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
	hands = append(hands, Hand{"Four of a Kind", []rune{
		cardRune(KING, 0),
		cardRune(KING, 1),
		cardRune(KING, 2),
		cardRune(KING, 3),
		cardRune(QUEEN, rand.Intn(4)),
	}})
	hands = append(hands, Hand{"Four of a Kind", []rune{
		cardRune(QUEEN, 0),
		cardRune(QUEEN, 1),
		cardRune(QUEEN, 2),
		cardRune(QUEEN, 3),
		cardRune(KING, rand.Intn(4)),
	}})
	hands = append(hands, Hand{"Four of a Kind", []rune{
		cardRune(KING, 0),
		cardRune(KING, 1),
		cardRune(KING, 2),
		cardRune(KING, 3),
		cardRune(ACE, rand.Intn(4)),
	}})
	hands = append(hands, Hand{"Four of a Kind", []rune{
		cardRune(ACE, 0),
		cardRune(ACE, 1),
		cardRune(ACE, 2),
		cardRune(ACE, 3),
		cardRune(TWO, rand.Intn(4)),
	}})
	hands = append(hands, Hand{"Four of a Kind", []rune{
		cardRune(TEN, 0),
		cardRune(TEN, 1),
		cardRune(TEN, 2),
		cardRune(TEN, 3),
		cardRune(ACE, rand.Intn(4)),
	}})
	hands = append(hands, Hand{"Four of a Kind", []rune{
		cardRune(TEN, 0),
		cardRune(TEN, 1),
		cardRune(TEN, 2),
		cardRune(TEN, 3),
		cardRune(TWO, rand.Intn(4)),
	}})
	hands = append(hands, Hand{"Four of a Kind", []rune{
		cardRune(TEN, 0),
		cardRune(TEN, 1),
		cardRune(TEN, 2),
		cardRune(TEN, 3),
		cardRune(KING, rand.Intn(4)),
	}})

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

	suits = rand.Perm(4)
	hands = append(hands, Hand{"Straight", []rune{
		cardRune(TWO, suits[0]), // Avoid flush
		cardRune(THREE, suits[1]),
		cardRune(FOUR, rand.Intn(4)),
		cardRune(FIVE, rand.Intn(4)),
		cardRune(SIX, rand.Intn(4)),
	}})
	suits = rand.Perm(4)
	hands = append(hands, Hand{"Straight", []rune{
		cardRune(FIVE, suits[0]), // Avoid flush
		cardRune(SIX, suits[1]),
		cardRune(SEVEN, rand.Intn(4)),
		cardRune(EIGHT, rand.Intn(4)),
		cardRune(NINE, rand.Intn(4)),
	}})
	suits = rand.Perm(4)
	hands = append(hands, Hand{"Straight", []rune{
		cardRune(EIGHT, rand.Intn(4)),
		cardRune(NINE, suits[0]),
		cardRune(TEN, rand.Intn(4)),
		cardRune(JACK, rand.Intn(4)),
		cardRune(QUEEN, suits[1]), // Avoid flush
	}})

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
			cardRune(KING, suit),
			cardRune(ACE, suit+1),
			cardRune(TWO, suit+1),
			cardRune(THREE, suit+1),
			cardRune(FOUR, suit+1),
		}})
		hands = append(hands, Hand{"High Card", []rune{
			cardRune(QUEEN, suit),
			cardRune(KING, suit),
			cardRune(ACE, suit+1),
			cardRune(TWO, suit+1),
			cardRune(THREE, suit+1),
		}})
		hands = append(hands, Hand{"High Card", []rune{
			cardRune(JACK, suit),
			cardRune(QUEEN, suit),
			cardRune(KING, suit),
			cardRune(ACE, suit+1),
			cardRune(TWO, suit+1),
		}})
	}

	// Flush, but could be mistaken for a straight.
	suit := rand.Intn(4)
	hands = append(hands, Hand{"Flush", []rune{
		cardRune(KING, suit),
		cardRune(ACE, suit),
		cardRune(TWO, suit),
		cardRune(THREE, suit),
		cardRune(FOUR, suit),
	}})
	suit = rand.Intn(4)
	hands = append(hands, Hand{"Flush", []rune{
		cardRune(QUEEN, suit),
		cardRune(KING, suit),
		cardRune(ACE, suit),
		cardRune(TWO, suit),
		cardRune(THREE, suit),
	}})
	suit = rand.Intn(4)
	hands = append(hands, Hand{"Flush", []rune{
		cardRune(JACK, suit),
		cardRune(QUEEN, suit),
		cardRune(KING, suit),
		cardRune(ACE, suit),
		cardRune(TWO, suit),
	}})
	suit = rand.Intn(4)
	hands = append(hands, Hand{"Flush", []rune{
		cardRune(ACE, suit),
		cardRune(SEVEN, suit),
		cardRune(EIGHT, suit),
		cardRune(NINE, suit),
		cardRune(TEN, suit),
	}})
	suit = rand.Intn(4)
	hands = append(hands, Hand{"Flush", []rune{
		cardRune(ACE, suit),
		cardRune(EIGHT, suit),
		cardRune(NINE, suit),
		cardRune(TEN, suit),
		cardRune(JACK, suit),
	}})
	suit = rand.Intn(4)
	hands = append(hands, Hand{"Flush", []rune{
		cardRune(ACE, suit),
		cardRune(NINE, suit),
		cardRune(TEN, suit),
		cardRune(JACK, suit),
		cardRune(QUEEN, suit),
	}})

	tests := make([]test, len(hands))
	for i, hand := range shuffle(hands) {
		tests[i] = test{string(shuffle(hand.Cards)), hand.Type}
	}

	const argc = 37 // Preserve original argc
	return outputTests(tests[:argc], tests[argc:2*argc], tests[len(tests)-argc:])
}
