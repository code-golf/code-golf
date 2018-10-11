package routes

import "math/rand"

func pangramGrep() (args []string, out string) {
	// They all start lowercase and valid.
	pangrams := [][]byte{
		[]byte(`all questions asked by five watched experts amaze the judge.`),
		[]byte(`a quick movement of the enemy will jeopardize six gunboats.`),
		[]byte(`back in june we delivered oxygen equipment of the same size.`),
		[]byte(`battle of thermopylae: quick javelin grazed wry xerxes.`),
		[]byte(`bored? craving a pub quiz fix? why, just come to the royal oak!`),
		[]byte(`cute, kind, jovial, foxy physique, amazing beauty? wowser!`),
		[]byte(`foxy parsons quiz and cajole the lovably dim wiki-girl.`),
		[]byte(`grumpy wizards make toxic brew for the evil queen and jack.`),
		[]byte(`how razorback-jumping frogs can level six piqued gymnasts!`),
		[]byte(`jackie will budget for the most expensive zoology equipment.`),
		[]byte(`jack quietly moved up front and seized the big ball of wax.`),
		[]byte(`jim quickly realized that the beautiful gowns are expensive.`),
		[]byte(`just poets wax boldly as kings and queens march over fuzz.`),
		[]byte(`quirky spud boys can jam after zapping five worthy polysixes.`),
		[]byte(`sixty zips were quickly picked from the woven jute bag.`),
		[]byte(`the wizard quickly jinxed the gnomes before they vaporized.`),
		[]byte(`"who am taking the ebonics quiz?", the prof jovially axed.`),
	}

	// Shuffle the whole set.
	for i := range pangrams {
		j := rand.Intn(i + 1)
		pangrams[i], pangrams[j] = pangrams[j], pangrams[i]
	}

	for _, pangram := range pangrams {
		valid := true

		// Make invalid by inc-ing all occurrences of random letter.
		if rand.Intn(2) == 0 {
			// a - y
			replace := 'a' + byte(rand.Intn(25))

			for i, letter := range pangram {
				if letter == replace {
					pangram[i]++
				}
			}

			valid = false
		}

		// Uppercase random letters.
		for i, letter := range pangram {
			if 'a' <= letter && letter <= 'z' && rand.Intn(2) == 0 {
				pangram[i] -= 32
			}
		}

		str := string(pangram)
		args = append(args, str)

		if valid {
			out += str + "\n"
		}
	}

	// Drop the trailing newline.
	out = out[:len(out)-1]

	return
}
