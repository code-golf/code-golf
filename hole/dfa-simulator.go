package hole

import (
	"math/rand/v2"
	"strings"
)

type DFAState struct {
	name     string
	accepted bool
}

type DFAStateTransition struct{ input, stateName string }

func solveDFA(g string) string {
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
			alphabetWord := alphabet[i]
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

func generateDFA() string {
	alphabet := shuffle([]string{
		"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m",
		"n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
		"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
	})
	states := shuffle([]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"})

	alphabetLength := rand.IntN(len(alphabet)-11) + 1
	alphabet = alphabet[:alphabetLength]

	var inputDFA strings.Builder

	inputDFA.WriteString("    ")
	inputDFA.WriteString(strings.Join(alphabet, " "))
	inputDFA.WriteByte('\n')

	stateLength := rand.IntN(8) + 1
	startState := rand.IntN(stateLength)
	for i := range stateLength {
		if i == startState {
			inputDFA.WriteByte('>')
		} else {
			inputDFA.WriteByte(' ')
		}

		if rand.IntN(2) == 0 {
			inputDFA.WriteByte('F')
		} else {
			inputDFA.WriteByte(' ')
		}

		inputDFA.WriteString(states[i])

		for range alphabetLength {
			inputDFA.WriteByte(' ')
			inputDFA.WriteString(states[rand.IntN(stateLength)])
		}
		inputDFA.WriteByte('\n')

	}
	inputDFA.WriteByte('"')
	for i := 0; i < rand.IntN(2*alphabetLength); i++ {
		inputDFA.WriteString(alphabet[rand.IntN(alphabetLength)])
	}
	inputDFA.WriteByte('"')

	return inputDFA.String()
}

func dfaSimulator() []Run {
	tests := fixedTests("dfa-simulator")

	for range 79 {
		dfa := generateDFA()
		tests = append(tests, test{dfa, solveDFA(dfa)})
	}

	return outputTests(shuffle(tests))
}
