package main

import (
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
		logInOrOut = `<a href=/logout id=logout title=Logout></a><div><img src="//avatars.githubusercontent.com/` + login + `?size=30">` + login + "</div>"
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
			"<a href=/leaderboards>Leaderboards</a>" +
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
			gzipWriter := gzip.NewWriter(w)
			defer gzipWriter.Close()

			if err := json.NewEncoder(gzipWriter).Encode(&out); err != nil {
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
			w.Write(cssGzip)
		}
		return
	case jsPath:
		header["Cache-Control"] = []string{"max-age=9999999,public"}
		header["Content-Type"] = []string{"application/javascript"}

		if strings.Contains(r.Header.Get("Accept-Encoding"), "br") {
			header["Content-Encoding"] = []string{"br"}
			w.Write(jsBr)
		} else {
			header["Content-Encoding"] = []string{"gzip"}
			w.Write(jsGzip)
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
		printLeaderboards(gzipWriter)
	} else if path == "/about" {
		printHeader(gzipWriter, login)
		gzipWriter.Write([]byte(about))
	} else if strings.HasPrefix(path, "/u/") && getUser(path[3:]) {
		printHeader(gzipWriter, login)
		gzipWriter.Write([]byte("<article><h1>" + path[3:] + "</h1>"))
	} else if preamble, ok := preambles[path[1:]]; ok {
		printHeader(gzipWriter, login)

		gzipWriter.Write([]byte(
			"<script async src=" + jsPath + "></script><div id=status><div>" +
				"<h2>Program Arguments</h2>" +
				"<textarea id=Arg readonly rows=1></textarea>" +
				"<h2>Standard Error</h2>" +
				"<textarea id=Err readonly rows=1></textarea>" +
				"<h2>Expected Output</h2>" +
				"<textarea id=Exp readonly rows=5></textarea>" +
				"<h2>Standard Output</h2>" +
				"<textarea id=Out preamble rows=5></textarea>" +
				"</div></div><article",
		))

		if userID != 0 {
			for lang, solution := range getUserSolutions(userID, path[1:]) {
				gzipWriter.Write([]byte(
					" data-" + lang + `="` +
					strings.Replace(solution, `"`, "&#34;", -1) + `"`,
				))
			}
		}

		gzipWriter.Write([]byte(
			">" + preamble +
				"<a class=tab href=#javascript>JavaScript</a> " +
				"<a class=tab href=#perl>Perl</a> " +
				"<a class=tab href=#perl6>Perl 6</a> " +
				"<a class=tab href=#php>PHP</a> " +
				"<a class=tab href=#python>Python</a> " +
				"<a class=tab href=#ruby>Ruby</a>" +
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

	timer := time.AfterFunc(750*time.Millisecond, func() { cmd.Process.Kill() })

	if err := cmd.Wait(); err != nil {
		println(err.Error())
	}

	timer.Stop()

	return string(bytes.TrimRightFunc(err.Bytes(), unicode.IsSpace)),
		string(bytes.TrimRightFunc(out.Bytes(), unicode.IsSpace))
}
