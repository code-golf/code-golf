package hole

import (
	"bufio"
	"bytes"
	"context"
	"embed"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
	"unicode"
)

const timeout = 5 * time.Second

//go:embed answers
var answers embed.FS

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

		// Disable implicit output by inserting a ';' before all newlines, except when
		// the next line begins with a space (for a continuation).
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

func getAnswer(holeID, code string) (args []string, answer string) {
	args = []string{}

	switch holeID {
	case "arabic-to-roman", "roman-to-arabic":
		args, answer = arabicToRoman(holeID == "roman-to-arabic")
	case "arrows":
		args, answer = arrows()
	case "brainfuck":
		args, answer = brainfuck()
	case "css-colors":
		args, answer = cssColors()
	case "ellipse-perimeters":
		args, answer = ellipsePerimeters()
	case "emojify":
		args, answer = emojify()
	case "fractions":
		args, answer = fractions()
	case "hexdump":
		args, answer = hexdump()
	case "isbn":
		args, answer = isbn()
	case "intersection":
		args, answer = intersection()
	case "jacobi-symbol":
		args, answer = jacobiSymbol()
	case "levenshtein-distance":
		args, answer = levenshteinDistance()
	case "lucky-tickets":
		args, answer = luckyTickets()
	case "maze":
		args, answer = maze()
	case "morse-decoder", "morse-encoder":
		args, answer = morse(holeID == "morse-decoder")
	case "musical-chords":
		args, answer = musicalChords()
	case "ordinal-numbers":
		args, answer = ordinalNumbers()
	case "pangram-grep":
		args, answer = pangramGrep()
	case "poker":
		args, answer = poker()
	case "qr-decoder", "qr-encoder":
		args, answer = qr(holeID == "qr-decoder")
	case "quine":
		answer = code
	case "reverse-polish-notation":
		args, answer = reversePolishNotation()
	case "rock-paper-scissors-spock-lizard":
		args, answer = rockPaperScissorsSpockLizard()
	case "seven-segment":
		args, answer = sevenSegment()
	case "spelling-numbers":
		args, answer = spellingNumbers()
	case "star-wars-opening-crawl":
		args, answer = starWarsOpeningCrawl()
	case "sudoku", "sudoku-v2":
		args, answer = sudoku(holeID == "sudoku-v2")
	case "ten-pin-bowling":
		args, answer = tenPinBowling()
	case "time-distance":
		args, answer = timeDistance()
	case "united-states":
		args, answer = unitedStates()
	case "turtle":
		args, answer = turtle()
	case "zodiac-signs":
		args, answer = zodiacSigns()
	case "zeckendorf-representation":
		args, answer = zeckendorfRepresentation()
	default:
		// ¯\_(ツ)_/¯ cannot embed file answers/√2.txt: invalid name √2.txt
		if holeID == "√2" {
			holeID = "root-2"
		}

		if b, err := answers.ReadFile("answers/" + holeID + ".txt"); err != nil {
			panic(err)
		} else {
			answer = string(bytes.TrimSuffix(b, []byte{'\n'}))
		}
	}

	return
}

func Play(ctx context.Context, holeID, langID, code string) (score Scorecard) {
	score.Args, score.Answer = getAnswer(holeID, code)

	var stderr, stdout bytes.Buffer
	var asmBytesRead, asmBytesWrite *os.File

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "/usr/bin/run-lang")
	cmd.Dir = "/langs/" + langID
	cmd.Env = []string{}
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWIPC | syscall.CLONE_NEWNET | syscall.CLONE_NEWNS | syscall.CLONE_NEWPID | syscall.CLONE_NEWUTS,
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
	case "bash":
		cmd.Args = []string{"/usr/bin/bash", "-s", "-"}
	case "brainfuck":
		cmd.Args = []string{"/usr/bin/brainfuck", "-c", code}
	case "c":
		cmd.Args = []string{"/usr/bin/tcc", "-run", "-"}
	case "c-sharp", "f-sharp":
		cmd.Args = []string{"/compiler/Compiler", "-"}
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
	case "haskell", "php":
		cmd.Args = []string{"/usr/bin/" + langID, "--"}
	case "hexagony":
		cmd.Args = []string{"/hexagony/Hexagony", "-d", "-"}
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
	case "sed":
		cmd.Args = []string{"/usr/bin/sed", "-E", "-z", "--", code}
	case "swift":
		cmd.Args = []string{"/usr/bin/swift", "-module-cache-path", "/tmp", "-"}
	default:
		cmd.Args = []string{"/usr/bin/" + langID, "-"}
	}

	// Args
	switch langID {
	case "brainfuck", "fish":
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
	case "brainfuck", "elixir", "fish", "javascript", "sed":
		// For these code is passed as an argument above.
	case "k":
		code = preprocessKCode(holeID, code)
		cmd.Stdin = strings.NewReader(code)
	case "php":
		code = "<?php " + code + " ;"
		fallthrough
	default:
		cmd.Stdin = strings.NewReader(code)
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
			stderr.WriteString("Killed for exceeding the 5s timeout.")
		} else {
			stderr.WriteString(err.Error())
			println(err.Error())
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

	if holeID == "quine" {
		score.Stdout = stdout.Next(maxLength)
	} else {
		// Trim trailing spaces per line.
		// FIXME This is all very hacky, but needed for Sierpiński.
		stdoutContents := stdout.Next(maxLength)

		// Postprocess sed output to turn null bytes into newlines
		if langID == "sed" {
			stdoutContents = bytes.ReplaceAll(stdoutContents, []byte("\x00"), []byte("\n"))
		}

		scanner := bufio.NewScanner(bytes.NewReader(stdoutContents))
		for scanner.Scan() {
			score.Stdout = append(
				score.Stdout, bytes.TrimRightFunc(scanner.Bytes(), unicode.IsSpace)...)
			score.Stdout = append(score.Stdout, '\n')
		}

		if err := scanner.Err(); err != nil {
			panic(err)
		}

		score.Stdout = bytes.TrimRightFunc(score.Stdout, unicode.IsSpace)
	}

	// ASCII-ify roman numerals
	if holeID == "arabic-to-roman" {
		score.Stdout = bytes.ReplaceAll(score.Stdout, []byte("Ⅰ"), []byte("I"))
		score.Stdout = bytes.ReplaceAll(score.Stdout, []byte("Ⅱ"), []byte("II"))
		score.Stdout = bytes.ReplaceAll(score.Stdout, []byte("Ⅲ"), []byte("III"))
		score.Stdout = bytes.ReplaceAll(score.Stdout, []byte("Ⅳ"), []byte("IV"))
		score.Stdout = bytes.ReplaceAll(score.Stdout, []byte("Ⅴ"), []byte("V"))
		score.Stdout = bytes.ReplaceAll(score.Stdout, []byte("Ⅵ"), []byte("VI"))
		score.Stdout = bytes.ReplaceAll(score.Stdout, []byte("Ⅶ"), []byte("VII"))
		score.Stdout = bytes.ReplaceAll(score.Stdout, []byte("Ⅷ"), []byte("VIII"))
		score.Stdout = bytes.ReplaceAll(score.Stdout, []byte("Ⅸ"), []byte("IX"))
		score.Stdout = bytes.ReplaceAll(score.Stdout, []byte("Ⅹ"), []byte("X"))
		score.Stdout = bytes.ReplaceAll(score.Stdout, []byte("Ⅺ"), []byte("XI"))
		score.Stdout = bytes.ReplaceAll(score.Stdout, []byte("Ⅻ"), []byte("XII"))
		score.Stdout = bytes.ReplaceAll(score.Stdout, []byte("Ⅼ"), []byte("L"))
		score.Stdout = bytes.ReplaceAll(score.Stdout, []byte("Ⅽ"), []byte("C"))
		score.Stdout = bytes.ReplaceAll(score.Stdout, []byte("Ⅾ"), []byte("D"))
		score.Stdout = bytes.ReplaceAll(score.Stdout, []byte("Ⅿ"), []byte("M"))
	}

	// Timeouts do not pass, no matter what they output
	if score.Timeout {
		score.Pass = false
		return
	}

	if len(score.Stdout) != 0 {
		// TODO Generalise a case insensitive flag, should it apply to others?
		if holeID == "css-colors" {
			score.Pass = strings.EqualFold(score.Answer, string(score.Stdout))
		} else {
			score.Pass = score.Answer == string(score.Stdout)
		}
	}

	return
}
