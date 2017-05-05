package main

import (
	"bytes"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/html"
)

type handler struct{}

const base = "views/base.html"

var holes = make(map[string]*template.Template)
var index = template.Must(template.ParseFiles(base, "views/index.html"))
var minfy = minify.New()

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	println(r.Method, r.URL.Path)

	vars := map[string]interface{}{"r": r}

	if hole, ok := holes[r.URL.Path]; ok {
		if r.Method == "POST" {
			code := r.FormValue("code")

			var stderr, stdout bytes.Buffer

			// TODO Need better IDs
			// TODO Use github.com/opencontainers/runc/libcontainer
			cmd := exec.Cmd{
				Args:   []string{"runc", "start", "id"},
				Dir:    "containers/perl",
				Path:   "/usr/bin/runc",
				Stderr: &stderr,
				Stdin:  strings.NewReader(code),
				Stdout: &stdout,
			}

			if err := cmd.Run(); err != nil {
				println(err.Error())
				vars["died"] = true
			}

			vars["stderr"] = string(stderr.Bytes())
			vars["stdout"] = string(bytes.TrimSpace(stdout.Bytes()))

			// TODO Move this logic into the template?
			// vars["Pass"] = string(vars["stdout"]) == fizzBuzzAnswer
		}

		render(w, hole, vars)
		return
	}

	// TODO Use a real asset pipeline.
	switch r.URL.Path {
	case "/":
		render(w, index, vars)
	case "/css/codemirror.css":
		bytes, _ := ioutil.ReadFile("css/codemirror.css")
		w.Header().Set("Content-Type", "text/css")
		w.Write(bytes)
	case "/js/codemirror.js":
		bytes, _ := ioutil.ReadFile("js/codemirror.js")
		w.Header().Set("Content-Type", "application/javascript")
		w.Write(bytes)
	case "/js/codemirror-perl.js":
		bytes, _ := ioutil.ReadFile("js/codemirror-perl.js")
		w.Header().Set("Content-Type", "application/javascript")
		w.Write(bytes)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func render(w http.ResponseWriter, tmpl *template.Template, args interface{}) {
	w.Header().Set("Content-Type", "text/html;charset=utf8")

	pipeR, pipeW := io.Pipe()

	go func() {
		defer pipeW.Close()

		if err := tmpl.Execute(pipeW, args); err != nil {
			panic(err)
		}
	}()

	if err := minfy.Minify("text/html", w, pipeR); err != nil {
		panic(err)
	}
}

func main() {
	minfy.AddFunc("text/html", html.Minify)

	holes["/fizz-buzz"] = template.Must(
		template.ParseFiles(base, "views/fizz-buzz.html"))

	println("Listening…")

	panic(http.ListenAndServe(":80", &handler{}))
}
