package hole

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"slices"
	"strings"
	"sync"
	"syscall"
	"time"
	"unicode"

	"github.com/agnivade/levenshtein"
	"github.com/buildkite/terminal-to-html/v3"
	"github.com/code-golf/code-golf/config"
	hungarianAlgorithm "github.com/oddg/hungarian-algorithm"
)

var timeout = 5 * time.Second

// Increase the timeout under e2e as the hardware is less powerful than live.
func init() {
	if _, e2e := os.LookupEnv("E2E"); e2e {
		timeout = 10 * time.Second
	}
}

//go:embed answers
var answers embed.FS

// All ASCII whitespace except newline, up to a newline or the end.
var stdoutTrimmer = regexp.MustCompile(`[\t\x0B\f\r ]+(?:\n|$)`)

// Run holds the results of running a given solution once.
type Run struct {
	Answer            string        `json:"answer"`
	ItemDelimiter     string        `json:"item_delimiter"`
	MultisetDelimiter string        `json:"multiset_delimiter"`
	Args              []string      `json:"args"`
	ExitCode          int           `json:"exit_code"`
	Pass              bool          `json:"pass"`
	Stderr            string        `json:"stderr"`
	Stdout            string        `json:"stdout"`
	Time              time.Duration `json:"time_ns"`
	Timeout           bool          `json:"timeout"`

	// This is a bit hacky, the only way to discover how long an assembly
	// solution is is to compile it so we store it here but don't JSON it.
	ASMBytes int `json:"-"`
}

func getClosestAnswer(anyAnswer, stdout, itemDelimiter, multisetDelimiter string) string {
	answerMultisets := []string{anyAnswer}
	stdoutMultisets := []string{stdout}
	if multisetDelimiter != "" {
		answerMultisets = strings.Split(anyAnswer, multisetDelimiter)
		stdoutMultisets = strings.Split(stdout, multisetDelimiter)
	}
	closestMultisets := make([]string, len(answerMultisets))

	for i, answerMultiset := range answerMultisets {
		stdoutMultiset := ""
		if i < len(stdoutMultisets) {
			stdoutMultiset = stdoutMultisets[i]
		}
		closestMultisets[i] = getClosestMultiset(answerMultiset, stdoutMultiset, itemDelimiter)
	}
	return strings.Join(closestMultisets, multisetDelimiter)
}

func getClosestMultiset(anyAnswer, stdout, itemDelimiter string) string {
	expectedItems := strings.Split(anyAnswer, itemDelimiter)
	expectedItemsReordered := make([]string, len(expectedItems))
	userItems := strings.Split(stdout, itemDelimiter)

	expectedItemsMap := make(map[string]int)
	for _, expected := range expectedItems {
		expectedItemsMap[expected]++
	}

	// Match items that are correct
	matches := 0
	for i, user := range userItems {
		if i < len(expectedItems) && expectedItemsMap[user] > 0 {
			expectedItemsReordered[i] = user
			expectedItemsMap[user]--
			userItems[i] = ""
			matches++
		}
	}

	// Process mismatched items
	if matches < len(expectedItems) {

		// Calculate indices of expected & user items that couldn't be matched be equality
		unmatchedExpectedIndices := []int{}
		unmatchedUserIndices := []int{}

		for i, expected := range expectedItems {
			if expectedItemsMap[expected] > 0 {
				unmatchedExpectedIndices = append(unmatchedExpectedIndices, i)
				expectedItemsMap[expected]--
			}
		}

		for i, user := range userItems {
			if user != "" {
				unmatchedUserIndices = append(unmatchedUserIndices, i)
			}
		}

		n := max(len(unmatchedExpectedIndices), len(unmatchedUserIndices))

		permutation := make([]int, n)
		for i := range permutation {
			permutation[i] = i
		}

		// If there are not many wrong items, try to match them
		// otherwise, use the above identity permutation
		if n <= 32 {
			dist := make([][]int, n)
			for i := range dist {
				dist[i] = make([]int, n)
				for j := range dist {
					if j >= len(unmatchedExpectedIndices) {
						dist[i][j] = len(userItems[unmatchedUserIndices[i]])
					} else if i >= len(unmatchedUserIndices) {
						dist[i][j] = len(expectedItems[unmatchedExpectedIndices[j]])
					} else {
						dist[i][j] = levenshtein.ComputeDistance(expectedItems[unmatchedExpectedIndices[j]], userItems[unmatchedUserIndices[i]])
					}
				}
			}

			permutation, _ = hungarianAlgorithm.Solve(dist)
		}

		k := 0
		for _, i := range permutation {
			if k >= len(expectedItemsReordered) {
				break
			}
			if i < len(unmatchedExpectedIndices) {
				for expectedItemsReordered[k] != "" {
					k++
				}
				expectedItemsReordered[k] = expectedItems[unmatchedExpectedIndices[i]]
			}
		}
	}

	return strings.Join(expectedItemsReordered, itemDelimiter)
}

// Play a given hole, in a given lang, with given code and return the runs.
func Play(
	ctx context.Context, hole *config.Hole, lang *config.Lang, code string,
) (runs []Run) {
	switch hole.ID {
	case "arabic-to-roman", "roman-to-arabic":
		runs = arabicToRoman(hole.ID == "roman-to-arabic")
	case "arrows":
		runs = arrows()
	case "billiards":
		runs = billiards()
	case "brainfuck":
		runs = brainfuck()
	case "card-number-validation":
		runs = cardNumberValidation()
	case "day-of-week":
		runs = dayOfWeek()
	case "dfa-simulator":
		runs = dfaSimulator()
	case "ellipse-perimeters":
		runs = ellipsePerimeters()
	case "forsyth-edwards-notation":
		runs = forsythEdwardsNotation()
	case "fractions":
		runs = fractions()
	case "game-of-life":
		runs = gameOfLife()
	case "gray-code-decoder", "gray-code-encoder":
		runs = grayCode(hole.ID == "gray-code-decoder")
	case "css-grid":
		runs = cssGrid()
	case "isbn":
		runs = isbn()
	case "intersection":
		runs = intersection()
	case "jacobi-symbol":
		runs = jacobiSymbol()
	case "levenshtein-distance":
		runs = levenshteinDistance()
	case "lucky-tickets":
		runs = luckyTickets()
	case "mahjong":
		runs = mahjong()
	case "maze":
		runs = maze()
	case "medal-tally":
		runs = medalTally()
	case "morse-decoder", "morse-encoder":
		runs = morse(hole.ID == "morse-decoder")
	case "musical-chords":
		runs = musicalChords()
	case "nfa-simulator":
		runs = nfaSimulator()
	case "ordinal-numbers":
		runs = ordinalNumbers()
	case "p-adic-expansion":
		runs = pAdicExpansion()
	case "pangram-grep":
		runs = pangramGrep()
	case "poker":
		runs = poker()
	case "qr-decoder", "qr-encoder":
		runs = qr(hole.ID == "qr-decoder")
	case "quadratic-formula":
		runs = quadraticFormula()
	case "quine":
		runs = []Run{{Args: []string{}, Answer: code}}
	case "repeating-decimals":
		runs = repeatingDecimals()
	case "reverse-polish-notation":
		runs = reversePolishNotation()
	case "reversi":
		runs = reversi()
	case "seven-segment":
		runs = sevenSegment()
	case "si-units":
		runs = siUnits()
	case "spelling-numbers":
		runs = spellingNumbers()
	case "sudoku", "sudoku-v2":
		runs = sudoku(hole.ID == "sudoku-v2")
	case "ten-pin-bowling":
		runs = tenPinBowling()
	case "time-distance":
		runs = timeDistance()
	case "transpose-sentence":
		runs = transposeSentence()
	case "turtle":
		runs = turtle()
	case "zeckendorf-representation":
		runs = zeckendorfRepresentation()
	case "zodiac-signs":
		runs = zodiacSigns()

	// Holes with fixed test cases.
	case "css-colors":
		runs = outputTests(shuffle(fixedTests(hole.ID)))
	case "emojify", "mnist", "rock-paper-scissors-spock-lizard",
		"united-states":
		runs = outputMultirunTests(fixedTests(hole.ID))
	case "floyd-steinberg-dithering", "hexdump", "proximity-grid", "star-wars-opening-crawl":
		runs = outputTestsWithSep("\n\n", shuffle(fixedTests(hole.ID)))

	// Holes with no arguments and a static answer.
	default:
		// ¯\_(ツ)_/¯ cannot embed file answers/√2.txt: invalid name √2.txt
		id := hole.ID
		if id == "√2" {
			id = "root-2"
		}

		if b, err := answers.ReadFile("answers/" + id + ".txt"); err != nil {
			panic(err)
		} else {
			answer := string(bytes.TrimSuffix(b, []byte{'\n'}))
			runs = []Run{{Args: []string{}, Answer: answer}}
		}
	}

	// Run all the runs in parallel to reduce the wall clock time.
	var wg sync.WaitGroup
	wg.Add(len(runs))

	for i := range runs {
		go func(run *Run) {
			if err := play(ctx, hole, lang, code, run); err != nil {
				log.Println(err)
			}

			wg.Done()
		}(&runs[i])
	}

	wg.Wait()

	return
}

func play(
	ctx context.Context, hole *config.Hole, lang *config.Lang, code string, run *Run,
) error {
	// Preprocess code.
	switch lang.ID {
	case "clojure":
		// Appending (print) prevents implicit output of the last form, if it
		// is not nil. This seems to be a quirk of the Babashka interpreter
		// that only occurs when providing code via a command line argument.
		code += "(print)"
	case "jq":
		// Prevent trivial quines. Error out and return early.
		if hole.ID == "quine" && json.Valid([]byte(code)) {
			run.Stderr = "Quine in jq must not be valid JSON."
			return nil
		}
	case "k":
		if hole.ID == "quine" {
			length := len(code)
			var newCode []byte

			// Disable implicit output by inserting a ';' before all newlines,
			// except when the next line begins with a space (for a continuation).
			for i := range length {
				x := code[i]
				if x != '\n' || i+1 < length && code[i+1] == ' ' {
					newCode = append(newCode, x)
				} else {
					newCode = append(newCode, ';', '\n')
				}
			}

			code = string(newCode)
		} else {
			code += "\n"
		}
	case "php":
		code = "<?php " + code + " ;"
	case "tex":
		// Prevent trivial quines. Error out and return early.
		if hole.ID == "quine" && !strings.Contains(code, `\`) {
			run.Stderr = `Quine in TeX must have at least one '\' character.`
			return nil
		}
	}

	var stderr, stdout bytes.Buffer

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "/usr/bin/run-lang")
	cmd.Dir = "/langs/" + lang.ID
	cmd.Env = append([]string{"HOME=/tmp"}, lang.Env...)
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	cmd.WaitDelay = time.Second
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWIPC | syscall.CLONE_NEWNET |
			syscall.CLONE_NEWNS | syscall.CLONE_NEWPID | syscall.CLONE_NEWUTS,
	}

	// Assembly bytes pipe.
	var asmBytesRead, asmBytesWrite *os.File
	if lang.ID == "assembly" {
		var err error
		if asmBytesRead, asmBytesWrite, err = os.Pipe(); err != nil {
			return err
		}

		cmd.ExtraFiles = []*os.File{asmBytesWrite}
	}

	// Language arguments. Clone because we intend to mutate.
	cmd.Args = slices.Clone(lang.Args)
	if hole.ID == "quine" && lang.ArgsQuine != nil {
		cmd.Args = slices.Clone(lang.ArgsQuine)
	}

	// Pass code via args or stdin.
	if i := slices.Index(cmd.Args, "$code"); i != -1 {
		cmd.Args[i] = code
	} else {
		cmd.Stdin = strings.NewReader(code)
	}

	// Run arguments.
	switch lang.ID {
	case "awk", "brainfuck", "fish":
		// Hole args passed through stdin for these langs separated by a null byte
		args := ""
		for _, arg := range run.Args {
			args += arg + "\x00"
		}
		cmd.Stdin = strings.NewReader(args)
	case "sed":
		// For sed we always need to append a null byte, even if no args exist
		args := strings.Join(run.Args, "\x00") + "\x00"
		cmd.Stdin = strings.NewReader(args)
	default:
		cmd.Args = append(cmd.Args, run.Args...)
	}

	err := cmd.Run()

	deadline, _ := ctx.Deadline()
	run.Time = timeout - time.Until(deadline)

	if err != nil {
		if err, ok := err.(*exec.ExitError); ok {
			run.ExitCode = err.ExitCode()
		}

		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			run.Timeout = true
			fmt.Fprint(&stderr, "Killed for exceeding the ", timeout, " timeout.")
		} else {
			stderr.WriteString(err.Error())
		}
	}

	// Actual byte count is printed by the assembler.
	if lang.ID == "assembly" {
		if _, err := fmt.Fscanf(asmBytesRead, "%d", &run.ASMBytes); err != nil {
			return err
		}
		asmBytesRead.Close()
	}

	const maxLength = 128 * 1024 // 128 KiB

	// Trim trailing whitespace.
	stderrBytes := bytes.TrimRightFunc(stderr.Next(maxLength), unicode.IsSpace)

	// Convert ANSI escapes into HTML, for coloured error messages.
	run.Stderr = terminal.Render(stderrBytes)

	// Bodge, suppress solitary "&nbsp;" that can be emitted.
	if run.Stderr == "&nbsp;" {
		run.Stderr = ""
	}

	stdoutBytes := stdout.Next(maxLength)

	// Postprocess sed output to turn null bytes into newlines.
	if lang.ID == "sed" {
		stdoutBytes = bytes.ReplaceAll(stdoutBytes, []byte("\x00"), []byte("\n"))
	}

	// Trim trailing whitespace on each line, and then trailing empty lines.
	// Quine solutions are obviously left untouched.
	if hole.ID == "quine" {
		run.Stdout = string(stdoutBytes)
	} else {
		run.Stdout = string(bytes.TrimRight(stdoutTrimmer.ReplaceAll(
			stdoutBytes, []byte{'\n'}), "\n"))
	}

	// Timeouts and whitespace only output never pass.
	if !run.Timeout && len(strings.TrimSpace(run.Stdout)) != 0 {
		if hole.ItemDelimiter != "" {
			run.Answer = getClosestAnswer(run.Answer, run.Stdout, hole.ItemDelimiter, hole.MultisetDelimiter)
		}

		if hole.CaseFold {
			run.Pass = strings.EqualFold(run.Answer, run.Stdout)
		} else {
			run.Pass = run.Answer == run.Stdout
		}
	}

	run.MultisetDelimiter = hole.MultisetDelimiter
	run.ItemDelimiter = hole.ItemDelimiter

	return nil
}
