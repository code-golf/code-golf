package routes

import (
	"bufio"
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"os/exec"
	"strings"
	"syscall"
	"time"
	"unicode"

	"github.com/JRaspass/code-golf/cookie"
	"github.com/buildkite/terminal"
	"github.com/julienschmidt/httprouter"
	"github.com/pmezard/go-difflib/difflib"
)

func solution(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var in struct {
		Code, Hole, Lang string
	}

	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		panic(err)
	}
	defer r.Body.Close()

	println(in.Code)

	var args []string
	var out struct {
		Arg, Diff, Err, Exp, Out string
		Argv                     []string
	}

	switch in.Hole {
	case "arabic-to-roman":
		args, out.Exp = arabicToRoman(false)
	case "brainfuck":
		args, out.Exp = brainfuck()
	case "morse-decoder", "morse-encoder":
		args, out.Exp = morse(in.Hole == "morse-decoder")
	case "pangram-grep":
		args, out.Exp = pangramGrep()
	case "poker":
		args, out.Exp = poker()
	case "quine":
		out.Exp = in.Code
	case "roman-to-arabic":
		args, out.Exp = arabicToRoman(true)
	case "seven-segment":
		args = make([]string, 1)
		args[0], out.Exp = sevenSegment()
	case "spelling-numbers":
		args, out.Exp = spellingNumbers()
	case "sudoku":
		args, out.Exp = sudoku()
	default:
		out.Exp = answers[in.Hole]
	}

	out.Err, out.Out = runCode(in.Hole, in.Lang, in.Code, args)
	out.Arg = strings.Join(args, " ")
	out.Argv = args

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
		         VALUES (NOW() AT TIME ZONE 'UTC', $1, $2, $3, $4)
		    ON CONFLICT ON CONSTRAINT solutions_pkey
		  DO UPDATE SET failing = false,
		                submitted = CASE
		                    WHEN solutions.failing OR LENGTH($4) < LENGTH(solutions.code)
		                    THEN NOW() AT TIME ZONE 'UTC'
		                    ELSE solutions.submitted
		                END,
		                code = CASE
		                    WHEN LENGTH($4) > LENGTH(solutions.code) AND NOT solutions.failing
		                    THEN solutions.code
		                    ELSE $4
		                END
		`, userID, in.Hole, in.Lang, in.Code); err != nil {
			panic(err)
		}

		awardTrophy(db, userID, "hello-world")

		if in.Hole == "fizz-buzz" {
			awardTrophy(db, userID, "interview-ready")
		}

		if queryBool(
			db,
			`SELECT COUNT(DISTINCT hole) > 18
			   FROM solutions
			  WHERE NOT failing AND user_id = $1`,
			userID,
		) {
			awardTrophy(db, userID, "the-watering-hole")
		}

		if queryBool(
			db,
			`SELECT COUNT(DISTINCT lang) = ARRAY_LENGTH(ENUM_RANGE(NULL::lang), 1)
			   FROM solutions
			  WHERE NOT failing AND user_id = $1`,
			userID,
		) {
			awardTrophy(db, userID, "polyglot")
		}

		if queryBool(
			db,
			"SELECT points > 9000 FROM points WHERE user_id = $1",
			userID,
		) {
			awardTrophy(db, userID, "its-over-9000")
		}

		switch in.Lang {
		case "php":
			awardTrophy(db, userID, "elephpant-in-the-room")
		case "perl", "perl6":
			if queryBool(
				db,
				`SELECT COUNT(*) = 2
				   FROM solutions
				  WHERE NOT failing
				    AND hole = $1
				    AND lang IN ('perl', 'perl6')
				    AND user_id = $2`,
				in.Hole,
				userID,
			) {
				awardTrophy(db, userID, "tim-toady")
			}
		case "python":
			if in.Hole == "quine" {
				awardTrophy(db, userID, "ouroboros")
			}
		}

	}

	w.Header()["Content-Type"] = []string{"application/json"}

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)

	if err := enc.Encode(&out); err != nil {
		panic(err)
	}
}

func awardTrophy(db *sql.DB, userID int, trophy string) {
	if _, err := db.Exec(
		"INSERT INTO trophies VALUES(NOW() AT TIME ZONE 'UTC', $1, $2) ON CONFLICT DO NOTHING",
		userID,
		trophy,
	); err != nil {
		panic(err)
	}
}

func queryBool(db *sql.DB, query string, args ...interface{}) (b bool) {
	if err := db.QueryRow(query, args...).Scan(&b); err != nil {
		panic(err)
	}

	return
}

func runCode(hole, lang, code string, args []string) (string, string) {
	var stderr, stdout bytes.Buffer

	if lang == "php" {
		code = "<?php " + code + " ?>"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "../../run-container")
	cmd.Dir = "containers/" + lang
	cmd.Stderr = &stderr
	cmd.Stdin = strings.NewReader(code)
	cmd.Stdout = &stdout
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWIPC | syscall.CLONE_NEWNET | syscall.CLONE_NEWNS | syscall.CLONE_NEWPID | syscall.CLONE_NEWUTS,
	}

	switch lang {
	case "bash":
		cmd.Args = []string{"/usr/bin/bash", "-s", "-"}
	case "haskell", "javascript", "php":
		cmd.Args = []string{"/usr/bin/" + lang, "--"}
	case "j":
		cmd.Args = []string{"/usr/bin/j", "/tmp/code.ijs"}
	case "julia":
		cmd.Args = []string{"/usr/bin/run-julia", "/tmp/code.jl"}
	case "perl6":
		cmd.Args = []string{
			"/usr/bin/moar",
			"--execname=perl6",
			"--libpath=/usr/share/nqp/lib",
			"--libpath=/usr/share/perl6/runtime",
			"/usr/share/perl6/runtime/perl6.moarvm",
			"-",
		}
	// Lua, Perl, Python, and Ruby are all sane.
	default:
		cmd.Args = []string{"/usr/bin/" + lang, "-"}
	}

	cmd.Args = append(cmd.Args, args...)

	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			stderr.WriteString("Killed for exceeding the 7s timeout.")
		} else {
			println(err.Error())
		}
	}

	var outBytes []byte

	// Trim trailing spaces per line.
	// FIXME This is all very hacky, but needed for Sierpiński.
	scanner := bufio.NewScanner(bytes.NewReader(stdout.Bytes()))
	for scanner.Scan() {
		outBytes = append(outBytes, bytes.TrimRightFunc(scanner.Bytes(), unicode.IsSpace)...)
		outBytes = append(outBytes, '\n')
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	// Trim trailing whitespace.
	errBytes := bytes.TrimRightFunc(stderr.Bytes(), unicode.IsSpace)

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
