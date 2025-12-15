package hole

import (
	"regexp"
	"slices"
	"strings"
)

var crawls []string

func init() {
	// Munge the data from "Star Wars Opening Crawl" into our format.
	// Drop first line, join the rest, lowercase, strip punctuation.
	punct := regexp.MustCompile("[.,;!?]")
	for _, t := range fixedTests("star-wars-opening-crawl") {
		crawl := strings.Join(strings.Split(t.in, "\n")[1:], " ")
		crawl = punct.ReplaceAllString(strings.ToLower(crawl), "")
		crawls = append(crawls, crawl)
	}
}

func solveSWGPT(corpusSanitized string, promptWords map[string]bool) map[string]string {
	markovFirst := map[string][]string{}
	markovCount := map[string]map[string]int{}

	words := strings.Split(corpusSanitized, " ")

	for i := range len(words) - 1 {
		if _, ok := markovCount[words[i]]; !ok {
			markovCount[words[i]] = make(map[string]int)
		}
		markovFirst[words[i]] = append(markovFirst[words[i]], words[i+1])
		markovCount[words[i]][words[i+1]]++
	}

	answer := make(map[string]string)

	for promptWord := range promptWords {
		bestWord := ""
		bestScore := 0
		for word, count := range markovCount[promptWord] {
			if count > bestScore || (count == bestScore && slices.Index(markovFirst[promptWord], word) < slices.Index(markovFirst[promptWord], bestWord)) {
				bestScore = count
				bestWord = word
			}
		}
		answer[promptWord] = bestWord
	}

	return answer
}

var _ = answerFunc("star-wars-gpt", func() []Answer {
	corpora := shuffle(slices.Clone(crawls))

	tests := make([]test, 15)
	testWordCount := 12

	for i := 0; i < 10; i += 2 {
		corpusSanitized := corpora[i] + " " + corpora[(i+1)%9]

		chosenWords := map[string]bool{}
		splitWords := strings.Split(corpusSanitized, " ")
		shuffledWords := shuffle(splitWords[:len(splitWords)-1]) // omit last word
		for j := 0; len(chosenWords) < testWordCount; j++ {
			if _, ok := chosenWords[shuffledWords[j]]; !ok {
				chosenWords[shuffledWords[j]] = true
			}
		}

		var inputSWGPT, outputSWGPT strings.Builder

		inputSWGPT.WriteString(corpusSanitized)
		inputSWGPT.WriteByte('\n')

		counter := 0
		for solveInputWord, solveOutputWord := range solveSWGPT(corpusSanitized, chosenWords) {
			inputSWGPT.WriteString(solveInputWord)
			if counter < testWordCount-1 {
				inputSWGPT.WriteByte('\n')
			}
			outputSWGPT.WriteString(solveOutputWord)
			if counter < testWordCount-1 {
				outputSWGPT.WriteByte('\n')
			}
			counter++
		}

		tests[i/2] = test{inputSWGPT.String(), outputSWGPT.String()}
	}

	staticTest1 := crawls[7] + " " + crawls[0] + "\nof\nto\nthe\ndispatched\nfirst\nmilitary\nengulfed\ngalaxy\nspark\ncertain\nevents\nseize"
	staticResult1 := "the\nthe\nfirst\ntwo\norder\ncontrol\nthe\nonly\nof\nthat\nthe\nmilitary"
	staticTest2 := crawls[3] + " " + crawls[6] + "\nenough\nin\nempire\nthe\nsave\nto\ndestroy\nand\nsinister\na\nis\nultimate"
	staticResult2 := "power\nhis\nduring\nempire's\nher\nthe\nan\nrestore\nagents\nperiod\na\nweapon"
	staticTest3 := crawls[2] + " " + crawls[6] + "\nwith\nin\nashes\ntheir\nrest\nof\nare\nthe\nknights\nfind\ncrumbling\nthere"
	staticResult3 := "their\na\nof\nvaluable\nuntil\nthe\nheroes\nrepublic\nlead\nher\nunder\nare"
	staticTest4 := crawls[1] + " " + crawls[2] + "\nthousand\nto\ncaptive\nin\nthe\nsith\njedi\ntwo\nrepublic\ngalaxy\non\nof"
	staticResult4 := "solar\nleave\nchancellor\nthe\nrepublic\nlord\nknights\njedi\nthis\nsenator\nthe\nthe"
	staticTest5 := crawls[1] + " " + crawls[6] + "\nand\ngain\nsenate\nbrother\na\npeace\nassist\nqueen\nhis\njedi\npilot\ngalactic"
	staticResult5 := "order\nhis\nseveral\nluke\nbrave\nand\nthe\nof\nabsence\nknights\non\nsenate"
	staticTest6 := crawls[1] + " " + crawls[2] + "\nchancellor\ntheir\na\nreturning\nthe\nthousand\nassist\ngalactic\nin\nqueen\non\nsolar"
	staticResult6 := "palpatine\nintentions\nstunning\nto\nrepublic\nsolar\nthe\nsenate\nthe\nof\nthe\nsystems"
	staticTest7 := crawls[3] + " " + crawls[8] + "\nintelligence\nthe\nultimate\nevil\nfor\nthreat\njedi\nwon\nempire\npower\ndead\nrebel"
	staticResult7 := "while\nempire's\nweapon\ngalactic\nbattle\nof\ntrains\ntheir\nduring\nto\nspeak\nspaceships"
	staticTest8 := crawls[6] + " " + crawls[3] + "\nthe\nstolen\nspace\nbattle\nvictory\nweapon\njustice\nto\nresistance\nin\nof\njedi"
	staticResult8 := "galaxy\nplans\nstation\nrebel\nagainst\nthe\nto\nthe\nshe\nhis\nthe\nhas"
	staticTest9 := crawls[4] + " " + crawls[3] + "\nfor\nprincess\nprobes\naboard\nestablished\nempire\ntroops\ngalactic\nthe\nhoth\nevading\nentire"
	staticResult9 := "the\nleia\ninto\nher\na\nduring\nhave\nempire\ndeath\nthe\nthe\nplanet"
	staticTest10 := crawls[2] + " " + crawls[5] + "\nswept\nmission\nas\nboth\nfreedom\nruthless\ngalactic\nnew\nskywalker\nrescue\nof\na"
	staticResult10 := "into\nto\nthe\nsides\nto\nsith\nsenate\narmored\nhas\nthe\nthe\nstunning"
	tests[5] = test{staticTest1, staticResult1}
	tests[6] = test{staticTest2, staticResult2}
	tests[7] = test{staticTest3, staticResult3}
	tests[8] = test{staticTest4, staticResult4}
	tests[9] = test{staticTest5, staticResult5}
	tests[10] = test{staticTest6, staticResult6}
	tests[11] = test{staticTest7, staticResult7}
	tests[12] = test{staticTest8, staticResult8}
	tests[13] = test{staticTest9, staticResult9}
	tests[14] = test{staticTest10, staticResult10}

	tests = shuffle(tests)

	const argc = 13 // Preserve original argc
	return outputTests(tests[:argc], tests[len(tests)-argc:])
})
