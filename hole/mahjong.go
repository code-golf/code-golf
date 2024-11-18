package hole

import (
	"math/rand/v2"
	"strings"
)

// Gets the numerical value of a tile
func getTileNumber(tile rune) int {
	if tile >= '🀇' && tile <= '🀡' {
		return int(tile-'🀇')%9 + 1
	} else {
		return 0
	}
}

func countHandTiles(hand string) map[rune]int {
	tileCounts := make(map[rune]int)
	for _, tile := range hand {
		tileCounts[tile]++
	}
	return tileCounts
}

// Determines whether a partial hand contains the specified number of pairs and melds
func hasGroupCounts(tileCounts map[rune]int, pairCount int, meldCount int) bool {
	if pairCount < 0 || meldCount < 0 {
		return false
	}
	if pairCount == 0 && meldCount == 0 {
		return true
	}
	for tile, count := range tileCounts {
		// Check for pair
		if count >= 2 {
			tileCounts[tile] -= 2
			if hasGroupCounts(tileCounts, pairCount-1, meldCount) {
				return true
			}
			tileCounts[tile] += 2
		}
		// Check for triplet
		if count >= 3 {
			tileCounts[tile] -= 3
			if hasGroupCounts(tileCounts, pairCount, meldCount-1) {
				return true
			}
			tileCounts[tile] += 3
		}
		// Check for sequence
		if getTileNumber(tile) >= 1 && getTileNumber(tile) <= 7 {
			if tileCounts[tile] >= 1 && tileCounts[tile+1] >= 1 && tileCounts[tile+2] >= 1 {
				tileCounts[tile]--
				tileCounts[tile+1]--
				tileCounts[tile+2]--
				if hasGroupCounts(tileCounts, pairCount, meldCount-1) {
					return true
				}
				tileCounts[tile]++
				tileCounts[tile+1]++
				tileCounts[tile+2]++
			}
		}

	}
	return false
}

func genRandomPair() string {
	tile := '🀀' + rune(rand.IntN(34))
	return string(tile) + string(tile)
}

func genRandomTriplet() string {
	tile := '🀀' + rune(rand.IntN(34))
	return string(tile) + string(tile) + string(tile)
}

func genRandomSequence() string {
	suit := rand.IntN(3)
	tile := '🀇' + rune(rand.IntN(7)+suit*9)
	return string(tile) + string(tile+1) + string(tile+2)
}

func genRandomStandardHand() string {
	var hand strings.Builder
	for range 4 {
		if rand.IntN(3) > 0 {
			hand.WriteString(genRandomTriplet())
		} else {
			hand.WriteString(genRandomSequence())
		}
	}
	hand.WriteString(genRandomPair())
	res := hand.String()
	if isHandValid(countHandTiles(res)) {
		return res
	} else {
		return genRandomStandardHand()
	}
}

func genRandomSevenPairs() string {
	var hand strings.Builder
	for range 7 {
		hand.WriteString(genRandomPair())
	}
	res := hand.String()
	if isHandValid(countHandTiles(res)) {
		return res
	} else {
		return genRandomSevenPairs()
	}
}

func genRandomThirteenOrphans() string {
	terminalsHonors := []rune("🀇🀏🀐🀘🀙🀡🀀🀁🀂🀃🀄🀅🀆")
	duplicateTile := terminalsHonors[rand.IntN(len(terminalsHonors))]
	return string(terminalsHonors) + string(duplicateTile)
}

func genValidHand() string {
	if rand.IntN(8) > 0 {
		return genRandomStandardHand()
	} else if rand.IntN(3) > 0 {
		return genRandomSevenPairs()
	} else {
		return genRandomThirteenOrphans()
	}
}

func genInvalidHand(mutationCount int) string {
	hand := genValidHand()
	runes := []rune(hand)
	for range mutationCount {
		mutPos := rand.IntN(len(runes))
		mutTile := '🀀' + rune(rand.IntN(34))
		runes[mutPos] = mutTile
	}
	hand = string(runes)

	// If the hand is still valid, try again
	if isHandValid(countHandTiles(hand)) {
		return genInvalidHand(mutationCount)
	}

	// Check that this mutation did not add more than four of a single tile
	tileCounts := countHandTiles(hand)
	for _, count := range tileCounts {
		if count > 4 {
			return genInvalidHand(mutationCount)
		}
	}

	return hand
}

func isSevenPairs(tileCounts map[rune]int) bool {
	for _, count := range tileCounts {
		if count != 2 && count != 0 {
			return false
		}
	}
	return true
}

func isHandValid(tileCounts map[rune]int) bool {
	// Check if there are more than 4 of some tile
	for _, count := range tileCounts {
		if count > 4 {
			return false
		}
	}

	// Check for standard hands
	if hasGroupCounts(tileCounts, 1, 4) {
		return true
	}

	// Check for 7 pairs
	if isSevenPairs(tileCounts) {
		return true
	}

	// Check for 13 orphans
	hasPair := false
	terminalsHonors := "🀇🀏🀐🀘🀙🀡🀀🀁🀂🀃🀄🀅🀆"
	for _, tile := range terminalsHonors {
		count := tileCounts[tile]
		if count == 2 {
			hasPair = true
		}
		if count < 1 {
			return false
		}
	}
	return hasPair
}

func mahjong() []Run {
	completeHands := []string{
		"🀀🀁🀂🀃🀄🀅🀆🀆🀇🀏🀐🀘🀙🀡",
		"🀀🀀🀁🀂🀃🀄🀅🀆🀇🀏🀐🀘🀙🀡",
		"🀀🀁🀂🀃🀃🀄🀅🀆🀇🀏🀐🀘🀙🀡",
		"🀀🀁🀂🀃🀄🀅🀆🀇🀏🀐🀘🀘🀙🀡",
		"🀀🀁🀂🀃🀄🀅🀆🀇🀏🀐🀘🀙🀙🀡",
		"🀀🀁🀂🀃🀄🀅🀆🀇🀏🀏🀐🀘🀙🀡",
		"🀙🀙🀙🀛🀛🀛🀝🀝🀝🀟🀟🀟🀡🀡",
		"🀄🀄🀄🀅🀅🀅🀆🀆🀆🀇🀈🀉🀊🀊",
		"🀀🀀🀀🀁🀁🀁🀂🀂🀂🀃🀃🀙🀚🀛",
		"🀀🀀🀀🀁🀁🀁🀂🀂🀂🀃🀃🀃🀐🀐",
		"🀀🀀🀀🀁🀁🀁🀂🀂🀄🀄🀄🀅🀅🀅",
		"🀀🀀🀀🀁🀁🀁🀂🀂🀂🀃🀃🀃🀄🀄",
		"🀀🀀🀁🀁🀂🀂🀃🀃🀄🀄🀅🀅🀆🀆",
		"🀇🀇🀇🀏🀏🀏🀐🀐🀙🀙🀙🀡🀡🀡",
		"🀅🀅🀅🀑🀑🀒🀒🀓🀓🀕🀕🀗🀗🀗",
		"🀇🀇🀇🀇🀈🀉🀊🀋🀌🀍🀎🀏🀏🀏",
		"🀙🀙🀙🀚🀚🀛🀜🀝🀞🀟🀠🀡🀡🀡",
		"🀐🀐🀐🀑🀒🀒🀓🀔🀕🀖🀗🀘🀘🀘",
		"🀇🀇🀈🀈🀉🀉🀉🀉🀊🀊🀋🀋🀌🀌",
		"🀑🀑🀒🀒🀓🀓🀓🀓🀔🀔🀕🀕🀗🀗",
		"🀐🀐🀐🀐🀑🀑🀒🀒🀜🀜🀝🀝🀞🀞",
		"🀇🀇🀇🀉🀉🀉🀋🀋🀋🀍🀍🀖🀗🀘",
		"🀀🀀🀋🀋🀌🀌🀍🀍🀜🀜🀝🀝🀞🀞",
	}

	incompleteHands := []string{
		"🀀🀁🀂🀃🀄🀅🀆🀇🀏🀐🀘🀙🀝🀡",
		"🀀🀁🀂🀃🀄🀅🀆🀇🀏🀐🀒🀘🀙🀡",
		"🀀🀁🀂🀃🀄🀅🀆🀇🀍🀏🀐🀘🀙🀡",
		"🀀🀁🀂🀃🀄🀅🀆🀇🀏🀐🀘🀙🀛🀛",
		"🀀🀂🀃🀄🀇🀊🀎🀐🀒🀕🀚🀚🀝🀡",
		"🀀🀁🀂🀂🀂🀃🀄🀅🀆🀇🀏🀐🀘🀙",
		"🀀🀂🀃🀄🀅🀆🀇🀇🀇🀏🀐🀘🀙🀡",
		"🀀🀁🀁🀁🀂🀃🀄🀅🀇🀏🀐🀘🀙🀡",
		"🀀🀁🀂🀇🀈🀉🀊🀋🀌🀍🀎🀏🀐🀐",
		"🀄🀅🀆🀇🀈🀉🀊🀋🀌🀍🀎🀏🀐🀐",
		"🀃🀄🀅🀇🀈🀉🀊🀋🀌🀍🀎🀏🀐🀐",
		"🀏🀐🀑🀙🀙🀙🀛🀛🀛🀝🀝🀝🀟🀟",
		"🀎🀏🀐🀙🀙🀙🀛🀛🀛🀝🀝🀝🀟🀟",
		"🀇🀇🀇🀉🀉🀉🀋🀋🀋🀍🀍🀘🀙🀚",
		"🀇🀇🀇🀉🀉🀉🀋🀋🀋🀍🀍🀗🀘🀙",
		"🀀🀇🀇🀇🀉🀉🀉🀋🀋🀋🀍🀍🀠🀡",
		"🀇🀈🀐🀐🀐🀒🀒🀒🀔🀔🀔🀖🀖🀡",
		"🀂🀂🀃🀃🀇🀇🀉🀉🀋🀋🀕🀕🀕🀕",
		"🀀🀀🀁🀁🀁🀁🀂🀂🀃🀃🀄🀄🀅🀅",
		"🀐🀐🀐🀐🀑🀑🀒🀒🀜🀜🀝🀝🀟🀟",
		"🀀🀍🀍🀍🀐🀐🀚🀚🀜🀜🀞🀞🀠🀠",
		"🀀🀀🀁🀁🀂🀂🀃🀃🀄🀄🀅🀅🀅🀆",
		"🀀🀀🀀🀁🀁🀂🀂🀃🀃🀃🀄🀄🀄🀄",
		"🀐🀐🀐🀐🀑🀑🀒🀓🀔🀕🀖🀗🀘🀘",
		"🀐🀐🀐🀑🀒🀓🀔🀔🀕🀕🀗🀘🀘🀘",
		"🀆🀆🀆🀖🀘🀛🀛🀛🀞🀞🀞🀟🀠🀡",
		"🀉🀊🀋🀋🀌🀌🀍🀎🀏🀏🀏🀙🀙🀙",
		"🀘🀕🀌🀕🀍🀕🀘🀎🀀🀌🀌🀘🀍🀀",
		"🀗🀇🀉🀈🀞🀋🀞🀊🀇🀉🀗🀉🀞🀗",
		"🀌🀋🀔🀌🀔🀊🀍🀋🀊🀊🀙🀋🀙🀙",
		"🀃🀌🀎🀄🀃🀄🀍🀒🀑🀍🀄🀐🀍🀓",
		"🀎🀏🀎🀁🀎🀏🀏🀍🀛🀍🀎🀛🀛🀁",
		"🀋🀞🀜🀜🀌🀟🀏🀍🀝🀍🀎🀜🀊🀎",
		"🀞🀄🀄🀞🀒🀓🀑🀌🀔🀄🀟🀞🀌🀟",
		"🀛🀎🀍🀜🀋🀌🀜🀛🀐🀛🀖🀖🀐🀖",
		"🀔🀋🀓🀊🀠🀑🀒🀡🀋🀑🀟🀉🀊🀊",
		"🀊🀊🀔🀠🀞🀔🀟🀟🀂🀡🀂🀂🀔🀊",
		"🀉🀛🀠🀘🀝🀉🀜🀕🀉🀕🀗🀠🀗🀖",
		"🀓🀔🀌🀌🀝🀚🀛🀌🀓🀔🀜🀞🀔🀛",
		"🀌🀈🀠🀠🀈🀉🀍🀇🀟🀋🀇🀠🀈🀟",
		"🀑🀟🀒🀡🀔🀠🀆🀔🀟🀟🀓🀆🀆🀞",
		"🀋🀇🀑🀖🀕🀕🀔🀓🀑🀉🀋🀕🀑🀈",
		"🀕🀚🀈🀗🀛🀋🀉🀛🀊🀌🀉🀜🀛🀖",
		"🀏🀓🀕🀕🀘🀘🀙🀙🀜🀜🀠🀠🀡🀡",
		"🀁🀂🀃🀄🀅🀆🀇🀇🀏🀐🀘🀙🀙🀡",
		"🀀🀀🀀🀈🀈🀐🀐🀘🀘🀙🀙🀡🀡🀡",
		"🀀🀀🀁🀂🀃🀄🀅🀆🀆🀇🀏🀐🀐🀘",
		"🀀🀀🀁🀂🀃🀄🀅🀆🀎🀏🀐🀘🀙🀡",
		"🀀🀀🀁🀂🀃🀄🀅🀍🀎🀏🀐🀘🀙🀡",
		"🀀🀀🀁🀂🀃🀄🀆🀆🀇🀏🀐🀘🀙🀡",
		"🀀🀀🀁🀂🀃🀆🀇🀏🀏🀐🀐🀙🀚🀡",
		"🀀🀀🀈🀈🀐🀐🀘🀘🀘🀘🀙🀙🀡🀡",
		"🀀🀁🀁🀂🀃🀄🀅🀆🀇🀇🀏🀐🀘🀙",
		"🀀🀁🀁🀂🀃🀅🀆🀇🀏🀐🀘🀙🀚🀡",
		"🀀🀁🀂🀂🀃🀄🀆🀆🀇🀏🀐🀘🀙🀡",
		"🀀🀁🀂🀃🀃🀄🀅🀆🀇🀏🀏🀐🀘🀙",
		"🀀🀁🀂🀃🀄🀄🀆🀇🀏🀐🀘🀘🀙🀡",
		"🀀🀁🀂🀃🀄🀅🀅🀆🀇🀐🀘🀙🀡🀡",
		"🀀🀁🀂🀃🀄🀅🀆🀇🀇🀏🀏🀐🀘🀡",
		"🀀🀁🀂🀃🀄🀅🀆🀎🀏🀐🀘🀙🀡🀡",
		"🀀🀁🀂🀃🀄🀅🀍🀎🀏🀐🀘🀙🀡🀡",
		"🀀🀁🀂🀃🀄🀌🀍🀎🀏🀐🀘🀙🀡🀡",
		"🀀🀁🀂🀃🀋🀌🀍🀎🀏🀐🀘🀙🀡🀡",
		"🀀🀁🀂🀊🀋🀌🀍🀎🀏🀐🀘🀙🀡🀡",
		"🀀🀁🀃🀄🀅🀆🀇🀏🀏🀐🀘🀙🀙🀡",
		"🀀🀁🀉🀊🀋🀌🀍🀎🀏🀐🀘🀙🀡🀡",
		"🀀🀈🀉🀊🀋🀌🀍🀎🀏🀐🀘🀙🀡🀡",
		"🀀🀈🀉🀊🀋🀌🀍🀎🀏🀗🀘🀙🀡🀡",
		"🀀🀈🀉🀊🀋🀌🀍🀎🀖🀗🀘🀙🀡🀡",
		"🀠🀆🀐🀚🀏🀁🀀🀇🀘🀄🀃🀅🀚🀂",
		"🀉🀙🀇🀖🀚🀇🀔🀕🀇🀛🀏🀈🀟🀞",
		"🀀🀀🀂🀃🀄🀅🀆🀇🀏🀐🀘🀙🀡🀞",
		"🀂🀀🀁🀂🀄🀅🀆🀇🀏🀐🀘🀙🀡🀎",
		"🀃🀀🀁🀂🀃🀄🀅🀇🀏🀐🀘🀙🀡🀔",
		"🀙🀀🀁🀂🀃🀄🀅🀆🀇🀏🀘🀙🀡🀛",
		"🀅🀀🀁🀂🀃🀄🀅🀆🀇🀏🀐🀙🀡🀐",
		"🀆🀄🀀🀂🀃🀏🀁🀡🀡🀡🀐🀙🀘🀄",
		"🀆🀘🀆🀄🀁🀙🀏🀇🀆🀐🀏🀃🀡🀂",
		"🀘🀅🀐🀂🀆🀃🀃🀀🀂🀙🀄🀂🀇🀏",
		"🀀🀀🀇🀈🀈🀈🀉🀊🀊🀌🀌🀐🀐🀐", // Integer wrapping exploit (1, 3, 1, 2, 0, 2)
	}

	const argc = 100 // Preserve original argc

	var tests []test

	// Start with some hardcoded complete hands.
	for _, hand := range completeHands {
		hand = string(shuffle([]rune(hand)))
		tests = append(tests, test{in: hand, out: hand})
	}

	// Append some hardcoded incomplete hands.
	for _, hand := range incompleteHands {
		hand = string(shuffle([]rune(hand)))
		tests = append(tests, test{in: hand})
	}

	// Fill-in the remaining with random hands.
	for range 2*argc - len(tests) {
		hand := genValidHand()
		mutCount := rand.IntN(4)
		if mutCount > 0 {
			hand = genInvalidHand(mutCount)
		}

		test := test{in: string(shuffle([]rune(hand)))}
		if mutCount == 0 {
			test.out = test.in
		}

		tests = append(tests, test)
	}

	shuffle(tests)
	return outputTests(tests[:argc], tests[argc:])
}
