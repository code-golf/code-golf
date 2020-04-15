package hole

import (
	"bufio"
	"bytes"
	"context"
	"os/exec"
	"strings"
	"syscall"
	"time"
	"unicode"
)

const timeout = 7 * time.Second

type Scorecard struct {
	Answer         string
	Args           []string
	Pass, Timeout  bool
	Stderr, Stdout []byte
	Took           time.Duration
}

func getAnswer(holeID, code string) ([]string, string) {
	var answer string
	var args []string
	switch holeID {
	case "arabic-to-roman", "roman-to-arabic":
		args, answer = arabicToRoman(holeID == "roman-to-arabic")
	case "brainfuck":
		args, answer = brainfuck()
	case "lucky-tickets":
		args, answer = luckyTickets()
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
	case "sudoku":
		args, answer = sudoku()
	case "ten-pin-bowling":
		args, answer = tenPinBowling()
	case "united-states":
		args, answer = unitedStates()
	default:
		answer = answers[holeID]
	}
	return args, answer
}

func Play(ctx context.Context, holeID, langID, code string) (score Scorecard) {
	score.Args, score.Answer = getAnswer(holeID, code)

	var stderr, stdout bytes.Buffer

	if langID == "php" {
		code = "<?php " + code + " ;"
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "/usr/bin/run-lang")
	cmd.Dir = "langs/" + langID
	cmd.Stderr = &stderr
	cmd.Stdin = strings.NewReader(code)
	cmd.Stdout = &stdout
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWIPC | syscall.CLONE_NEWNET | syscall.CLONE_NEWNS | syscall.CLONE_NEWPID | syscall.CLONE_NEWUTS,
	}

	switch langID {
	case "bash":
		cmd.Args = []string{"/usr/bin/bash", "-s", "-"}
	case "c":
		cmd.Args = []string{"/usr/bin/tcc", "-run", "-"}
	case "go":
		cmd.Args = []string{"/usr/bin/run-go"}
	case "haskell", "javascript", "php":
		cmd.Args = []string{"/usr/bin/" + langID, "--"}
	case "j":
		cmd.Args = []string{"/usr/bin/j", "/tmp/code.ijs"}
	case "julia":
		cmd.Args = []string{"/usr/bin/run-julia", "/tmp/code.jl"}
	case "nim":
		cmd.Args = []string{"/usr/bin/run_nim"}
	default:
		cmd.Args = []string{"/usr/bin/" + langID, "-"}
	}

	cmd.Args = append(cmd.Args, score.Args...)

	err := cmd.Run()

	deadline, _ := ctx.Deadline()
	score.Took = timeout - time.Until(deadline)

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			score.Timeout = true
			stderr.WriteString("Killed for exceeding the 7s timeout.")
		} else {
			stderr.WriteString(err.Error())
			println(err.Error())
		}
	}

	// Trim trailing whitespace.
	score.Stderr = bytes.TrimRightFunc(stderr.Bytes(), unicode.IsSpace)

	if holeID == "quine" {
		score.Stdout = stdout.Bytes()
	} else {
		// Trim trailing spaces per line.
		// FIXME This is all very hacky, but needed for Sierpiński.
		scanner := bufio.NewScanner(bytes.NewReader(stdout.Bytes()))
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

	score.Pass = score.Answer == string(score.Stdout) && len(score.Stdout) != 0

	return
}
