package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unicode"

	"github.com/buildkite/terminal"
)

const base = "views/base.html"

var hmacKey []byte

func init() {
	var err error
	if hmacKey, err = base64.RawURLEncoding.DecodeString(os.Getenv("HMAC_KEY")); err != nil {
		panic(err)
	}
}

func readCookie(r *http.Request) (id int, login string) {
	if cookie, err := r.Cookie("__Host-user"); err == nil {
		if i := strings.LastIndexByte(cookie.Value, ':'); i != -1 {
			mac := hmac.New(sha256.New, hmacKey)
			mac.Write([]byte(cookie.Value[:i]))

			if subtle.ConstantTimeCompare(
				[]byte(cookie.Value[i+1:]),
				[]byte(base64.RawURLEncoding.EncodeToString(mac.Sum(nil))),
			) == 1 {
				j := strings.IndexByte(cookie.Value, ':')

				id, _ = strconv.Atoi(cookie.Value[:j])
				login = cookie.Value[j+1 : i]
			}
		}
	}

	return
}

func printHeader(w io.WriteCloser, login string) {
	var logInOrOut string

	if login == "" {
		logInOrOut = `<a href="//github.com/login/oauth/authorize?client_id=7f6709819023e9215205&scope=user:email" id=login>Login with GitHub</a>`
	} else {
		logInOrOut = `<a href=/logout id=logout title=Logout></a><a href="/u/` + login + `" id=me><img src="//avatars.githubusercontent.com/` + login + `?s=30">` + login + "</a>"
	}

	w.Write([]byte(
		"<!doctype html>" +
			"<link rel=stylesheet href=" + cssPath + ">" +
			"<meta name=theme-color content=#222>" +
			`<meta name=viewport content="width=device-width">` +
			"<title>Code-Golf</title>" +
			"<header><nav>" +
			"<a href=/>Home</a>" +
			"<a href=/about>About</a>" +
			"<a href=/scores/all/all>Scores</a>" +
			logInOrOut +
			"</nav></header>",
	))
}

func codeGolf(w http.ResponseWriter, r *http.Request) {
	header := w.Header()
	path := r.URL.Path

	// Routes that don't return HTML.
	switch path {
	case "/callback":
		if user := githubAuth(r.FormValue("code")); user.ID != 0 {
			data := strconv.Itoa(user.ID) + ":" + user.Login

			mac := hmac.New(sha256.New, hmacKey)
			mac.Write([]byte(data))

			header["Set-Cookie"] = []string{
				"__Host-user=" + data + ":" +
					base64.RawURLEncoding.EncodeToString(mac.Sum(nil)) +
					";HttpOnly;Path=/;SameSite=Lax;Secure",
			}
		}

		http.Redirect(w, r, "/", 302)
		return
	case "/favicon.ico":
		w.Write(favicon)
		return
	case "/logout":
		header["Set-Cookie"] = []string{"__Host-user=;MaxAge=0;Path=/;Secure"}
		http.Redirect(w, r, "/", 302)
		return
	case "/solution":
		if r.Method == http.MethodPost {
			type In struct {
				Code, Hole, Lang string
			}

			type Out struct {
				Arg, Err, Exp, Out string
			}

			var in In

			if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
				panic(err)
			}
			defer r.Body.Close()

			var args []string
			var out Out

			if in.Hole == "arabic-to-roman-numerals" {
				for i := 0; i < 20; i++ {
					i := rand.Intn(3998) + 1 // 1 - 3999 inclusive.

					out.Exp += arabicToRoman(i) + "\n"
					args = append(args, strconv.Itoa(i))
				}

				// Drop the trailing newline.
				out.Exp = out.Exp[:len(out.Exp)-1]
			} else if in.Hole == "seven-segment" {
				args = make([]string, 1)
				args[0], out.Exp = sevenSegment()
			} else if in.Hole == "spelling-numbers" {
				args, out.Exp = spellingNumbers()
			} else {
				out.Exp = answers[in.Hole]
			}

			out.Err, out.Out = runCode(in.Lang, in.Code, args)
			out.Arg = strings.Join(args, " ")

			// Save the solution if it passes and the user is logged in.
			if out.Exp == out.Out {
				if userID, _ := readCookie(r); userID != 0 {
					addSolution(userID, in.Hole, in.Lang, in.Code)
				}
			}

			header["Content-Encoding"] = []string{"gzip"}
			header["Content-Type"] = []string{"application/json"}
			gzipWriter := gzip.NewWriter(w)
			defer gzipWriter.Close()

			enc := json.NewEncoder(gzipWriter)
			enc.SetEscapeHTML(false)

			if err := enc.Encode(&out); err != nil {
				panic(err)
			}
			return
		}
	case cssPath:
		header["Cache-Control"] = []string{"max-age=9999999,public"}
		header["Content-Type"] = []string{"text/css"}

		if strings.Contains(r.Header.Get("Accept-Encoding"), "br") {
			header["Content-Encoding"] = []string{"br"}
			w.Write(cssBr)
		} else {
			header["Content-Encoding"] = []string{"gzip"}
			w.Write(cssGz)
		}
		return
	case jsHolePath:
		header["Cache-Control"] = []string{"max-age=9999999,public"}
		header["Content-Type"] = []string{"application/javascript"}

		if strings.Contains(r.Header.Get("Accept-Encoding"), "br") {
			header["Content-Encoding"] = []string{"br"}
			w.Write(jsHoleBr)
		} else {
			header["Content-Encoding"] = []string{"gzip"}
			w.Write(jsHoleGz)
		}
		return
	case jsScoresPath:
		header["Cache-Control"] = []string{"max-age=9999999,public"}
		header["Content-Type"] = []string{"application/javascript"}

		if strings.Contains(r.Header.Get("Accept-Encoding"), "br") {
			header["Content-Encoding"] = []string{"br"}
			w.Write(jsScoresBr)
		} else {
			header["Content-Encoding"] = []string{"gzip"}
			w.Write(jsScoresGz)
		}
		return
	}

	header["Content-Encoding"] = []string{"gzip"}
	header["Content-Type"] = []string{"text/html;charset=utf8"}
	header["Content-Security-Policy"] = []string{
		"connect-src 'self';" +
			"default-src 'none';" +
			"img-src 'self' data: https://avatars.githubusercontent.com;" +
			"script-src 'self';" +
			"style-src 'self'",
	}

	userID, login := readCookie(r)

	gzipWriter := gzip.NewWriter(w)
	defer gzipWriter.Close()

	// Routes that do return HTML.
	if path == "/" {
		header["Strict-Transport-Security"] = []string{headerHSTS}

		printHeader(gzipWriter, login)
		printLeaderboards(gzipWriter, userID)
	} else if path == "/about" {
		printHeader(gzipWriter, login)
		gzipWriter.Write([]byte(about))
	} else if strings.HasPrefix(path, "/scores/") {
		parts := strings.Split(path[8:], "/")

		if len(parts) != 2 {
			http.Redirect(w, r, "/scores/all/all", 301)
			return
		}

		var hole, lang string

		switch parts[0] {
		case "all":
		case "99-bottles-of-beer",
			"arabic-to-roman-numerals",
			"e",
			"emirp-numbers",
			"evil-numbers",
			"fibonacci",
			"fizz-buzz",
			"odious-numbers",
			"pascals-triangle",
			"prime-numbers",
			"seven-segment",
			"sierpiński-triangle",
			"spelling-numbers",
			"π":
			hole = parts[0]
		default:
			http.Redirect(w, r, "/scores/all/all", 301)
			return
		}

		switch parts[1] {
		case "all":
		case "javascript", "perl", "perl6", "php", "python", "ruby":
			lang = parts[1]
		default:
			http.Redirect(w, r, "/scores/all/all", 301)
			return
		}

		printHeader(gzipWriter, login)
		printScores(gzipWriter, hole, lang, userID)
	} else if strings.HasPrefix(path, "/u/") && getUser(path[3:]) {
		printHeader(gzipWriter, login)
		gzipWriter.Write([]byte("<article><h1>" + path[3:] + "</h1>"))
	} else if preamble, ok := preambles[path[1:]]; ok {
		printHeader(gzipWriter, login)

		gzipWriter.Write([]byte(
			"<script async src=" + jsHolePath + "></script><div id=status><div>" +
				"<h2>Program Arguments</h2><pre id=Arg></pre>" +
				"<h2>Standard Error</h2><pre id=Err></pre>" +
				"<h2>Expected Output</h2><pre id=Exp></pre>" +
				"<h2>Standard Output</h2><pre id=Out></pre>" +
				"</div></div><article",
		))

		if userID == 0 {
			gzipWriter.Write([]byte(
				"><div id=alert>Please " +
					`<a href="//github.com/login/oauth/authorize?` +
					`client_id=7f6709819023e9215205&scope=user:email">` +
					"Login with GitHub</a> in order to save solutions " +
					"and appear on the leaderboards.</div",
			))
		} else {
			for lang, solution := range getUserSolutions(userID, path[1:]) {
				gzipWriter.Write([]byte(
					" data-" + lang + `="` +
						strings.Replace(solution, `"`, "&#34;", -1) + `"`,
				))
			}
		}

		gzipWriter.Write([]byte(
			">" + preamble +
				"<a class=tab href=#javascript>JS<span></span></a>" +
				"<a class=tab href=#perl>Perl<span></span></a>" +
				"<a class=tab href=#perl6>Perl 6<span></span></a>" +
				"<a class=tab href=#php>PHP<span></span></a>" +
				"<a class=tab href=#python>Python<span></span></a>" +
				"<a class=tab href=#ruby>Ruby<span></span></a>" +
				"<input type=submit value=Run>",
		))
	} else {
		w.WriteHeader(http.StatusNotFound)

		printHeader(gzipWriter, login)
		gzipWriter.Write([]byte("<article><h1>404 Not Found"))
	}
}

func runCode(lang, code string, args []string) (string, string) {
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
	default:
		cmd.Args = []string{"/usr/bin/" + lang, "-", "--"}
	}

	cmd.Args = append(cmd.Args, args...)

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	timer := time.AfterFunc(
		2*time.Second,
		func() {
			cmd.Process.Kill()
			err.WriteString("Killed for exceeding 2s timeout.")
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
	outBytes = bytes.TrimRightFunc(outBytes, unicode.IsSpace)

	// Escape HTML & convert ANSI to HTML in stderr.
	errBytes = terminal.Render(errBytes)

	// Escape HTML in stdout
	outBytes = bytes.Replace(outBytes, []byte{'<'}, []byte("&lt;"), -1)
	outBytes = bytes.Replace(outBytes, []byte{'>'}, []byte("&gt;"), -1)

	// ASCII-ify roman numerals
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

	return string(errBytes), string(outBytes)
}
