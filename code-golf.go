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
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/html"
)

const base = "views/base.html"

var hmacKey []byte
var holes = make(map[string]*template.Template)
var index = template.Must(template.ParseFiles(base, "views/index.html"))
var minfy = minify.New()

func init() {
	minfy.AddFunc("text/html", html.Minify)

	holes["fizz-buzz"] = template.Must(
		template.ParseFiles(base, "views/fizz-buzz.html"))

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
	// Skip over the initial forward slash.
	switch path := r.URL.Path[1:]; path {
	case "":
		w.Header().Set("Link","</roboto-v16>;as=font;crossorigin;rel=preload")
		w.Header().Set("Strict-Transport-Security", headerHSTS)

		vars := map[string]interface{}{"r": r}

		_, vars["login"] = readCookie(r)

		render(w, index, vars)
	case "callback":
		if user := githubAuth(r.FormValue("code")); user.ID != 0 {
			data := strconv.Itoa(user.ID) + ":" + user.Login

			mac := hmac.New(sha256.New, hmacKey)
			mac.Write([]byte(data))

			w.Header().Set(
				"Set-Cookie",
				"__Host-user="+data+":"+
					base64.RawURLEncoding.EncodeToString(mac.Sum(nil))+
					";HttpOnly;Path=/;SameSite=Lax;Secure",
			)
		}

		http.Redirect(w, r, "/", 302)
	case "logout":
		w.Header().Set("Set-Cookie", "__Host-user=;MaxAge=0;Path=/;Secure")
		http.Redirect(w, r, "/", 302)
	case "roboto-v16":
		w.Header().Set("Cache-Control", "max-age=9999999,public")
		w.Header().Set("Content-Type", "font/woff2")
		w.Write(roboto)
	case "roboto-mono-v4":
		w.Header().Set("Cache-Control", "max-age=9999999,public")
		w.Header().Set("Content-Type", "font/woff2")
		w.Write(robotoMono)
	case cssHash:
		w.Header().Set("Cache-Control", "max-age=9999999,public")
		w.Header().Set("Content-Type", "text/css")

		if strings.Contains(r.Header.Get("Accept-Encoding"), "br") {
			w.Header().Set("Content-Encoding", "br")
			w.Write(cssBr)
		} else {
			w.Header().Set("Content-Encoding", "gzip")
			w.Write(cssGzip)
		}
	case jsHash:
		w.Header().Set("Cache-Control", "max-age=9999999,public")
		w.Header().Set("Content-Type", "application/javascript")

		if strings.Contains(r.Header.Get("Accept-Encoding"), "br") {
			w.Header().Set("Content-Encoding", "br")
			w.Write(jsBr)
		} else {
			w.Header().Set("Content-Encoding", "gzip")
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

		if tmpl, ok := holes[hole]; ok {
			switch lang {
			case "javascript", "perl", "perl6", "php", "python", "ruby":
				w.Header().Set("Link",
					"</roboto-v16>;as=font;crossorigin;rel=preload")
				w.Header().Set("Link",
					"</roboto-mono-v4>;as=font;crossorigin;rel=preload")

				vars := map[string]interface{}{"lang": lang, "r": r}

				_, vars["login"] = readCookie(r)

				if r.Method == http.MethodPost {
					vars["code"] = r.FormValue("code")
					vars["pass"] = fizzBuzzAnswer == runCode(lang, r.FormValue("code"))
				} else {
					vars["code"] = examples[lang]
				}

				render(w, tmpl, vars)
			default:
				http.Redirect(w, r, "/"+hole+"/perl", 302)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

func render(w http.ResponseWriter, tmpl *template.Template, vars map[string]interface{}) {
	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Set("Content-Type", "text/html;charset=utf8")
	w.Header().Set(
		"Content-Security-Policy",
		"default-src 'none';font-src 'self';img-src https://avatars.githubusercontent.com;script-src 'self';style-src 'self'",
	)

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

func runCode(lang, code string) string {
	var out bytes.Buffer

	// TODO Need better IDs
	// TODO Use github.com/opencontainers/runc/libcontainer
	cmd := exec.Cmd{
		Args:   []string{"runc", "start", "id"},
		Dir:    "containers/" + lang,
		Path:   "/usr/bin/runc",
		Stderr: os.Stdout,
		Stdin:  strings.NewReader(code),
		Stdout: &out,
	}

	if err := cmd.Run(); err != nil {
		println(err.Error())
	}

	return string(bytes.TrimSpace(out.Bytes()))
}
