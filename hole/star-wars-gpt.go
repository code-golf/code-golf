package hole

import (
	"slices"
	"strings"
)

func solveSWGPT(corpusSanitized string, promptWords map[string]bool) map[string]string {
	markovFirst := map[string][]string{}
	markovCount := map[string]map[string]int{}

	words := strings.Split(corpusSanitized, " ")

	for i := 0; i < len(words)-1; i++ {
		_, ok := markovCount[words[i]]
		if !ok {
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
	corpora := shuffle([]string{
		`Turmoil has engulfed the Galactic Republic. The taxation of trade routes to outlying star systems is in dispute. Hoping to resolve the matter with a blockade of deadly battleships, the greedy Trade Federation has stopped all shipping to the small planet of Naboo. While the Congress of the Republic endlessly debates this alarming chain of events, the Supreme Chancellor has secretly dispatched two Jedi Knights, the guardians of peace and justice in the galaxy, to settle the conflict....`,
		`War! The Republic is crumbling under attacks by the ruthless Sith Lord, Count Dooku. There are heroes on both sides. Evil is everywhere. In a stunning move, the fiendish droid leader, General Grievous, has swept into the Republic capital and kidnapped Chancellor Palpatine, leader of the Galactic Senate. As the Separatist Droid Army attempts to flee the besieged capital with their valuable hostage, two Jedi Knights lead a desperate mission to rescue the captive Chancellor....`,
		`There is unrest in the Galactic Senate. Several thousand solar systems have declared their intentions to leave the Republic. This separatist movement, under the leadership of the mysterious Count Dooku, has made it difficult for the limited number of Jedi Knights to maintain peace and order in the galaxy. Senator Amidala, the former Queen of Naboo, is returning to the Galactic Senate to vote on the critical issue of creating an ARMY OF THE REPUBLIC to assist the overwhelmed Jedi....`,
		`It is a period of civil war. Rebel spaceships, striking from a hidden base, have won their first victory against the evil Galactic Empire. During the battle, Rebel spies managed to steal secret plans to the Empire's ultimate weapon, the DEATH STAR, an armored space station with enough power to destroy an entire planet. Pursued by the Empire's sinister agents, Princess Leia races home aboard her starship, custodian of the stolen plans that can save her people and restore freedom to the galaxy....`,
		`It is a dark time for the Rebellion. Although the Death Star has been destroyed, Imperial troops have driven the Rebel forces from their hidden base and pursued them across the galaxy. Evading the dreaded Imperial Starfleet, a group of freedom fighters led by Luke Skywalker has established a new secret base on the remote ice world of Hoth. The evil lord Darth Vader, obsessed with finding young Skywalker, has dispatched thousands of remote probes into the far reaches of space....`,
		`Luke Skywalker has returned to his home planet of Tatooine in an attempt to rescue his friend Han Solo from the clutches of the vile gangster Jabba the Hutt. Little does Luke know that the GALACTIC EMPIRE has secretly begun construction on a new armored space station even more powerful than the first dreaded Death Star. When completed, this ultimate weapon will spell certain doom for the small band of rebels struggling to restore freedom to the galaxy...`,
		`Luke Skywalker has vanished. In his absence, the sinister FIRST ORDER has risen from the ashes of the Empire and will not rest until Skywalker, the last Jedi, has been destroyed. With the support of the REPUBLIC, General Leia Organa leads a brave RESISTANCE. She is desperate to find her brother Luke and gain his help in restoring peace and justice to the galaxy. Leia has sent her most daring pilot on a secret mission to Jakku, where an old ally has discovered a clue to Luke's whereabouts....`,
		`The FIRST ORDER reigns. Having decimated the peaceful Republic, Supreme Leader Snoke now deploys his merciless legions to seize military control of the galaxy. Only General Leia Organa's band of RESISTANCE fighters stand against the rising tyranny, certain that Jedi Master Luke Skywalker will return and restore a spark of hope to the fight. But the Resistance has been exposed. As the First Order speeds toward the rebel base, the brave heroes mount a desperate escape....`,
		`The dead speak! The galaxy has heard a mysterious broadcast, a threat of REVENGE in the sinister voice of the late EMPEROR PALPATINE. GENERAL LEIA ORGANA dispatches secret agents to gather intelligence, while REY, the last hope of the Jedi, trains for battle against the diabolical FIRST ORDER. Meanwhile, Supreme Leader KYLO REN rages in search of the phantom Emperor, determined to destroy any threat to his power....`,
	})

	tests := make([]test, 11)
	testWordCount := 12

	for i := 0; i < 10; i += 2 {
		corpus := corpora[i] + " " + corpora[(i+1)%9]
		corpusSanitized := corpus
		for _, c := range ".,;!?" {
			corpusSanitized = strings.ReplaceAll(corpusSanitized, string(c), "")
		}
		corpusSanitized = strings.ToLower(corpusSanitized)

		chosenWords := map[string]bool{}
		splitWords := strings.Split(corpusSanitized, " ")
		shuffledWords := shuffle(splitWords[:len(splitWords)-1]) // omit last word
		for j := 0; len(chosenWords) < testWordCount; j++ {
			_, ok := chosenWords[shuffledWords[j]]
			if !ok {
				chosenWords[shuffledWords[j]] = true
			}
		}

		var inputSWGPT strings.Builder
		var outputSWGPT strings.Builder

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
	staticTest1 := "the first order reigns having decimated the peaceful republic supreme leader snoke now deploys his merciless legions to seize military control of the galaxy only general leia organa's band of resistance fighters stand against the rising tyranny certain that jedi master luke skywalker will return and restore a spark of hope to the fight but the resistance has been exposed as the first order speeds toward the rebel base the brave heroes mount a desperate escape turmoil has engulfed the galactic republic the taxation of trade routes to outlying star systems is in dispute hoping to resolve the matter with a blockade of deadly battleships the greedy trade federation has stopped all shipping to the small planet of naboo while the congress of the republic endlessly debates this alarming chain of events the supreme chancellor has secretly dispatched two jedi knights the guardians of peace and justice in the galaxy to settle the conflict\nof\nto\nthe\ndispatched\nfirst\nmilitary\nengulfed\ngalaxy\nspark\ncertain\nevents\nseize"
	staticResult1 := "the\nthe\nfirst\ntwo\norder\ncontrol\nthe\nonly\nof\nthat\nthe\nmilitary"
	staticTest2 := "it is a period of civil war rebel spaceships striking from a hidden base have won their first victory against the evil galactic empire during the battle rebel spies managed to steal secret plans to the empire's ultimate weapon the death star an armored space station with enough power to destroy an entire planet pursued by the empire's sinister agents princess leia races home aboard her starship custodian of the stolen plans that can save her people and restore freedom to the galaxy luke skywalker has vanished in his absence the sinister first order has risen from the ashes of the empire and will not rest until skywalker the last jedi has been destroyed with the support of the republic general leia organa leads a brave resistance she is desperate to find her brother luke and gain his help in restoring peace and justice to the galaxy leia has sent her most daring pilot on a secret mission to jakku where an old ally has discovered a clue to luke's whereabouts\nenough\nin\nempire\nthe\nsave\nto\ndestroy\nand\nsinister\na\nis\nultimate"
	staticResult2 := "power\nhis\nduring\nempire's\nher\nthe\nan\nrestore\nagents\nperiod\na\nweapon"
	staticTest3 := "war the republic is crumbling under attacks by the ruthless sith lord count dooku there are heroes on both sides evil is everywhere in a stunning move the fiendish droid leader general grievous has swept into the republic capital and kidnapped chancellor palpatine leader of the galactic senate as the separatist droid army attempts to flee the besieged capital with their valuable hostage two jedi knights lead a desperate mission to rescue the captive chancellor luke skywalker has vanished in his absence the sinister first order has risen from the ashes of the empire and will not rest until skywalker the last jedi has been destroyed with the support of the republic general leia organa leads a brave resistance she is desperate to find her brother luke and gain his help in restoring peace and justice to the galaxy leia has sent her most daring pilot on a secret mission to jakku where an old ally has discovered a clue to luke's whereabouts\nwith\nin\nashes\ntheir\nrest\nof\nare\nthe\nknights\nfind\ncrumbling\nthere"
	staticResult3 := "their\na\nof\nvaluable\nuntil\nthe\nheroes\nrepublic\nlead\nher\nunder\nare"
	staticTest4 := "there is unrest in the galactic senate several thousand solar systems have declared their intentions to leave the republic this separatist movement under the leadership of the mysterious count dooku has made it difficult for the limited number of jedi knights to maintain peace and order in the galaxy senator amidala the former queen of naboo is returning to the galactic senate to vote on the critical issue of creating an army of the republic to assist the overwhelmed jedi war the republic is crumbling under attacks by the ruthless sith lord count dooku there are heroes on both sides evil is everywhere in a stunning move the fiendish droid leader general grievous has swept into the republic capital and kidnapped chancellor palpatine leader of the galactic senate as the separatist droid army attempts to flee the besieged capital with their valuable hostage two jedi knights lead a desperate mission to rescue the captive chancellor\nthousand\nto\ncaptive\nin\nthe\nsith\njedi\ntwo\nrepublic\ngalaxy\non\nof"
	staticResult4 := "solar\nleave\nchancellor\nthe\nrepublic\nlord\nknights\njedi\nthis\nsenator\nthe\nthe"
	staticTest5 := "there is unrest in the galactic senate several thousand solar systems have declared their intentions to leave the republic this separatist movement under the leadership of the mysterious count dooku has made it difficult for the limited number of jedi knights to maintain peace and order in the galaxy senator amidala the former queen of naboo is returning to the galactic senate to vote on the critical issue of creating an army of the republic to assist the overwhelmed jedi luke skywalker has vanished in his absence the sinister first order has risen from the ashes of the empire and will not rest until skywalker the last jedi has been destroyed with the support of the republic general leia organa leads a brave resistance she is desperate to find her brother luke and gain his help in restoring peace and justice to the galaxy leia has sent her most daring pilot on a secret mission to jakku where an old ally has discovered a clue to luke's whereabouts\nand\ngain\nsenate\nbrother\na\npeace\nassist\nqueen\nhis\njedi\npilot\ngalactic"
	staticResult5 := "order\nhis\nseveral\nluke\nbrave\nand\nthe\nof\nabsence\nknights\non\nsenate"
	staticTest6 := "there is unrest in the galactic senate several thousand solar systems have declared their intentions to leave the republic this separatist movement under the leadership of the mysterious count dooku has made it difficult for the limited number of jedi knights to maintain peace and order in the galaxy senator amidala the former queen of naboo is returning to the galactic senate to vote on the critical issue of creating an army of the republic to assist the overwhelmed jedi war the republic is crumbling under attacks by the ruthless sith lord count dooku there are heroes on both sides evil is everywhere in a stunning move the fiendish droid leader general grievous has swept into the republic capital and kidnapped chancellor palpatine leader of the galactic senate as the separatist droid army attempts to flee the besieged capital with their valuable hostage two jedi knights lead a desperate mission to rescue the captive chancellor\nchancellor\ntheir\na\nreturning\nthe\nthousand\nassist\ngalactic\nin\nqueen\non\nsolar"
	staticResult6 := "palpatine\nintentions\nstunning\nto\nrepublic\nsolar\nthe\nsenate\nthe\nof\nthe\nsystems"
	tests[5] = test{staticTest1, staticResult1}
	tests[6] = test{staticTest2, staticResult2}
	tests[7] = test{staticTest3, staticResult3}
	tests[8] = test{staticTest4, staticResult4}
	tests[9] = test{staticTest5, staticResult5}
	tests[10] = test{staticTest6, staticResult6}

	return outputTests(shuffle(tests))
})
