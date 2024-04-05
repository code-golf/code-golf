package hole

import (
	"math/rand/v2"
	"slices"
)

type cardValue int

// Off-by-one (two == 1), see cardRune() for details.
const (
	ace cardValue = iota
	two
	three
	four
	five
	six
	seven
	eight
	nine
	ten
	jack
	queen
	king
)

func cardRune[T int | cardValue](number T, suit int) rune {
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
			cardRune(cards[2], rand.IntN(4)),
			cardRune(cards[3], rand.IntN(4)),
			cardRune(cards[4], rand.IntN(4)),
		}})
	}
	suits := rand.Perm(4)
	hands = append(hands, Hand{"High Card", []rune{
		cardRune(ace, suits[0]),
		cardRune(two, suits[1]), // Avoid flush
		cardRune(three, rand.IntN(4)),
		cardRune(four, rand.IntN(4)),
		cardRune(six, rand.IntN(4)),
	}})
	suits = rand.Perm(4)
	hands = append(hands, Hand{"High Card", []rune{
		cardRune(two, rand.IntN(4)),
		cardRune(three, suits[0]),
		cardRune(four, rand.IntN(4)),
		cardRune(six, rand.IntN(4)),
		cardRune(eight, suits[1]), // Avoid flush
	}})
	suits = rand.Perm(4)
	hands = append(hands, Hand{"High Card", []rune{
		cardRune(three, suits[0]),
		cardRune(four, suits[1]), // Avoid flush
		cardRune(five, rand.IntN(4)),
		cardRune(six, rand.IntN(4)),
		cardRune(eight, rand.IntN(4)),
	}})
	suits = rand.Perm(4)
	hands = append(hands, Hand{"High Card", []rune{
		cardRune(four, suits[0]),
		cardRune(five, suits[1]), // Avoid flush
		cardRune(six, rand.IntN(4)),
		cardRune(seven, rand.IntN(4)),
		cardRune(nine, rand.IntN(4)),
	}})
	suits = rand.Perm(4)
	hands = append(hands, Hand{"High Card", []rune{
		cardRune(four, suits[0]),
		cardRune(six, suits[1]),
		cardRune(eight, suits[0]),
		cardRune(ten, suits[1]),
		cardRune(queen, rand.IntN(4)),
	}})
	suits = rand.Perm(4)
	hands = append(hands, Hand{"High Card", []rune{
		cardRune(five, suits[0]),
		cardRune(six, suits[1]),
		cardRune(seven, suits[0]),
		cardRune(nine, suits[1]),
		cardRune(jack, rand.IntN(4)),
	}})
	suits = rand.Perm(4)
	hands = append(hands, Hand{"High Card", []rune{
		cardRune(six, suits[0]),
		cardRune(nine, suits[1]),
		cardRune(ten, suits[0]),
		cardRune(jack, suits[1]),
		cardRune(queen, rand.IntN(4)),
	}})
	suits = rand.Perm(4)
	hands = append(hands, Hand{"High Card", []rune{
		cardRune(seven, suits[0]),
		cardRune(nine, suits[1]),
		cardRune(jack, suits[0]),
		cardRune(queen, suits[1]),
		cardRune(king, rand.IntN(4)),
	}})
	suits = rand.Perm(4)
	hands = append(hands, Hand{"High Card", []rune{
		cardRune(eight, suits[0]),
		cardRune(ten, suits[0]),
		cardRune(jack, suits[1]),
		cardRune(queen, rand.IntN(4)),
		cardRune(king, suits[1]),
	}})
	suits = rand.Perm(4)
	hands = append(hands, Hand{"High Card", []rune{
		cardRune(five, suits[0]),
		cardRune(six, suits[1]),
		cardRune(eight, suits[0]),
		cardRune(nine, suits[1]),
		cardRune(ten, rand.IntN(4)),
	}})
	suits = rand.Perm(4)
	hands = append(hands, Hand{"High Card", []rune{
		cardRune(six, suits[0]),
		cardRune(eight, suits[1]),
		cardRune(nine, rand.IntN(4)),
		cardRune(ten, suits[1]),
		cardRune(jack, suits[0]),
	}})
	suits = rand.Perm(4)
	hands = append(hands, Hand{"High Card", []rune{
		cardRune(five, suits[0]),
		cardRune(six, suits[1]),
		cardRune(seven, suits[0]),
		cardRune(eight, suits[1]),
		cardRune(ten, rand.IntN(4)),
	}})
	suits = rand.Perm(4)
	hands = append(hands, Hand{"High Card", []rune{
		cardRune(six, suits[0]),
		cardRune(seven, suits[0]),
		cardRune(eight, suits[1]),
		cardRune(ten, rand.IntN(4)),
		cardRune(jack, suits[1]),
	}})
	suits = rand.Perm(4)
	hands = append(hands, Hand{"High Card", []rune{
		cardRune(ace, suits[0]),
		cardRune(three, suits[0]),
		cardRune(five, suits[0]),
		cardRune(eight, suits[0]),
		cardRune(ten, suits[1]),
	}})

	// Pair
	for i := 0; i < handCount; i++ {
		cards := rand.Perm(13)
		suits := rand.Perm(4)
		hands = append(hands, Hand{"Pair", []rune{
			cardRune(cards[0], suits[0]),
			cardRune(cards[0], suits[1]),
			cardRune(cards[1], rand.IntN(4)),
			cardRune(cards[2], rand.IntN(4)),
			cardRune(cards[3], rand.IntN(4)),
		}})
	}
	suits = rand.Perm(4)
	hands = append(hands, Hand{"Pair", []rune{
		cardRune(nine, suits[0]),
		cardRune(nine, suits[1]),
		cardRune(ten, rand.IntN(4)),
		cardRune(queen, rand.IntN(4)),
		cardRune(king, rand.IntN(4)),
	}})
	suits = rand.Perm(4)
	hands = append(hands, Hand{"Pair", []rune{
		cardRune(nine, suits[0]),
		cardRune(nine, suits[1]),
		cardRune(jack, rand.IntN(4)),
		cardRune(queen, rand.IntN(4)),
		cardRune(king, rand.IntN(4)),
	}})
	suits = rand.Perm(4)
	hands = append(hands, Hand{"Pair", []rune{
		cardRune(ten, suits[0]),
		cardRune(ten, suits[1]),
		cardRune(jack, rand.IntN(4)),
		cardRune(queen, rand.IntN(4)),
		cardRune(king, rand.IntN(4)),
	}})
	suits = rand.Perm(4)
	hands = append(hands, Hand{"Pair", []rune{
		cardRune(king, suits[0]),
		cardRune(king, suits[1]),
		cardRune(jack, rand.IntN(4)),
		cardRune(ten, rand.IntN(4)),
		cardRune(nine, rand.IntN(4)),
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
			cardRune(cards[2], rand.IntN(4)),
		}})
	}
	suit1 := rand.Perm(4)
	suit2 := rand.Perm(4)
	hands = append(hands, Hand{"Two Pair", []rune{
		cardRune(king, suit1[0]),
		cardRune(king, suit1[1]),
		cardRune(nine, suit2[0]),
		cardRune(nine, suit2[1]),
		cardRune(ten, rand.IntN(4)),
	}})
	suit1 = rand.Perm(4)
	suit2 = rand.Perm(4)
	hands = append(hands, Hand{"Two Pair", []rune{
		cardRune(king, suit1[0]),
		cardRune(king, suit1[1]),
		cardRune(four, suit2[0]),
		cardRune(four, suit2[1]),
		cardRune(nine, rand.IntN(4)),
	}})
	suit1 = rand.Perm(4)
	suit2 = rand.Perm(4)
	hands = append(hands, Hand{"Two Pair", []rune{
		cardRune(queen, suit1[0]),
		cardRune(queen, suit1[1]),
		cardRune(ten, suit2[0]),
		cardRune(ten, suit2[1]),
		cardRune(jack, rand.IntN(4)),
	}})

	// Three of a Kind
	for i := 0; i < handCount; i++ {
		cards := rand.Perm(13)
		suits := rand.Perm(4)
		hands = append(hands, Hand{"Three of a Kind", []rune{
			cardRune(cards[0], suits[0]),
			cardRune(cards[0], suits[1]),
			cardRune(cards[0], suits[2]),
			cardRune(cards[1], rand.IntN(4)),
			cardRune(cards[2], rand.IntN(4)),
		}})
	}
	suits = rand.Perm(4)
	card := 1 + rand.IntN(11)
	hands = append(hands, Hand{"Three of a Kind", []rune{
		cardRune(card, suits[0]),
		cardRune(card, suits[1]),
		cardRune(card, suits[2]),
		cardRune(card+1, rand.IntN(4)),
		cardRune(card-1, rand.IntN(4)),
	}})
	suits = rand.Perm(4)
	hands = append(hands, Hand{"Three of a Kind", []rune{
		cardRune(king, suits[0]),
		cardRune(king, suits[1]),
		cardRune(king, suits[2]),
		cardRune(ace, rand.IntN(4)),
		cardRune(queen, rand.IntN(4)),
	}})
	suits = rand.Perm(4)
	hands = append(hands, Hand{"Three of a Kind", []rune{
		cardRune(ace, suits[0]),
		cardRune(ace, suits[1]),
		cardRune(ace, suits[2]),
		cardRune(king, rand.IntN(4)),
		cardRune(two, rand.IntN(4)),
	}})

	// Four of a Kind
	for i := 0; i < handCount; i++ {
		cards := rand.Perm(13)
		hands = append(hands, Hand{"Four of a Kind", []rune{
			cardRune(cards[0], 0),
			cardRune(cards[0], 1),
			cardRune(cards[0], 2),
			cardRune(cards[0], 3),
			cardRune(cards[1], rand.IntN(4)),
		}})
	}
	hands = append(hands, Hand{"Four of a Kind", []rune{
		cardRune(king, 0),
		cardRune(king, 1),
		cardRune(king, 2),
		cardRune(king, 3),
		cardRune(queen, rand.IntN(4)),
	}})
	hands = append(hands, Hand{"Four of a Kind", []rune{
		cardRune(queen, 0),
		cardRune(queen, 1),
		cardRune(queen, 2),
		cardRune(queen, 3),
		cardRune(king, rand.IntN(4)),
	}})
	hands = append(hands, Hand{"Four of a Kind", []rune{
		cardRune(king, 0),
		cardRune(king, 1),
		cardRune(king, 2),
		cardRune(king, 3),
		cardRune(ace, rand.IntN(4)),
	}})
	hands = append(hands, Hand{"Four of a Kind", []rune{
		cardRune(ace, 0),
		cardRune(ace, 1),
		cardRune(ace, 2),
		cardRune(ace, 3),
		cardRune(two, rand.IntN(4)),
	}})
	hands = append(hands, Hand{"Four of a Kind", []rune{
		cardRune(ten, 0),
		cardRune(ten, 1),
		cardRune(ten, 2),
		cardRune(ten, 3),
		cardRune(ace, rand.IntN(4)),
	}})
	hands = append(hands, Hand{"Four of a Kind", []rune{
		cardRune(ten, 0),
		cardRune(ten, 1),
		cardRune(ten, 2),
		cardRune(ten, 3),
		cardRune(two, rand.IntN(4)),
	}})
	hands = append(hands, Hand{"Four of a Kind", []rune{
		cardRune(ten, 0),
		cardRune(ten, 1),
		cardRune(ten, 2),
		cardRune(ten, 3),
		cardRune(king, rand.IntN(4)),
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
	suit1 = rand.Perm(4)
	suit2 = rand.Perm(4)
	hands = append(hands, Hand{"Full House", []rune{
		cardRune(four, suit1[0]),
		cardRune(four, suit1[1]),
		cardRune(king, suit2[0]),
		cardRune(king, suit2[1]),
		cardRune(king, suit2[2]),
	}})

	// Flush
	for i := 0; i < handCount; i++ {
		cards := rand.Perm(13)
		if straightCheck(cards[:5]) {
			i--
			continue
		}
		suit := rand.IntN(4)
		hands = append(hands, Hand{"Flush", []rune{
			cardRune(cards[0], suit),
			cardRune(cards[1], suit),
			cardRune(cards[2], suit),
			cardRune(cards[3], suit),
			cardRune(cards[4], suit),
		}})
	}

	flushHand := func(a, b, c, d, e cardValue) Hand {
		suit := rand.IntN(4)

		return Hand{"Flush", []rune{
			cardRune(a, suit),
			cardRune(b, suit),
			cardRune(c, suit),
			cardRune(d, suit),
			cardRune(e, suit),
		}}
	}

	hands = append(
		hands,
		flushHand(ace, two, three, seven, jack),
		flushHand(ace, two, four, seven, nine),
		flushHand(ace, four, five, seven, nine),
		flushHand(ace, four, five, jack, queen),
		flushHand(ace, four, seven, ten, queen),
		flushHand(ace, four, seven, ten, king),
		flushHand(ace, four, nine, ten, queen),
		flushHand(ace, five, six, eight, nine),
		flushHand(ace, five, seven, queen, king),
		flushHand(ace, seven, eight, nine, ten),
		flushHand(ace, eight, nine, ten, king),
		flushHand(ace, nine, ten, queen, king),
		flushHand(two, three, four, eight, nine),
		flushHand(two, three, five, seven, jack),
		flushHand(two, three, eight, ten, king),
		flushHand(two, three, nine, queen, king),
		flushHand(two, five, six, eight, king),
		flushHand(three, four, six, queen, king),
		flushHand(three, four, six, eight, queen),
		flushHand(three, five, seven, nine, queen),
		flushHand(three, five, nine, queen, king),
		flushHand(four, eight, ten, jack, queen),
		flushHand(five, seven, eight, nine, jack),
		flushHand(six, eight, nine, ten, jack),
	)

	// Straight
	for lowCard := range 10 {
		suits := rand.Perm(4)
		hands = append(hands, Hand{"Straight", []rune{
			cardRune(lowCard, suits[0]),
			cardRune(lowCard+1, suits[1]), // Avoid flush
			cardRune(lowCard+2, rand.IntN(4)),
			cardRune(lowCard+3, rand.IntN(4)),
			cardRune((lowCard+4)%13, rand.IntN(4)),
		}})
	}

	// Straight Flush
	for lowCard := range 9 {
		suit := rand.IntN(4)
		var hand []rune
		for card := lowCard; card < lowCard+5; card++ {
			hand = append(hand, cardRune(card%13, suit))
		}
		hands = append(hands, Hand{"Straight Flush", hand})
	}

	// Royal Flush
	for suit := range 4 {
		var hand []rune
		for card := 9; card < 14; card++ {
			hand = append(hand, cardRune(card%13, suit))
		}
		hands = append(hands, Hand{"Royal Flush", hand})
	}

	// For a flush, the highest minus lowest codepoint is at most 13, but this is not sufficient
	// for detecting a flush. Generate hands that meet this criteria that aren't flushes.
	for suit := range 3 {
		// Start near the top of the range on one of the lower three suits.
		// The hand will have cards with two different suits and five different face values.
		// The start card is at least 10 to avoid a straight.
		startCard := 12 - rand.IntN(3)
		// The offset of the end card from the start card must be at least 4.
		// It should be at most 11 to ensure that the range in codepoints doesn't exceed 13.
		endOffset := 4 + rand.IntN(7)
		offsets := rand.Perm(endOffset)
		hand := []rune{cardRune(startCard, suit)}
		for _, offset := range offsets[:4] {
			card := startCard + offset + 1
			hand = append(hand, cardRune(card%13, suit+card/13))
		}
		hands = append(hands, Hand{"High Card", hand})
	}

	// High Card, but could be mistaken for a straight.
	for suit := range 3 {
		hands = append(hands, Hand{"High Card", []rune{
			cardRune(king, suit),
			cardRune(ace, suit+1),
			cardRune(two, suit+1),
			cardRune(three, suit+1),
			cardRune(four, suit+1),
		}})
		hands = append(hands, Hand{"High Card", []rune{
			cardRune(queen, suit),
			cardRune(king, suit),
			cardRune(ace, suit+1),
			cardRune(two, suit+1),
			cardRune(three, suit+1),
		}})
		hands = append(hands, Hand{"High Card", []rune{
			cardRune(jack, suit),
			cardRune(queen, suit),
			cardRune(king, suit),
			cardRune(ace, suit+1),
			cardRune(two, suit+1),
		}})
		hands = append(hands, Hand{"High Card", []rune{
			cardRune(three, suit+1),
			cardRune(eight, suit),
			cardRune(jack, suit),
			cardRune(queen, suit),
			cardRune(king, suit),
		}})
	}

	// Flush, but could be mistaken for a straight.
	hands = append(
		hands,
		flushHand(king, ace, two, three, four),
		flushHand(queen, king, ace, two, three),
		flushHand(jack, queen, king, ace, two),
		flushHand(ace, two, three, four, seven),
		flushHand(ace, two, three, four, eight),
		flushHand(ace, two, three, eight, queen),
		flushHand(ace, two, ten, jack, queen),
		flushHand(ace, three, four, five, seven),
		flushHand(ace, three, five, seven, nine),
		flushHand(ace, four, five, six, seven),
		flushHand(ace, five, six, seven, eight),
		flushHand(ace, seven, eight, nine, ten),
		flushHand(ace, eight, nine, ten, jack),
		flushHand(ace, nine, ten, jack, queen),
		flushHand(ace, nine, ten, jack, king),
		flushHand(ace, nine, jack, queen, king),
		flushHand(two, four, six, eight, ten),
		flushHand(two, nine, jack, queen, king),
		flushHand(two, ten, jack, queen, king),
		flushHand(three, four, five, seven, eight),
		flushHand(three, five, seven, nine, jack),
		flushHand(three, ten, jack, queen, king),
		flushHand(four, five, six, seven, king),
		flushHand(four, ten, jack, queen, king),
		flushHand(five, ten, jack, queen, king),
		flushHand(six, seven, eight, jack, queen),
		flushHand(six, seven, eight, nine, king),
		flushHand(six, seven, jack, queen, king),
		flushHand(six, eight, nine, ten, queen),
		flushHand(six, ten, jack, queen, king),
		flushHand(seven, ten, jack, queen, king),
		flushHand(eight, nine, ten, jack, king),
		flushHand(eight, nine, ten, queen, king),
	)

	tests := make([]test, len(hands))
	for i, hand := range shuffle(hands) {
		tests[i] = test{string(shuffle(hand.Cards)), hand.Type}
	}

	const argc = 37 // Preserve original argc
	return outputTests(tests[:argc], tests[argc:2*argc], tests[2*argc:3*argc], tests[len(tests)-argc:])
}
