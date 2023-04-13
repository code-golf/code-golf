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
	"syscall"
	"time"
	"unicode"
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

// All whitespace except newline, up to a newline or the end.
var stdoutTrimmer = regexp.MustCompile(`[^\S\n]+(?:\n|$)`)

var romanToASCII = strings.NewReplacer(
	"Ⅰ", "I", "Ⅱ", "II", "Ⅲ", "III", "Ⅳ", "IV", "Ⅴ", "V",
	"Ⅵ", "VI", "Ⅶ", "VII", "Ⅷ", "VIII", "Ⅸ", "IX", "Ⅹ", "X",
	"Ⅺ", "XI", "Ⅻ", "XII", "Ⅼ", "L", "Ⅽ", "C", "Ⅾ", "D", "Ⅿ", "M",
)

type Scorecard struct {
	ASMBytes, ExitCode int
	Answer             string
	Args               []string
	Pass, Timeout      bool
	Stderr, Stdout     []byte
	Took               time.Duration
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

// Play a given hole, in a given lang, with given code and return a Scorecard.
func Play(ctx context.Context, holeID, langID, code string) (score Scorecard) {
	var scores []Scorecard

	switch holeID {
	case "arabic-to-roman", "roman-to-arabic":
		scores = arabicToRoman(holeID == "roman-to-arabic")
	case "arrows":
		scores = arrows()
	case "brainfuck":
		scores = brainfuck()
	case "css-colors":
		scores = cssColors()
	case "day-of-week":
		scores = dayOfWeek()
	case "ellipse-perimeters":
		scores = ellipsePerimeters()
	case "emojify":
		scores = emojify()
	case "forsyth-edwards-notation":
		scores = forsythEdwardsNotation()
	case "fractions":
		scores = fractions()
	case "game-of-life":
		scores = gameOfLife()
	case "gray-code-encoder", "gray-code-decoder":
		scores = grayCode(holeID == "gray-code-decoder")
	case "hexdump":
		scores = hexdump()
	case "isbn":
		scores = isbn()
	case "intersection":
		scores = intersection()
	case "jacobi-symbol":
		scores = jacobiSymbol()
	case "levenshtein-distance":
		scores = levenshteinDistance()
	case "lucky-tickets":
		scores = luckyTickets()
	case "mahjong":
		scores = mahjong()
	case "maze":
		scores = maze()
	case "morse-decoder", "morse-encoder":
		scores = morse(holeID == "morse-decoder")
	case "musical-chords":
		scores = musicalChords()
	case "ordinal-numbers":
		scores = ordinalNumbers()
	case "p-adic-expansion":
		scores = pAdicExpansion()
	case "pangram-grep":
		scores = pangramGrep()
	case "poker":
		scores = poker()
	case "proximity-grid":
		scores = proximityGrid()
	case "qr-decoder", "qr-encoder":
		scores = qr(holeID == "qr-decoder")
	case "quine":
		scores = []Scorecard{{Args: []string{}, Answer: code}}
	case "repeating-decimals":
		scores = repeatingDecimals()
	case "reverse-polish-notation":
		scores = reversePolishNotation()
	case "rock-paper-scissors-spock-lizard":
		scores = rockPaperScissorsSpockLizard()
	case "seven-segment":
		scores = sevenSegment()
	case "si-units":
		scores = siUnits()
	case "spelling-numbers":
		scores = spellingNumbers()
	case "star-wars-opening-crawl":
		scores = starWarsOpeningCrawl()
	case "sudoku", "sudoku-v2":
		scores = sudoku(holeID == "sudoku-v2")
	case "ten-pin-bowling":
		scores = tenPinBowling()
	case "time-distance":
		scores = timeDistance()
	case "united-states":
		scores = unitedStates()
	case "turtle":
		scores = turtle()
	case "zodiac-signs":
		scores = zodiacSigns()
	case "zeckendorf-representation":
		scores = zeckendorfRepresentation()
	default:
		// ¯\_(ツ)_/¯ cannot embed file answers/√2.txt: invalid name √2.txt
		if holeID == "√2" {
			holeID = "root-2"
		}

		if b, err := answers.ReadFile("answers/" + holeID + ".txt"); err != nil {
			panic(err)
		} else {
			answer := string(bytes.TrimSuffix(b, []byte{'\n'}))
			scores = []Scorecard{{Args: []string{}, Answer: answer}}
		}
	}

	// Fast path, only one scorecard? No need for goroutines and channels.
	if len(scores) == 1 {
		play(ctx, holeID, langID, code, &scores[0])
		return scores[0]
	}

	done := make(chan Scorecard)

	for _, score := range scores {
		go func(score Scorecard) {
			defer func() {
				if r := recover(); r != nil {
					log.Println(r)
				}
			}()

			play(ctx, holeID, langID, code, &score)
			done <- score
		}(score)
	}

	// TODO Maybe return all runs (rather than last or failing) to the UI.
	for range scores {
		score = <-done

		// We failed! Return that run.
		if !score.Pass {
			break
		}
	}

	return // Return the last run.
}

func play(ctx context.Context, holeID, langID, code string, score *Scorecard) {
	var stderr, stdout bytes.Buffer
	var asmBytesRead, asmBytesWrite *os.File

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "/usr/bin/run-lang")
	cmd.Dir = "/langs/" + langID
	cmd.Env = []string{}
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	cmd.WaitDelay = time.Second
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWIPC | syscall.CLONE_NEWNET |
			syscall.CLONE_NEWNS | syscall.CLONE_NEWPID | syscall.CLONE_NEWUTS,
	}

	// Interpreter
	switch langID {
	case "assembly":
		var err error
		if asmBytesRead, asmBytesWrite, err = os.Pipe(); err != nil {
			panic(err)
		}

		cmd.Args = []string{"/usr/bin/defasm", "--size-out=3", "-w", "-r"}
		cmd.ExtraFiles = []*os.File{asmBytesWrite}
	case "awk":
		cmd.Args = []string{"/usr/bin/gawk", "-v", "RS=\\0", code}
	case "bash":
		cmd.Args = []string{"/usr/bin/bash", "-s", "-"}
	case "brainfuck":
		cmd.Args = []string{"/usr/bin/brainfuck", "-c", code}
	case "c":
		cmd.Args = []string{"/usr/bin/tcc", "-run", "-"}
	case "crystal":
		cmd.Args = []string{"/usr/bin/crystal", "run", "--stdin-filename", "code.cr", "--"}
		cmd.Env = []string{"CRYSTAL_CACHE_DIR=/tmp", "PATH=/usr/bin:/bin"}
	case "d":
		cmd.Args = []string{"/usr/bin/ldc2", "--enable-color=true", "--run", "-"}
		cmd.Env = []string{"PATH=/usr/bin"}
	case "elixir":
		cmd.Args = []string{"/usr/local/bin/elixir", "-e", code, "--"}
		cmd.Env = []string{"LANG=C.UTF-8", "PATH=/usr/local/bin:/usr/bin:/bin"}
	case "fish":
		cmd.Args = []string{"/usr/bin/fish", "--no-prng", "-c", code, "-u"}
	case "golfscript":
		cmd.Args = []string{"/usr/bin/golfscript", "-n", "-e", code}
		if holeID == "quine" {
			cmd.Args = append(cmd.Args, "-q")
		}
		cmd.Args = append(cmd.Args, "--")
	case "haskell", "php":
		cmd.Args = []string{"/usr/bin/" + langID, "--"}
	case "hexagony":
		cmd.Args = []string{"/usr/bin/hexagony", "-d", "-"}
	case "j":
		cmd.Args = []string{"/usr/bin/j", "/tmp/code.ijs"}
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
	case "powershell":
		cmd.Args = []string{"/usr/bin/powershell"}

		// Require explicit output for Quine to prevent trivial solutions.
		if holeID == "quine" {
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
		if holeID == "quine" {
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
	default:
		cmd.Args = []string{"/usr/bin/" + langID, "-"}
	}

	// Args
	switch langID {
	case "awk", "brainfuck", "fish":
		// Hole args passed through stdin for these langs separated by a null byte
		args := ""
		for _, arg := range score.Args {
			args += arg + "\x00"
		}
		cmd.Stdin = strings.NewReader(args)
	case "sed":
		// For sed we always need to append a null byte, even if no args exist
		args := strings.Join(score.Args, "\x00") + "\x00"
		cmd.Stdin = strings.NewReader(args)
	default:
		cmd.Args = append(cmd.Args, score.Args...)
	}

	// Code
	switch langID {
	case "awk", "brainfuck", "elixir", "fish", "golfscript", "javascript",
		"perl", "sed", "tex":
		// For these langs, code is passed as an argument above.
	case "k":
		cmd.Stdin = strings.NewReader(preprocessKCode(holeID, code))
	case "php":
		cmd.Stdin = strings.NewReader("<?php " + code + " ;")
	default:
		cmd.Stdin = strings.NewReader(code)
	}

	// We do not allow quine in TeX to pass without at least one backslash
	if langID == "tex" && holeID == "quine" && !strings.Contains(code, `\`) {
		// don't even run the code; just mark error and return
		score.Pass = false
		score.Stderr = []byte(`Quine in TeX must have at least one '\' character.`)
		return
	}

	err := cmd.Run()

	deadline, _ := ctx.Deadline()
	score.Took = timeout - time.Until(deadline)

	if err != nil {
		if err, ok := err.(*exec.ExitError); ok {
			score.ExitCode = err.ExitCode()
		}

		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			score.Timeout = true
			fmt.Fprint(&stderr, "Killed for exceeding the ", timeout, " timeout.")
		} else {
			stderr.WriteString(err.Error())
		}
	}

	// Actual byte count is printed by the assembler.
	if langID == "assembly" {
		if _, err := fmt.Fscanf(asmBytesRead, "%d", &score.ASMBytes); err != nil {
			panic(err)
		}
		asmBytesRead.Close()
	}

	const maxLength = 128 * 1024 // 128 KiB

	// Trim trailing whitespace.
	score.Stderr = bytes.TrimRightFunc(stderr.Next(maxLength), unicode.IsSpace)

	stdoutContents := stdout.Next(maxLength)

	// Postprocess sed output to turn null bytes into newlines.
	if langID == "sed" {
		stdoutContents = bytes.ReplaceAll(stdoutContents, []byte("\x00"), []byte("\n"))
	}

	// Trim trailing whitespace on each line, and then trailing empty lines.
	// Quine solutions are obviously left untouched.
	if holeID == "quine" {
		score.Stdout = stdoutContents
	} else {
		score.Stdout = bytes.TrimRight(stdoutTrimmer.ReplaceAll(
			stdoutContents, []byte{'\n'}), "\n")
	}

	// ASCII-ify roman numerals
	if holeID == "arabic-to-roman" {
		score.Stdout = []byte(romanToASCII.Replace(string(score.Stdout)))
	}

	// Timeouts do not pass, no matter what they output
	if score.Timeout {
		score.Pass = false
		return
	}

	// We do not allow stdout with only whitespace to pass to prevent suspicious sed "quines"
	if len(bytes.TrimRightFunc(score.Stdout, unicode.IsSpace)) != 0 {
		switch holeID {
		case "css-colors":
			// TODO Generalise case insensitivity, should it apply to others?
			score.Pass = strings.EqualFold(score.Answer, string(score.Stdout))
		default:
			score.Pass = score.Answer == string(score.Stdout)
		}
	}
}
