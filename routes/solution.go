package routes

import (
	"bufio"
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unicode"

	"github.com/buildkite/terminal"
	"github.com/jraspass/code-golf/cookie"
	"github.com/julienschmidt/httprouter"
	"github.com/pmezard/go-difflib/difflib"
)

func solution(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	type In struct {
		Code, Hole, Lang string
	}

	type Out struct {
		Arg, Diff, Err, Exp, Out string
	}

	var in In

	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		panic(err)
	}
	defer r.Body.Close()

	var args []string
	var out Out

	switch in.Hole {
	case "arabic-to-roman":
		for i := 0; i < 20; i++ {
			i := rand.Intn(3998) + 1 // 1 - 3999 inclusive.

			out.Exp += arabicToRoman(i) + "\n"
			args = append(args, strconv.Itoa(i))
		}

		// Drop the trailing newline.
		out.Exp = out.Exp[:len(out.Exp)-1]
	case "quine":
		out.Exp = in.Code
	case "roman-to-arabic":
		for i := 0; i < 20; i++ {
			i := rand.Intn(3998) + 1 // 1 - 3999 inclusive.

			out.Exp += strconv.Itoa(i) + "\n"
			args = append(args, arabicToRoman(i))
		}

		// Drop the trailing newline.
		out.Exp = out.Exp[:len(out.Exp)-1]
	case "seven-segment":
		args = make([]string, 1)
		args[0], out.Exp = sevenSegment()
	case "spelling-numbers":
		args, out.Exp = spellingNumbers()
	default:
		out.Exp = answers[in.Hole]
	}

	out.Err, out.Out = runCode(in.Hole, in.Lang, in.Code, args)
	out.Arg = strings.Join(args, " ")

	out.Diff, _ = difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
		A:        difflib.SplitLines(out.Exp),
		B:        difflib.SplitLines(out.Out),
		Context:  3,
		FromFile: "Exp",
		ToFile:   "Out",
	})

	// Save the solution if the user is logged in and it passes.
	if userID, _ := cookie.Read(r); userID != 0 && out.Exp == out.Out && out.Out != "" {
		// Update the code if it's the same length or less, but only update
		// the submitted time if the solution is shorter. This avoids a user
		// moving down the leaderboard by matching their personal best.
		if _, err := db.Exec(`
		    INSERT INTO solutions
		         VALUES (NOW(), $1, $2, $3, $4)
		    ON CONFLICT ON CONSTRAINT solutions_pkey
		  DO UPDATE SET submitted = CASE
		                    WHEN LENGTH($4) < LENGTH(solutions.code)
		                    THEN NOW()
		                    ELSE solutions.submitted
		                END,
		                code = CASE
		                    WHEN LENGTH($4) > LENGTH(solutions.code)
		                    THEN solutions.code
		                    ELSE $4
		                END
		`, userID, in.Hole, in.Lang, in.Code); err != nil {
			panic(err)
		}
	}

	w.Header()["Content-Type"] = []string{"application/json"}

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)

	if err := enc.Encode(&out); err != nil {
		panic(err)
	}
}

func runCode(hole, lang, code string, args []string) (string, string) {
	var err, out bytes.Buffer

	if lang == "php" {
		code = "<?php " + code + " ?>"
	}

	cmd := exec.Cmd{
		Dir:    "containers/" + lang,
		Path:   "../../run-container",
		Stderr: &err,
		Stdin:  strings.NewReader(code),
		Stdout: &out,
		SysProcAttr: &syscall.SysProcAttr{
			Cloneflags: syscall.CLONE_NEWIPC | syscall.CLONE_NEWNET | syscall.CLONE_NEWNS | syscall.CLONE_NEWPID | syscall.CLONE_NEWUTS,
		},
	}

	switch lang {
	case "javascript":
		cmd.Args = []string{"/usr/bin/js", "-f", "-", "--"}
	// The perls seems to show -- in @ARGV :-(
	case "perl":
		cmd.Args = []string{"/usr/bin/perl", "-"}
	case "perl6":
		cmd.Args = []string{
			"/usr/bin/moar",
			"--execname=perl6",
			"--libpath=/usr/share/nqp/lib",
			"--libpath=/usr/share/perl6/runtime",
			"/usr/share/perl6/runtime/perl6.moarvm",
			"-",
		}
	case "php":
		cmd.Args = []string{"/usr/bin/php", "--"}
	// Ruby also seem to show -- in ARGV :-(
	// TODO Try special casing just the langs that need --
	case "ruby":
		cmd.Args = []string{"/usr/bin/ruby", "-"}
	default:
		cmd.Args = []string{"/usr/bin/" + lang, "-", "--"}
	}

	cmd.Args = append(cmd.Args, args...)

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	timer := time.AfterFunc(
		5*time.Second,
		func() {
			cmd.Process.Kill()
			err.WriteString("Killed for exceeding the 5s timeout.")
		},
	)

	if err := cmd.Wait(); err != nil {
		println(err.Error())
	}

	timer.Stop()

	var outBytes []byte

	// Trim trailing spaces per line.
	// FIXME This is all very hacky, but needed for Sierpiński.
	scanner := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	for scanner.Scan() {
		outBytes = append(outBytes, bytes.TrimRightFunc(scanner.Bytes(), unicode.IsSpace)...)
		outBytes = append(outBytes, '\n')
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	// Trim trailing whitespace.
	errBytes := bytes.TrimRightFunc(err.Bytes(), unicode.IsSpace)

	if hole != "quine" {
		outBytes = bytes.TrimRightFunc(outBytes, unicode.IsSpace)
	}

	// Escape HTML & convert ANSI to HTML in stderr.
	errBytes = terminal.Render(errBytes)

	// ASCII-ify roman numerals
	if hole == "arabic-to-roman" {
		outBytes = bytes.Replace(outBytes, []byte("Ⅰ"), []byte("I"), -1)
		outBytes = bytes.Replace(outBytes, []byte("Ⅱ"), []byte("II"), -1)
		outBytes = bytes.Replace(outBytes, []byte("Ⅲ"), []byte("III"), -1)
		outBytes = bytes.Replace(outBytes, []byte("Ⅳ"), []byte("IV"), -1)
		outBytes = bytes.Replace(outBytes, []byte("Ⅴ"), []byte("V"), -1)
		outBytes = bytes.Replace(outBytes, []byte("Ⅵ"), []byte("VI"), -1)
		outBytes = bytes.Replace(outBytes, []byte("Ⅶ"), []byte("VII"), -1)
		outBytes = bytes.Replace(outBytes, []byte("Ⅷ"), []byte("VIII"), -1)
		outBytes = bytes.Replace(outBytes, []byte("Ⅸ"), []byte("IX"), -1)
		outBytes = bytes.Replace(outBytes, []byte("Ⅹ"), []byte("X"), -1)
		outBytes = bytes.Replace(outBytes, []byte("Ⅺ"), []byte("XI"), -1)
		outBytes = bytes.Replace(outBytes, []byte("Ⅻ"), []byte("XII"), -1)
		outBytes = bytes.Replace(outBytes, []byte("Ⅼ"), []byte("L"), -1)
		outBytes = bytes.Replace(outBytes, []byte("Ⅽ"), []byte("C"), -1)
		outBytes = bytes.Replace(outBytes, []byte("Ⅾ"), []byte("D"), -1)
		outBytes = bytes.Replace(outBytes, []byte("Ⅿ"), []byte("M"), -1)
	}

	return string(errBytes), string(outBytes)
}
