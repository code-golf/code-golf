package hole

import (
	"bytes"
	"context"
	"embed"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"syscall"
	"time"
	"unicode"

	"github.com/code-golf/code-golf/config"
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
	Answer   string        `json:"answer"`
	Args     []string      `json:"args"`
	ExitCode int           `json:"exit_code"`
	Pass     bool          `json:"pass"`
	Stderr   string        `json:"stderr"`
	Stdout   string        `json:"stdout"`
	Time     time.Duration `json:"time_ns"`
	Timeout  bool          `json:"timeout"`

	// This is a bit hacky, the only way to discover how long an assembly
	// solution is is to compile it so we store it here but don't JSON it.
	ASMBytes int `json:"-"`
}

func preprocessKCode(holeID, code string) string {
	if holeID == "quine" {
		length := len(code)
		var newCode []byte

		// Disable implicit output by inserting a ';' before all newlines,
		// except when the next line begins with a space (for a continuation).
		for i := 0; i < length; i++ {
			x := code[i]
			if x != '\n' || i+1 < length && code[i+1] == ' ' {
				newCode = append(newCode, x)
			} else {
				newCode = append(newCode, ';', '\n')
			}
		}

		return string(newCode)
	} else {
		return code + "\n"
	}
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
	case "brainfuck":
		runs = brainfuck()
	case "css-colors":
		runs = cssColors()
	case "day-of-week":
		runs = dayOfWeek()
	case "dfa-simulator":
		runs = dfaSimulator()
	case "ellipse-perimeters":
		runs = ellipsePerimeters()
	case "emojify":
		runs = emojify()
	case "forsyth-edwards-notation":
		runs = forsythEdwardsNotation()
	case "fractions":
		runs = fractions()
	case "game-of-life":
		runs = gameOfLife()
	case "gray-code-encoder", "gray-code-decoder":
		runs = grayCode(hole.ID == "gray-code-decoder")
	case "hexdump":
		runs = hexdump()
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
	case "proximity-grid":
		runs = proximityGrid()
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
	case "rock-paper-scissors-spock-lizard":
		runs = rockPaperScissorsSpockLizard()
	case "seven-segment":
		runs = sevenSegment()
	case "si-units":
		runs = siUnits()
	case "spelling-numbers":
		runs = spellingNumbers()
	case "star-wars-opening-crawl":
		runs = starWarsOpeningCrawl()
	case "sudoku", "sudoku-v2":
		runs = sudoku(hole.ID == "sudoku-v2")
	case "ten-pin-bowling":
		runs = tenPinBowling()
	case "time-distance":
		runs = timeDistance()
	case "united-states":
		runs = unitedStates()
	case "turtle":
		runs = turtle()
	case "zodiac-signs":
		runs = zodiacSigns()
	case "zeckendorf-representation":
		runs = zeckendorfRepresentation()
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
	var stderr, stdout bytes.Buffer
	var asmBytesRead, asmBytesWrite *os.File

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "/usr/bin/run-lang")
	cmd.Dir = "/langs/" + lang.ID
	cmd.Env = []string{}
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	cmd.WaitDelay = time.Second
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWIPC | syscall.CLONE_NEWNET |
			syscall.CLONE_NEWNS | syscall.CLONE_NEWPID | syscall.CLONE_NEWUTS,
	}

	// Interpreter
	switch lang.ID {
	case "assembly":
		var err error
		if asmBytesRead, asmBytesWrite, err = os.Pipe(); err != nil {
			return err
		}

		cmd.Args = []string{"/usr/bin/defasm", "--size-out=3", "-w", "-r"}
		cmd.ExtraFiles = []*os.File{asmBytesWrite}
	case "awk":
		cmd.Args = []string{"/usr/bin/gawk", "-v", "RS=\\0", code}
	case "bash":
		cmd.Args = []string{"/usr/bin/bash", "-s", "-"}
	case "brainfuck":
		cmd.Args = []string{"/usr/bin/brainfuck", "-xc", code}
	case "c":
		cmd.Args = []string{"/usr/bin/tcc", "-run", "-"}
	case "clojure":
		// Appending (print) prevents implicit output of the last form, if it is not nil.
		// This seems to be a quirk of the Babashka interpreter that only occurs when
		// providing code via a command line argument.
		cmd.Args = []string{"/usr/bin/clojure", "-e", code + "(print)"}
	case "coconut":
		cmd.Args = []string{"/usr/bin/coconut", "--quiet", "--target", "sys", "--keep-lines", "--argv"}
	case "crystal":
		cmd.Args = []string{"/usr/bin/crystal", "run", "--stdin-filename", "code.cr", "--"}
		cmd.Env = []string{"CRYSTAL_CACHE_DIR=/tmp", "PATH=/usr/bin:/bin"}
	case "d":
		cmd.Args = []string{"/usr/bin/ldc2", "--enable-color=true", "--run", "-"}
		cmd.Env = []string{"PATH=/usr/bin"}
	case "elixir":
		cmd.Args = []string{"/usr/local/bin/elixir", "-e", code, "--"}
		cmd.Env = []string{"LANG=C.UTF-8", "PATH=/usr/local/bin:/usr/bin:/bin"}
	case "factor":
		cmd.Args = []string{"/factor/factor", "/proc/self/fd/0"}
		cmd.Env = []string{"XDG_CACHE_HOME=/tmp"}
	case "fish":
		cmd.Args = []string{"/usr/bin/fish", "--no-prng", "-c", code, "-u"}
	case "forth":
		cmd.Args = []string{"/usr/bin/forth", "/proc/self/fd/0"}
	case "golfscript":
		cmd.Args = []string{"/usr/bin/golfscript", "-n", "-e", code}
		if hole.ID == "quine" {
			cmd.Args = append(cmd.Args, "-q")
		}
		cmd.Args = append(cmd.Args, "--")
	case "hexagony":
		cmd.Args = []string{"/usr/bin/hexagony", "-d", "-"}
	case "j":
		cmd.Args = []string{"/usr/bin/j", "/tmp/code.ijs"}
	case "janet":
		cmd.Args = []string{"/usr/bin/janet", "/proc/self/fd/0"}
	case "k":
		cmd.Args = []string{"/usr/bin/kwrapper", "/tmp/code.k"}
	case "javascript":
		cmd.Args = []string{"/usr/bin/d8", "-e", code, "--"}
	case "julia":
		cmd.Args = []string{"/usr/bin/julia", "--color=yes", "/proc/self/fd/0"}
		cmd.Env = []string{"HOME=/"}
	case "nim":
		cmd.Args = []string{"/usr/bin/nim", "--colors:on", "-o:/tmp/code", "-r", "c", "-"}
	case "ocaml":
		cmd.Args = []string{"/usr/bin/ocaml", "/proc/self/fd/0"}
	case "perl":
		cmd.Args = []string{"/usr/bin/perl", "-E", code, "--"}
	case "php":
		cmd.Args = []string{"/usr/bin/php", "--"}
	case "powershell":
		cmd.Args = []string{"/usr/bin/powershell"}

		// Require explicit output for Quine to prevent trivial solutions.
		if hole.ID == "quine" {
			cmd.Args = append(cmd.Args, "--explicit")
		}
	case "prolog":
		cmd.Args = []string{"/usr/bin/prolog", "-g", "halt", "/tmp/code.pl"}
	case "python":
		// Force the stdout and stderr streams to be unbuffered.
		cmd.Args = []string{"/usr/bin/python", "-u", "-"}
	case "r":
		cmd.Args = []string{"/usr/bin/Rscript", "-"}

		// Disable implicit output for Quine to prevent trivial solutions.
		if hole.ID == "quine" {
			cmd.Args = []string{"/usr/bin/Rscript", "-e", "source('stdin')"}
		}
	case "sed":
		cmd.Args = []string{"/usr/bin/sed", "-E", "-z", "--sandbox", "-u", "--", code}
	case "swift":
		cmd.Args = []string{"/usr/bin/swift", "-module-cache-path", "/tmp", "-"}
	case "tcl":
		cmd.Args = []string{"/usr/bin/tcl", "/proc/self/fd/0"}
	case "tex":
		cmd.Args = []string{"/usr/bin/tex", code}

		// Require a backslash for Quine to prevent trivial solutions.
		// Don't even run the code; just mark error and return.
		if hole.ID == "quine" && !strings.Contains(code, `\`) {
			run.Stderr = `Quine in TeX must have at least one '\' character.`
			return nil
		}
	default:
		cmd.Args = []string{"/usr/bin/" + lang.ID, "-"}
	}

	// Args
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

	// Code
	switch lang.ID {
	case "awk", "brainfuck", "elixir", "fish", "golfscript", "javascript",
		"perl", "sed", "tex":
		// For these langs, code is passed as an argument above.
	case "k":
		cmd.Stdin = strings.NewReader(preprocessKCode(hole.ID, code))
	case "php":
		cmd.Stdin = strings.NewReader("<?php " + code + " ;")
	default:
		cmd.Stdin = strings.NewReader(code)
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
	run.Stderr = string(bytes.TrimRightFunc(stderr.Next(maxLength), unicode.IsSpace))

	stdoutContents := stdout.Next(maxLength)

	// Postprocess sed output to turn null bytes into newlines.
	if lang.ID == "sed" {
		stdoutContents = bytes.ReplaceAll(stdoutContents, []byte("\x00"), []byte("\n"))
	}

	// Trim trailing whitespace on each line, and then trailing empty lines.
	// Quine solutions are obviously left untouched.
	if hole.ID == "quine" {
		run.Stdout = string(stdoutContents)
	} else {
		run.Stdout = string(bytes.TrimRight(stdoutTrimmer.ReplaceAll(
			stdoutContents, []byte{'\n'}), "\n"))
	}

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
