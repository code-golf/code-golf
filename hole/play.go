package hole

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"slices"
	"strings"
	"syscall"
	"time"
	"unicode"

	"github.com/buildkite/terminal-to-html/v3"
	"github.com/code-golf/code-golf/config"
	"golang.org/x/sync/errgroup"
)

var timeout = 5 * time.Second

// Increase the timeout under e2e as the hardware is less powerful than live.
func init() {
	if _, e2e := os.LookupEnv("E2E"); e2e {
		timeout = 10 * time.Second
	}
}

// All ASCII whitespace except newline, up to a newline or the end.
var perLineTrimmer = regexp.MustCompile(`[\t\x0B\f\r ]+(?:\n|$)`)

func trimPerLine(bytesSlice []byte) string {
	return string(bytes.TrimRight(perLineTrimmer.ReplaceAll(
		bytesSlice, []byte{'\n'}), "\n"))
}

// Run holds the results of running a given solution once.
type Run struct {
	Answer                string        `json:"answer"`
	Code                  string        `json:"-"`
	MultisetItemDelimiter string        `json:"multiset_item_delimiter"`
	OutputDelimiter       string        `json:"output_delimiter"`
	Args                  []string      `json:"args"`
	ExitCode              int           `json:"exit_code"`
	Pass                  bool          `json:"pass"`
	Stderr                string        `json:"stderr"`
	Stdout                string        `json:"stdout"`
	Time                  time.Duration `json:"time_ns"`
	Timeout               bool          `json:"timeout"`

	// This is a bit hacky, the only way to discover how long an assembly
	// solution is is to compile it so we store it here but don't JSON it.
	ASMBytes int `json:"-"`
}

// Play a given hole, in a given lang, with given code and return the runs.
func Play(
	ctx context.Context, hole *config.Hole, lang *config.Lang, code string,
) ([]Run, error) {
	var answers []Answer

	switch hole.ID {

	// Holes with fixed test cases.
	case "css-colors":
		answers = outputTests(shuffle(fixedTests(hole.ID)))
	case "emojify", "flags", "rock-paper-scissors-spock-lizard", "tic-tac-toe", "united-states":
		answers = outputMultirunTests(fixedTests(hole.ID))
	case "floyd-steinberg-dithering", "hexdump", "minesweeper", "proximity-grid", "star-wars-opening-crawl":
		answers = outputTestsWithSep("\n\n", shuffle(fixedTests(hole.ID)))

	// Holes with a static answer or answer func.
	default:
		if hole.AnswerFunc != nil {
			answers = hole.AnswerFunc()
		} else {
			answers = []Answer{{Args: []string{}, Answer: hole.Answer}}
		}
	}

	runs := make([]Run, len(answers))

	// Run all the runs in parallel to reduce the wall clock time.
	var g errgroup.Group
	for i, answer := range answers {
		runs[i] = Run{Args: answer.Args, Answer: answer.Answer}

		g.Go(func() error { return play(ctx, hole, lang, code, &runs[i]) })
	}

	return runs, g.Wait()
}

func play(
	ctx context.Context, hole *config.Hole, lang *config.Lang, code string, run *Run,
) error {
	err := runCode(ctx, hole, lang, code, run)
	if err != nil {
		return err
	}

	run.Code = code
	run.OutputDelimiter = hole.OutputDelimiter
	run.MultisetItemDelimiter = hole.MultisetItemDelimiter
	run.Answer = holeJudges[hole.ID](*run)

	// Timeouts and whitespace only output never pass.
	if !run.Timeout && len(strings.TrimSpace(run.Stdout)) != 0 {
		if hole.CaseFold {
			run.Pass = strings.EqualFold(run.Answer, run.Stdout)
		} else {
			run.Pass = run.Answer == run.Stdout
		}
	}

	return nil
}

func runCode(
	ctx context.Context, hole *config.Hole, lang *config.Lang, code string, run *Run,
) error {
	// Preprocess code.
	if strings.Contains(code, "\x00") {
		run.Stderr = "Solutions must not contain a literal null byte."
		return nil
	}

	switch lang.ID {
	case "05ab1e":
		// Prevent trivial quines. Error out and return early.
		if hole.ID == "quine" && len(code) > 0 && !strings.Contains(code, `"`) && !strings.Contains(code, "”") {
			run.Stderr = `Quine in 05AB1E must have at least one '"' or '”' character.`
			return nil
		}
	case "cjam":
		// Prevent trivial quines. Error out and return early.
		if hole.ID == "quine" && !strings.Contains(code, "`") {
			run.Stderr = "Quine in CJam must have at least one '`' character."
			return nil
		}
	case "clojure":
		// Appending (print) prevents implicit output of the last form, if it
		// is not nil. This seems to be a quirk of the Babashka interpreter
		// that only occurs when providing code via a command line argument.
		//
		// Add a newline in to avoid commenting it out via ;, and
		// do it twice to avoid commenting it out via #_.
		code += "\n(print)(print)"
	case "go":
		// Prevent trivial quines. Error out and return early.
		if hole.ID == "quine" && strings.Contains(code, "//go:embed") {
			run.Stderr = `Quine in Go must not use "embed".`
			return nil
		}
	case "iogii", "stax":
		// Prevent trivial quines. Error out and return early.
		if hole.ID == "quine" && len(code) > 0 && regexp.MustCompile(`^\d+\n?$`).MatchString(code) {
			run.Stderr = "Quine in " + lang.Name + " must not consist solely of numeric characters."
			return nil
		}
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
	case "kotlin":
		if hole.ID == "quine" {
			// Appending `Unit` on a newline suppresses implicit output of expressions
			// in Kotlin. The '\n' guarantees we're not appending a ';' to another ';'.
			code += "\nUnit"
		}
	case "php":
		code = "<?php " + code + " ;"
	case "racket":
		if hole.ID == "quine" {
			// Inserting `(current-print (λ (x) (void)))` before the code in the editor
			// suppresses the implicit output of expressions in Racket.
			code = "(current-print (λ (x) (void)))" + code
		}
	case "tex":
		// Prevent trivial quines. Error out and return early.
		if hole.ID == "quine" && !strings.Contains(code, `\`) {
			run.Stderr = `Quine in TeX must have at least one '\' character.`
			return nil
		}
	case "uiua":
		// Prevent trivial quines. Error out and return early.
		if hole.ID == "quine" && len(code) > 0 && (!strings.Contains(code, "&p") || strings.Contains(code, `"`)) {
			run.Stderr = "Quine in Uiua must use `&p` (without backticks) and cannot contain double quotes."
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
		defer asmBytesRead.Close()
		defer asmBytesWrite.Close()

		cmd.ExtraFiles = []*os.File{asmBytesWrite}
	}

	// Language arguments. Clone because we intend to mutate.
	if hole.ID == "quine" && lang.ArgsQuine != nil {
		cmd.Args = slices.Clone(lang.ArgsQuine)
	} else {
		cmd.Args = slices.Clone(lang.Args)
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
		// Explicitly close the writer in case defasm died before it could.
		// This prevents a very long wait in the upcoming fscanf.
		asmBytesWrite.Close()

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

	// Postprocess output in apl or sed.
	// Convert apl's carriage returns or sed's null bytes to newlines.
	if lang.ID == "apl" || lang.ID == "sed" {
		stdoutByte := "\r"

		if lang.ID == "sed" {
			stdoutByte = "\x00"
		}

		stdoutBytes = bytes.ReplaceAll(stdoutBytes, []byte(stdoutByte), []byte("\n"))
	}

	// Trim trailing whitespace on each line, and then trailing empty lines.
	// Quine solutions are obviously left untouched.
	if hole.ID == "quine" {
		run.Stdout = string(stdoutBytes)
	} else {
		run.Stdout = trimPerLine(stdoutBytes)
	}

	return nil
}
