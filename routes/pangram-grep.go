package routes

import (
	"math/rand"
	"strings"
)

func pangramGrep() (args []string, out string) {
	// They all start lowercase and valid.
	pangrams := [][]byte{
		[]byte("6>_4\"gv9lb?2!ic7=-m'fd30ph].o%@w+[8unk&t1es<az(x;$^y#)q,rj\\5/*:"),
		[]byte(`a large fawn jumped quickly over white zinc boxes.`),
		[]byte(`all questions asked by five watched experts amaze the judge.`),
		[]byte(`a quick movement of the enemy will jeopardize six gunboats.`),
		[]byte(`back in june we delivered oxygen equipment of the same size.`),
		[]byte(`battle of thermopylae: quick javelin grazed wry xerxes.`),
		[]byte(`bored? craving a pub quiz fix? why, just come to the royal oak!`),
		[]byte(`bprsjzfwdqyaxgckilvunthemo`),
		[]byte(`brawny gods just flocked up to quiz and vex him.`),
		[]byte(`cute, kind, jovial, foxy physique, amazing beauty? wowser!`),
		[]byte(`fix problem quickly with galvanized jets.`),
		[]byte(`foxy parsons quiz and cajole the lovably dim wiki-girl.`),
		[]byte(`grumpy wizards make toxic brew for the evil queen and jack.`),
		[]byte(`hey zach, should i program a hex editor in java? why not sql or brainf--k!`),
		[]byte(`how razorback-jumping frogs can level six piqued gymnasts!`),
		[]byte(`jackie will budget for the most expensive zoology equipment.`),
		[]byte(`jack quietly moved up front and seized the big ball of wax.`),
		[]byte(`jim quickly realized that the beautiful gowns are expensive.`),
		[]byte(`just poets wax boldly as kings and queens march over fuzz.`),
		[]byte(`my faxed joke won a pager in the cable tv quiz show.`),
		[]byte(`quirky spud boys can jam after zapping five worthy polysixes.`),
		[]byte(`sixty zips were quickly picked from the woven jute bag.`),
		[]byte(`the quick brown fox jumps over the lazy dog.`),
		[]byte(`the wizard quickly jinxed the gnomes before they vaporized.`),
		[]byte(`when zombies arrive, quickly fax judge pat.`),
		[]byte(`"who am taking the ebonics quiz?", the prof jovially axed.`),
	}

	// Shuffle the whole set.
	for i := range pangrams {
		j := rand.Intn(i + 1)
		pangrams[i], pangrams[j] = pangrams[j], pangrams[i]
	}

	for i, pangram := range pangrams {
		clone := make([]byte, len(pangram))
		copy(clone, pangram)

		// Replace letter `i` with a different random letter.
		old := 'a' + byte(i)
		new := 'a' + byte((i+rand.Intn(25)+1)%26)
		for j, letter := range clone {
			if letter == old {
				clone[j] = new
			}
		}

		pangrams = append(pangrams, clone)
	}

	// Uppercase random letters.
	for _, pangram := range pangrams {
		for j, letter := range pangram {
			if 'a' <= letter && letter <= 'z' && rand.Intn(2) == 0 {
				pangram[j] -= 32
			}
		}
	}

	// Shuffle the whole set.
	for i := range pangrams {
		j := rand.Intn(i + 1)
		pangrams[i], pangrams[j] = pangrams[j], pangrams[i]
	}

outer:
	for _, pangram := range pangrams {
		str := string(pangram)
		args = append(args, str)

		for c := 'a'; c <= 'z'; c++ {
			if !strings.ContainsRune(str, c) && !strings.ContainsRune(str, c-32) {
				continue outer
			}
		}

		out += str + "\n"
	}

	// Drop the trailing newline.
	out = out[:len(out)-1]

	return
}
