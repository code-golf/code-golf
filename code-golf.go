package main

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"html/template"
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
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/html"
)

const base = "views/base.html"

var hmacKey []byte
var holeTmpl = template.Must(template.ParseFiles(base, "views/hole.html"))
var index = template.Must(template.ParseFiles(base, "views/index.html"))
var minfy = minify.New()

func init() {
	minfy.AddFunc("text/html", html.Minify)

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

func codeGolf(w http.ResponseWriter, r *http.Request) {
	header := w.Header()

	// Skip over the initial forward slash.
	switch path := r.URL.Path[1:]; path {
	case "":
		header["Strict-Transport-Security"] = []string{headerHSTS}

		vars := map[string]interface{}{"r": r}

		_, vars["login"] = readCookie(r)

		render(w, index, vars)
		printLeaderboards(w)
	case "callback":
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
	case "favicon.ico":
		w.Write(favicon)
	case "logout":
		header["Set-Cookie"] = []string{"__Host-user=;MaxAge=0;Path=/;Secure"}
		http.Redirect(w, r, "/", 302)
	case cssHash:
		header["Cache-Control"] = []string{"max-age=9999999,public"}
		header["Content-Type"] = []string{"text/css"}

		if strings.Contains(r.Header.Get("Accept-Encoding"), "br") {
			header["Content-Encoding"] = []string{"br"}
			w.Write(cssBr)
		} else {
			header["Content-Encoding"] = []string{"gzip"}
			w.Write(cssGzip)
		}
	case jsHash:
		header["Cache-Control"] = []string{"max-age=9999999,public"}
		header["Content-Type"] = []string{"application/javascript"}

		if strings.Contains(r.Header.Get("Accept-Encoding"), "br") {
			header["Content-Encoding"] = []string{"br"}
			w.Write(jsBr)
		} else {
			header["Content-Encoding"] = []string{"gzip"}
			w.Write(jsGzip)
		}
	default:
		var hole, lang string

		if i := strings.IndexByte(path, '/'); i != -1 {
			hole = path[:i]
			lang = path[i+1:]
		} else {
			hole = path
		}

		if preamble, ok := preambles[hole]; ok {
			switch lang {
			case "javascript", "perl", "perl6", "php", "python", "ruby":
				vars := map[string]interface{}{
					"lang":     lang,
					"preamble": template.HTML(preamble),
					"r":        r,
				}

				var userID int
				userID, vars["login"] = readCookie(r)

				if r.Method == http.MethodPost {
					code := strings.Replace(r.FormValue("code"), "\r", "", -1)
					vars["code"] = code

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
						vars["pass"] = true

						addSolution(userID, hole, lang, code)
					} else {
						var diffString string

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

						vars["diff"] = template.HTML(diffString)
					}
				} else if code := getSolutionCode(userID, hole, lang); code != "" {
					vars["code"] = code
				} else if hole == "fizz-buzz" {
					vars["code"] = examples[lang]
				}

				render(w, holeTmpl, vars)
			default:
				http.Redirect(w, r, "/"+hole+"/perl6", 302)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

func render(w http.ResponseWriter, tmpl *template.Template, vars map[string]interface{}) {
	header := w.Header()

	header["Content-Encoding"] = []string{"gzip"}
	header["Content-Type"] = []string{"text/html;charset=utf8"}
	header["Content-Security-Policy"] = []string{
		"default-src 'none';img-src data: https://avatars.githubusercontent.com https://code-golf.io/favicon.ico;script-src 'self';style-src 'self'",
	}

	vars["cssHash"] = cssHash
	vars["jsHash"] = jsHash

	pipeR, pipeW := io.Pipe()

	go func() {
		defer pipeW.Close()

		if err := tmpl.Execute(pipeW, vars); err != nil {
			panic(err)
		}
	}()

	var buf bytes.Buffer

	if err := minfy.Minify("text/html", &buf, pipeR); err != nil {
		panic(err)
	}

	writer := gzip.NewWriter(w)
	writer.Write(buf.Bytes())
	writer.Close()
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
