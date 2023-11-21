package hole

import (
	"math/rand"
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
	tile := '🀀' + rune(rand.Intn(34))
	return string(tile) + string(tile)
}

func genRandomTriplet() string {
	tile := '🀀' + rune(rand.Intn(34))
	return string(tile) + string(tile) + string(tile)
}

func genRandomSequence() string {
	suit := rand.Intn(3)
	tile := '🀇' + rune(rand.Intn(7)+suit*9)
	return string(tile) + string(tile+1) + string(tile+2)
}

func genRandomStandardHand() string {
	var hand strings.Builder
	for i := 0; i < 4; i++ {
		if rand.Intn(3) > 0 {
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
	for i := 0; i < 7; i++ {
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
	duplicateTile := terminalsHonors[rand.Intn(len(terminalsHonors))]
	return string(terminalsHonors) + string(duplicateTile)
}

func genValidHand() string {
	if rand.Intn(8) > 0 {
		return genRandomStandardHand()
	} else if rand.Intn(3) > 0 {
		return genRandomSevenPairs()
	} else {
		return genRandomThirteenOrphans()
	}
}

func genInvalidHand(mutationCount int) string {
	hand := genValidHand()
	runes := []rune(hand)
	for i := 0; i < mutationCount; i++ {
		mutPos := rand.Intn(len(runes))
		mutTile := '🀀' + rune(rand.Intn(34))
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
	runs := make([]Run, 2)

	args := make([]string, 100)
	var answer strings.Builder

	for i := range args {
		hand := genValidHand()
		mutCount := rand.Intn(4)
		if mutCount > 0 {
			hand = genInvalidHand(mutCount)
		}
		runes := []rune(hand)
		rand.Shuffle(len(runes), func(i, j int) {
			runes[i], runes[j] = runes[j], runes[i]
		})
		args[i] = string(runes)

		if mutCount == 0 {
			if answer.Len() > 0 {
				answer.WriteByte('\n')
			}
			answer.WriteString(string(runes))
		}
	}

	runs[0] = Run{Args: args, Answer: answer.String()}

	// For the last run, use a set of static test cases
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
	}

	tests := append(append([]string{}, completeHands...), incompleteHands...)
	testValidity := make([]bool, len(completeHands)+len(incompleteHands))

	for i := 0; i < len(completeHands); i++ {
		testValidity[i] = true
	}

	for i := len(completeHands); i < len(testValidity); i++ {
		testValidity[i] = false
	}

	// Shuffle complete and incomplete hands
	rand.Shuffle(len(tests), func(i, j int) {
		tests[i], tests[j] = tests[j], tests[i]
		testValidity[i], testValidity[j] = testValidity[j], testValidity[i]
	})

	args2 := make([]string, len(tests))
	var answer2 strings.Builder

	for i, t := range tests {
		// Shuffle tiles within each hand
		runes := []rune(t)
		rand.Shuffle(len(runes), func(i, j int) {
			runes[i], runes[j] = runes[j], runes[i]
		})
		args2[i] = string(runes)

		// Add hand to answer, if complete
		if testValidity[i] {
			if answer2.Len() > 0 {
				answer2.WriteByte('\n')
			}
			answer2.WriteString(string(runes))
		}
	}

	runs[1] = Run{Args: args2, Answer: answer2.String()}

	return runs
}
