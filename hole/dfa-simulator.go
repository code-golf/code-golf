package hole

import (
	"math/rand"
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

	alphabetLength := rand.Intn(len(alphabet)-11) + 1
	alphabet = alphabet[:alphabetLength]

	var inputDFA strings.Builder

	inputDFA.WriteString("    ")
	inputDFA.WriteString(strings.Join(alphabet, " "))
	inputDFA.WriteByte('\n')

	stateLength := rand.Intn(8) + 1
	startState := rand.Intn(stateLength)
	for i := 0; i < stateLength; i++ {
		if i == startState {
			inputDFA.WriteByte('>')
		} else {
			inputDFA.WriteByte(' ')
		}

		if rand.Intn(2) == 0 {
			inputDFA.WriteByte('F')
		} else {
			inputDFA.WriteByte(' ')
		}

		inputDFA.WriteString(states[i])

		for j := 0; j < alphabetLength; j++ {
			inputDFA.WriteByte(' ')
			inputDFA.WriteString(states[rand.Intn(stateLength)])
		}
		inputDFA.WriteByte('\n')

	}
	inputDFA.WriteByte('"')
	for i := 0; i < rand.Intn(2*alphabetLength); i++ {
		inputDFA.WriteString(alphabet[rand.Intn(alphabetLength)])
	}
	inputDFA.WriteByte('"')

	return inputDFA.String()
}

func dfaSimulator() []Run {
	tests := fixedTests("dfa-simulator")

	for i := 0; i < 79; i++ {
		dfa := generateDFA()
		tests = append(tests, test{dfa, solveDFA(dfa)})
	}

	return outputTests(shuffle(tests))
}
