package main

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"html"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/sergi/go-diff/diffmatchpatch"
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
			"<script async src=" + jsPath + "></script>" +
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

	// Routes that don't return HTML.
	switch r.URL.Path {
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
	case "/favicon.ico":
		w.Write(favicon)
	case "/logout":
		header["Set-Cookie"] = []string{"__Host-user=;MaxAge=0;Path=/;Secure"}
		http.Redirect(w, r, "/", 302)
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
	}

	header["Content-Encoding"] = []string{"gzip"}
	header["Content-Type"] = []string{"text/html;charset=utf8"}
	header["Content-Security-Policy"] = []string{
		"default-src 'none';img-src data: https://avatars.githubusercontent.com https://code-golf.io/favicon.ico;script-src 'self';style-src 'self'",
	}

	userID, login := readCookie(r)

	// Routes that do return HTML.
	switch r.URL.Path {
	case "/":
		header["Strict-Transport-Security"] = []string{headerHSTS}

		gzipWriter := gzip.NewWriter(w)
		defer gzipWriter.Close()

		printHeader(gzipWriter, login)
		printLeaderboards(gzipWriter)
	default:
		var hole, lang string

		// FIXME
		path := r.URL.Path[1:]

		if i := strings.IndexByte(path, '/'); i != -1 {
			hole = path[:i]
			lang = path[i+1:]
		} else {
			hole = path
		}

		if preamble, ok := preambles[hole]; ok {
			switch lang {
			case "javascript", "perl", "perl6", "php", "python", "ruby":
				var code, diffString string

				gzipWriter := gzip.NewWriter(w)
				defer gzipWriter.Close()

				printHeader(gzipWriter, login)

				if r.Method == http.MethodPost {
					code = strings.Replace(r.FormValue("code"), "\r", "", -1)

					var answer string
					var args []string

					if hole == "arabic-to-roman-numerals" {
						for i := 0; i < 20; i++ {
							i := rand.Intn(3998) + 1 // 1 - 3999 inclusive.

							answer += arabicToRoman(i) + "\n"
							args = append(args, strconv.Itoa(i))
						}

						// Drop the trailing newline.
						answer = answer[:len(answer)-1]
					} else {
						answer = answers[hole]
					}

					output := runCode(lang, code, args)

					if answer == output {
						gzipWriter.Write([]byte("<div id=success>Success</div>"))

						addSolution(userID, hole, lang, code)
					} else {
						gzipWriter.Write([]byte("<div id=failure>Failure</div>"))

						for _, diff := range diffmatchpatch.New().DiffMain(
							answer, output, false,
						) {
							switch diff.Type {
							case diffmatchpatch.DiffInsert:
								diffString += "<ins>" + diff.Text + "</ins>"
							case diffmatchpatch.DiffDelete:
								diffString += "<del>" + diff.Text + "</del>"
							case diffmatchpatch.DiffEqual:
								diffString += diff.Text
							}
						}

						diffString = "<pre>" + diffString + "</pre>"
					}
				} else {
					code = getSolutionCode(userID, hole, lang)
				}

				// TODO Only Fizz Buzz has example code ATM.
				if hole == "fizz-buzz" && code == ""  {
					code = examples[lang]
				}

				gzipWriter.Write([]byte(
					"<article>" + diffString + preamble +
						"<a class=lang href=javascript>JavaScript</a> " +
						"<a class=lang href=perl>Perl</a> " +
						"<a class=lang href=perl6>Perl 6</a> " +
						"<a class=lang href=php>PHP</a> " +
						"<a class=lang href=python>Python</a> " +
						"<a class=lang href=ruby>Ruby</a>" +
						"<form method=post>" +
						"<textarea name=code>" +
						html.EscapeString(code) +
						"</textarea>" +
						"<input type=submit value=Run>" +
						"</form>",
				))
			default:
				http.Redirect(w, r, "/"+hole+"/perl6", 302)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)

			gzipWriter := gzip.NewWriter(w)
			defer gzipWriter.Close()

			printHeader(gzipWriter, login)
			gzipWriter.Write([]byte("<article><h1>404 Not Found"))
		}
	}
}

func runCode(lang, code string, args []string) string {
	var out bytes.Buffer

	if lang == "php" {
		code = "<?php " + code + " ?>"
	}

	cmd := exec.Cmd{
		Dir:    "containers/" + lang,
		Path:   "../../run-container",
		Stderr: os.Stdout,
		Stdin:  strings.NewReader(code),
		Stdout: &out,
		SysProcAttr: &syscall.SysProcAttr{
			Cloneflags: syscall.CLONE_NEWIPC | syscall.CLONE_NEWNET | syscall.CLONE_NEWNS | syscall.CLONE_NEWPID | syscall.CLONE_NEWUTS,
		},
	}

	// $binary, @new_argv where $new_argv[0] is used for hostname too.
	switch lang {
	case "javascript":
		cmd.Args = []string{"/usr/bin/js", "javascript", "-f"}
	case "perl6":
		cmd.Args = []string{
			"/usr/bin/moar",
			"perl6",
			"--execname=perl6",
			"--libpath=/usr/share/nqp/lib",
			"--libpath=/usr/share/perl6/runtime",
			"/usr/share/perl6/runtime/perl6.moarvm",
			"-",
		}
	case "python":
		cmd.Args = []string{"/usr/bin/python3.6", "python"}
	default:
		cmd.Args = []string{"/usr/bin/" + lang, lang}
	}

	cmd.Args = append(cmd.Args, append([]string{"-"}, args...)...)

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	timer := time.AfterFunc(500*time.Millisecond, func() { cmd.Process.Kill() })

	if err := cmd.Wait(); err != nil {
		println(err.Error())
	}

	timer.Stop()

	return string(bytes.TrimSpace(out.Bytes()))
}
