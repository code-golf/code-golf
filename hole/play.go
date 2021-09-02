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

const timeout = 7 * time.Second

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
	case "intersection":
		args, answer = intersection()
	case "levenshtein-distance":
		args, answer = levenshteinDistance()
	case "lucky-tickets":
		args, answer = luckyTickets()
	case "maze":
		args, answer = maze()
	case "morse-decoder", "morse-encoder":
		args, answer = morse(holeID == "morse-decoder")
	case "ordinal-numbers":
		args, answer = ordinalNumbers()
	case "pangram-grep":
		args, answer = pangramGrep()
	case "poker":
		args, answer = poker()
	case "quine":
		answer = code
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
	case "united-states":
		args, answer = unitedStates()
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

		cmd.Args = []string{"/usr/bin/defasm", "--size-out=3", "-r"}
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
	case "fish":
		cmd.Args = []string{"/usr/bin/fish", "--no-prng", "-c", code, "-u"}
	case "haskell", "php":
		cmd.Args = []string{"/usr/bin/" + langID, "--"}
	case "hexagony":
		cmd.Args = []string{"/hexagony/Hexagony", "-d", "-"}
	case "j":
		cmd.Args = []string{"/usr/bin/j", "/tmp/code.ijs"}
	case "javascript":
		cmd.Args = []string{"/usr/bin/d8", "-e", code, "--"}
	case "julia":
		cmd.Args = []string{"/usr/bin/run-julia", "/tmp/code.jl"}
	case "nim":
		cmd.Args = []string{"/usr/bin/nim", "-o:/tmp/code", "-r", "c", "-"}
	case "powershell":
		cmd.Args = []string{"/interpreter/Interpreter", "-"}

		// Require explicit output for Quine to prevent trivial solutions.
		if holeID == "quine" {
			cmd.Args[1] = "--explicit"
			cmd.Args = append(cmd.Args, "-")
		}
	case "python":
		// Force the stdout and stderr streams to be unbuffered.
		cmd.Args = []string{"/usr/bin/python", "-u", "-"}
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
	default:
		cmd.Args = append(cmd.Args, score.Args...)
	}

	// Code
	switch langID {
	case "brainfuck", "fish", "javascript":
		// For these code is passed as an argument above.
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
			stderr.WriteString("Killed for exceeding the 7s timeout.")
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
		scanner := bufio.NewScanner(bytes.NewReader(stdout.Next(maxLength)))
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
		score.Stdout = bytes.Replace(score.Stdout, []byte("Ⅰ"), []byte("I"), -1)
		score.Stdout = bytes.Replace(score.Stdout, []byte("Ⅱ"), []byte("II"), -1)
		score.Stdout = bytes.Replace(score.Stdout, []byte("Ⅲ"), []byte("III"), -1)
		score.Stdout = bytes.Replace(score.Stdout, []byte("Ⅳ"), []byte("IV"), -1)
		score.Stdout = bytes.Replace(score.Stdout, []byte("Ⅴ"), []byte("V"), -1)
		score.Stdout = bytes.Replace(score.Stdout, []byte("Ⅵ"), []byte("VI"), -1)
		score.Stdout = bytes.Replace(score.Stdout, []byte("Ⅶ"), []byte("VII"), -1)
		score.Stdout = bytes.Replace(score.Stdout, []byte("Ⅷ"), []byte("VIII"), -1)
		score.Stdout = bytes.Replace(score.Stdout, []byte("Ⅸ"), []byte("IX"), -1)
		score.Stdout = bytes.Replace(score.Stdout, []byte("Ⅹ"), []byte("X"), -1)
		score.Stdout = bytes.Replace(score.Stdout, []byte("Ⅺ"), []byte("XI"), -1)
		score.Stdout = bytes.Replace(score.Stdout, []byte("Ⅻ"), []byte("XII"), -1)
		score.Stdout = bytes.Replace(score.Stdout, []byte("Ⅼ"), []byte("L"), -1)
		score.Stdout = bytes.Replace(score.Stdout, []byte("Ⅽ"), []byte("C"), -1)
		score.Stdout = bytes.Replace(score.Stdout, []byte("Ⅾ"), []byte("D"), -1)
		score.Stdout = bytes.Replace(score.Stdout, []byte("Ⅿ"), []byte("M"), -1)
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
