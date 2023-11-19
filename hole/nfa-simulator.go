package hole

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
)

type NFAState struct {
	name     string
	accepted bool
}

type (
	NFAStateTransition struct{ input, stateName string }
	NFAResult          struct{ inputSegment, stateName string }
)

func solveNFA(test string) string {
	// parse
	var startState NFAState

	stateMap := map[NFAStateTransition][]NFAState{}
	nameToState := map[string]NFAState{}

	lines := strings.Split(test, "\n")
	inputNames := strings.Split(lines[0], "|")[1:]

	inputs := make([]string, len(inputNames))
	for i, name := range inputNames {
		inputs[i] = strings.TrimSpace(name)
	}

	firstTestRow := 0
	for i, row := range lines[1:] {
		// stop when the first test case is detected
		if row == "ε" || len(row) < 4 || []rune(row)[3] != ' ' {
			firstTestRow = i + 1
			break
		}
	}

	for _, stateRow := range lines[1:firstTestRow] {
		stateName := string([]rune(stateRow)[2])
		newState := NFAState{stateName, []rune(stateRow)[1] == 'F'}
		nameToState[stateName] = newState

		if []rune(stateRow)[0] == '→' {
			startState = newState
		}
	}

	for _, stateRow := range lines[1:firstTestRow] {
		stateName := string([]rune(stateRow)[2])

		rowData := strings.Split(stateRow, "|")
		for i, potentialStates := range rowData[1 : len(rowData)-1] {
			input := inputs[i]
			newStateTransition := NFAStateTransition{input, stateName}
			stateMap[newStateTransition] = make([]NFAState, 0)

			if potentialStates == " ∅ " {
				continue
			}

			stateTransitionData := strings.Split(potentialStates[1:len(potentialStates)-1], ",")
			for _, targetState := range stateTransitionData {
				stateMap[newStateTransition] = append(stateMap[newStateTransition], nameToState[targetState])
			}
		}
	}

	var output strings.Builder

	for _, line := range lines[firstTestRow:] {
		nfaResult := []NFAState{startState}
		if line != "ε" {
			inputString := line
			nfaResult = recurseNFA(startState, inputString, stateMap, map[NFAResult][]NFAState{})
		}

		// format output
		accept := false
		resultStateNamesMap := map[string]bool{}
		for _, state := range nfaResult {
			accept = accept || state.accepted
			resultStateNamesMap[state.name] = true
		}

		if len(nfaResult) == 0 {
			output.WriteString("∅")
		} else {
			resultStateNames := make([]string, 0)
			for stateName := range resultStateNamesMap {
				resultStateNames = append(resultStateNames, stateName)
			}

			sort.Strings(resultStateNames)
			output.WriteString(fmt.Sprintf("{%s}", strings.Join(resultStateNames, ",")))
		}

		if accept {
			output.WriteString(" Accept\n")
		} else {
			output.WriteString(" Reject\n")
		}
	}

	resultString := output.String()
	return resultString[:len(resultString)-1]
}

// solve recursively
func recurseNFA(currentState NFAState, input string, stateMap map[NFAStateTransition][]NFAState, resultMap map[NFAResult][]NFAState) []NFAState {
	if len(input) == 1 {
		return stateMap[NFAStateTransition{input, currentState.name}]
	}

	// hit cache to avoid redundant recursion
	nfaResultObject := NFAResult{input, currentState.name}
	lookupNFAResult, ok := resultMap[nfaResultObject]
	if ok {
		return lookupNFAResult
	}

	results := make([]NFAState, 0)
	for _, nextState := range stateMap[NFAStateTransition{string(input[0]), currentState.name}] {
		result := recurseNFA(nextState, input[1:], stateMap, resultMap)

		// cache result
		resultMap[nfaResultObject] = result
		for _, nfaState := range result {
			results = append(results, nfaState)
		}
	}

	return results
}

func generateNFA() string {
	alphabet := shuffle([]string{
		"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m",
		"n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
		"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
	})
	states := shuffle([]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"})

	alphabetLength := rand.Intn(len(alphabet)-18) + 1
	alphabet = alphabet[:alphabetLength]

	var inputNFA strings.Builder

	inputNFA.WriteString("    | ")
	inputNFA.WriteString(strings.Join(alphabet, " | "))
	inputNFA.WriteString(" |\n")

	stateLength := rand.Intn(6) + 1
	startState := rand.Intn(stateLength)

	shuffleStates := make([]string, stateLength)
	for j, state := range states[:stateLength] {
		shuffleStates[j] = state
	}

	for i := 0; i < stateLength; i++ {
		if i == startState {
			inputNFA.WriteString("→")
		} else {
			inputNFA.WriteByte(' ')
		}

		if rand.Intn(2) == 0 {
			inputNFA.WriteByte('F')
		} else {
			inputNFA.WriteByte(' ')
		}

		inputNFA.WriteString(states[i])
		inputNFA.WriteString(" |")
		for j := 0; j < alphabetLength; j++ {

			if rand.Intn(4) == 0 {
				inputNFA.WriteString(" ∅ ")
			} else {
				shuffle(shuffleStates)
				possibleStatesCount := rand.Intn(len(shuffleStates)) + 1

				chosenStates := make([]string, possibleStatesCount)
				for j, state := range shuffleStates[:possibleStatesCount] {
					chosenStates[j] = state
				}

				sort.Strings(chosenStates)

				inputNFA.WriteString("{" + strings.Join(chosenStates, ",") + "}")

			}
			inputNFA.WriteString("|")

		}
		inputNFA.WriteByte('\n')

	}

	inputStringCount := rand.Intn(4) + 1

	for m := 0; m < inputStringCount; m++ {
		inputLength := rand.Intn(2 * alphabetLength)
		if inputLength == 0 {
			inputNFA.WriteString("ε")
		} else {
			for i := 0; i < inputLength; i++ {
				inputNFA.WriteString(alphabet[rand.Intn(alphabetLength)])
			}
		}
		if m+1 < inputStringCount {
			inputNFA.WriteByte('\n')
		}
	}

	return inputNFA.String()
}

func nfaSimulator() []Run {
	args := []string{
		"    | a | b | c |\n→ 0 |{0}|{0}|{0,1}| \n  1 |{2}| ∅ |  ∅ |\n  2 | ∅ |{3}|  ∅ |\n F3 | ∅ | ∅ | ∅ |\nacbcab\nε\nacbca",
		"    | a | b | c |\n→ 0 |{0}|{0}|{0,1}| \n  1 |{2}| ∅ |  ∅ |\n  2 | ∅ |{3}|  ∅ |\n F3 |{3}|{3}|{3}|\nacbcababc",
		"    | a | b | c |\n→ 0 | ∅ | ∅ |{1}| \n  1 |{2}| ∅ |  ∅ |\n  2 | ∅ |{3}|  ∅ |\n F3 |{3}|{3}|{3}|\ncab\nacbcab",
		"    | a | b | c |\n→ 0 |{0}|{0}|{0,1}| \n  1 |{2}|{2}|{2}|\n  2 |{3}|{3}|{3}|\n F3 | ∅ | ∅ | ∅ |\nacbcabcba\ncabca",
		"    | w | e | b | a | y | x |\n→ 1 |{1,2}|{1,5}|{1}|{1}|{1}|{1}|\n  2 | ∅ |{3}| ∅ | ∅ | ∅ | ∅ |\n  3 | ∅ | ∅ |{4}| ∅ | ∅ | ∅ |\n F4 | ∅ | ∅ | ∅ | ∅ | ∅ | ∅ |\n  5 | ∅ | ∅ |{6}| ∅ | ∅ | ∅ |\n  6 | ∅ | ∅ | ∅ |{7}| ∅ | ∅ |\n  7 | ∅ | ∅ | ∅ | ∅ |{8}| ∅ |\n F8 | ∅ | ∅ | ∅ | ∅ | ∅ | ∅ |\nebay\nwwweb\nxwe",
	}
	results := []string{
		"{0,3} Accept\n{0} Reject\n{0,2} Reject",
		"{0,1,3} Accept",
		"{3} Accept\n∅ Reject",
		"{0,3} Accept\n{0,2} Reject",
		"{1,8} Accept\n{1,4,6} Accept\n{1,3,5} Reject",
	}

	for i := 0; i < 12; i++ {
		dfa := generateNFA()
		args = append(args, dfa)
		results = append(results, solveNFA(dfa))
	}

	rand.Shuffle(len(args), func(i, j int) {
		args[i], args[j] = args[j], args[i]
		results[i], results[j] = results[j], results[i]
	})

	return []Run{{Args: args, Answer: strings.Join(results, "\n\n")}}
}
