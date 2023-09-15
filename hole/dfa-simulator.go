package hole

import (
	"math/rand"
	"strconv"
	"strings"
)

type DFAState struct {
	name     string
	accepted bool
}

type DFAStateTransition struct {
	input     string
	stateName string
}

func solve(g string) string {
	var startState DFAState
	stateMap := map[DFAStateTransition]DFAState{}
	nameToState := map[string]DFAState{}

	data := strings.Split(g, "\n")
	alphabet := strings.Split(data[0][4:], " ")

	stateRows := data[1 : len(data)-1]

	for _, state := range stateRows {
		stateName := string(state[2])
		newState := DFAState{stateName, state[1] == 'F'}
		nameToState[stateName] = newState

		if state[0] == '>' {
			startState = newState
		}
	}

	for _, state := range stateRows {
		transitions := strings.Split(state[2:], " ")
		stateName := transitions[0]

		for i, transition := range transitions[1:] {
			alphabetWord := string(alphabet[i])
			newStateTransition := DFAStateTransition{alphabetWord, stateName}
			stateMap[newStateTransition] = nameToState[transition]
		}
	}

	currentState := startState
	for _, testString := range data[len(data)-1] {
		if string(testString) != "\"" {
			currentState = stateMap[DFAStateTransition{string(testString), currentState.name}]
		}
	}
	if currentState.accepted {
		return currentState.name + " Accept"
	}

	return currentState.name + " Reject"
}

func generate() string {
	alphabet := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	alphabetLength := rand.Intn(36)
	alphabet = alphabet[:alphabetLength]

	inputDFA := "    " + strings.Join(alphabet, " ") + "\n"

	stateLength := rand.Intn(10) + 1
	startState := rand.Intn(stateLength)
	for i := 0; i < stateLength; i++ {
		stateRow := ""

		if i == startState {
			stateRow += ">"
		} else {
			stateRow += " "
		}

		if rand.Intn(2) == 0 {
			stateRow += "F"
		} else {
			stateRow += " "
		}

		stateRow += strconv.Itoa(i) + " "

		transitions := []string{}
		for j := 0; j < alphabetLength; j++ {
			transitions = append(transitions, strconv.Itoa(rand.Intn(stateLength)))
		}
		stateRow += strings.Join(transitions, " ") + "\n"
		inputDFA += stateRow

	}
	inputDFA += "\""
	for i := 0; i < rand.Intn(2*alphabetLength); i++ {
		inputDFA += alphabet[rand.Intn(alphabetLength)]
	}
	inputDFA += "\""

	return inputDFA
}

func dfaSimulator() []Run {
	args := []string{
		"    a\n> 0 1\n F1 2\n F2 0\n\"aaaaaaa\"",
		"    a\n> 0 1\n F1 2\n F2 0\n\"aaaaaaaa\"",
		"    a\n> 0 1\n F1 2\n F2 0\n\"aaa\"",
		"    a\n> 0 1\n F1 2\n F2 0\n\"\"",
		"    a b c d r\n> 0 0 0 0 0 0\n\"abracadabra\"",
		"    a b c d r\n>F0 0 0 0 0 0\n\"abracadabra\"",
		"    a b c d r\n>F9 9 9 9 9 9\n\"abracadabra\"",
		"    a b c d r\n> 0 0 0 0 0 1\n  1 2 0 0 0 1\n F2 0 0 0 0 1\n\"abracadabra\"",
		"    a b c d r\n> 0 0 0 0 0 1\n  1 2 0 0 0 1\n F2 0 0 0 0 1\n\"barra\"",
		"    a b c d e f\n> 0 0 0 0 1 0 0\n  1 0 0 0 1 0 2\n  2 3 0 0 1 0 0\n F3 3 3 3 3 3 3\n\"dfa\"",
		"    a b c d e f\n> 0 0 0 0 1 0 0\n  1 0 0 0 1 0 2\n  2 3 0 0 1 0 0\n F3 3 3 3 3 3 3\n\"aabacadafad\"",
		"    a b c d e f\n> 0 0 0 0 1 0 0\n  1 0 0 0 1 0 2\n  2 3 0 0 1 0 0\n F3 3 3 3 3 3 3\n\"aabacadfad\"",
		"    a b c d e f g h i j k l m n o p q r s t u v w x y z\n>F0 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0\n  1 0 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1\n\"bananajazz\"",
		"    a b c d e f g h i j k l m n o p q r s t u v w x y z\n>F0 0 0 0 0 0 0 0 0 0 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 1\n  1 1 1 1 1 1 1 1 1 1 0 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 0\n\"bananajazz\"",
		"    a b c d e f g h i j k l m n o p q r s t u v w x y z\n> 0 0 0 0 0 0 0 0 0 0 0 0 0 0 1 0 0 0 0 0 0 0 0 0 0 0 1\n F1 1 1 1 1 1 1 1 1 1 1 1 1 1 0 1 1 1 1 1 1 1 1 1 1 1 0\n\"panamajazz\"",
		"    a b c d e f g h i j k l m n o p q r s t u v w x y z\n> 0 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 1\n F1 0 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 0\n\"panamajazz\"",
		"    0 1 2 3 4 5 6 7 8 9 a b c d e f\n F4 4 4 4 4 4 4 4 4 4 4 4 4 4 4 4 4\n  3 0 0 0 0 0 0 0 0 0 0 0 0 1 0 4 0\n  2 0 0 0 0 0 0 0 0 0 0 0 0 1 3 0 0\n  1 2 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0\n> 0 0 0 0 0 0 0 0 0 0 0 0 0 1 0 0 0\n\"123456789c0de801f\"",
		"    0 1 2 3 4 5 6 7 8 9 a b c d e f\n F4 4 4 4 4 4 4 4 4 4 4 4 4 4 4 4 4\n  3 0 0 0 0 0 0 0 0 0 0 0 0 1 0 4 0\n  2 0 0 0 0 0 0 0 0 0 0 0 0 1 3 0 0\n  1 2 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0\n> 0 0 0 0 0 0 0 0 0 0 0 0 0 1 0 0 0\n\"123456789c0d3801f\"",
		"    0 1 2 3 4 5 6 7 8 9 a b c d e f\n F4 4 4 4 4 4 4 4 4 4 4 4 4 4 4 4 4\n  3 0 0 0 0 0 0 0 0 0 0 0 0 1 0 4 0\n  2 0 0 0 0 0 0 0 0 0 0 0 0 1 3 0 0\n  1 2 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0\n> 0 0 0 0 0 0 0 0 0 0 0 0 0 1 0 0 0\n\"123456789801f2c0d\"",
		"    0 1 2 3 4 5 6 7 8 9 a b c d e f\n F4 4 4 4 4 4 4 4 4 4 4 4 4 4 4 4 4\n  3 3 3 3 3 3 3 3 3 3 3 3 3 3 3 4 3\n  2 2 2 2 2 2 2 2 2 2 2 2 2 2 3 2 2\n  1 2 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1\n> 0 0 0 0 0 0 0 0 0 0 0 0 0 1 0 0 0\n\"123456789c0d3801f\"",
		"    0 1 2 3 4 5 6 7 8 9 a b c d e f\n F4 4 4 4 4 4 4 4 4 4 4 4 4 4 4 4 4\n  3 3 3 3 3 3 3 3 3 3 3 3 3 3 3 4 3\n  2 2 2 2 2 2 2 2 2 2 2 2 2 2 3 2 2\n  1 2 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1\n> 0 0 0 0 0 0 0 0 0 0 0 0 0 1 0 0 0\n\"123456789c0d38e1f\"",
		"    0 1\n> 0 2 1\n  1 3 1\n  2 2 2\n F3 3 1\n\"01010101\"",
		"    0 1\n> 0 2 1\n  1 3 1\n F2 2 2\n  3 3 1\n\"01010101\"",
		"    0 1\n> 0 2 1\n  1 3 1\n  2 2 2\n F3 3 1\n\"10101010\"",
		"    0 1\n> 0 2 1\n  1 3 1\n F2 2 2\n  3 3 1\n\"10101010\"",
		"    0 1\n> 0 2 1\n  1 3 1\n  2 2 2\n F3 3 1\n\"01010101\"",
		"    0 1\n>F0 2 1\n  1 3 1\n F2 2 2\n  3 3 1\n\"01010101\"",
		"    0 1\n> 0 2 1\n  1 3 1\n  2 2 2\n F3 3 1\n\"01010100\"",
		"    0 1\n> 0 2 1\n F1 3 1\n F2 2 2\n  3 3 1\n\"01010100\"",
		"    0 1\n> 0 2 1\n  1 3 1\n  2 2 2\n F3 3 1\n\"11010100\"",
		"    0 1\n> 0 2 1\n F1 3 1\n  2 2 2\n  3 3 1\n\"11010100\"",
		"    0 1\n> 0 2 1\n  1 3 1\n  2 2 2\n F3 3 1\n\"11010101\"",
		"    0 1\n>F0 2 1\n F1 3 1\n F2 2 2\n  3 3 1\n\"11010101\"",
		"    0 1\n F0 1 1\n>F1 3 0\n  3 3 3\n\"\"",
		"    0 1\n F0 1 1\n>F1 3 0\n  3 3 3\n\"10101\"",
		"    0 1\n F0 1 1\n>F1 3 0\n  3 3 3\n\"01010\"",
		"    0 1\n F0 1 1\n>F1 3 0\n  3 3 3\n\"101000\"",
	}

	results := []string{
		"1 Accept",
		"2 Accept",
		"0 Reject",
		"0 Reject",
		"0 Reject",
		"0 Accept",
		"0 Accept",
		"2 Accept",
		"2 Accept",
		"3 Accept",
		"1 Reject",
		"3 Accept",
		"0 Accept",
		"1 Reject",
		"1 Accept",
		"0 Reject",
		"4 Accept",
		"0 Reject",
		"3 Reject",
		"3 Reject",
		"4 Accept",
		"2 Reject",
		"2 Accept",
		"3 Accept",
		"3 Reject",
		"2 Reject",
		"2 Accept",
		"2 Reject",
		"2 Accept",
		"3 Accept",
		"3 Reject",
		"1 Reject",
		"1 Accept",
		"1 Accept",
		"0 Accept",
		"3 Reject",
		"3 Reject",
	}

	for i := 0; i < 100; i++ {
		generatedTest := generate()
		args = append(args, generatedTest)
		results = append(results, solve(generatedTest))
	}

	rand.Shuffle(len(args), func(i, j int) {
		args[i], args[j] = args[j], args[i]
		results[i], results[j] = results[j], results[i]
	})

	return []Run{{Args: args, Answer: strings.Join(results, "\n")}}
}
