package main

import (
	"bytes"
	"html/template"
	"io"
	"io/ioutil"
	"os/exec"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/html"
	"github.com/valyala/fasthttp"
)

const base = "views/base.html"

var holes = make(map[string]*template.Template)
var index = template.Must(template.ParseFiles(base, "views/index.html"))
var minfy = minify.New()

func handler(ctx *fasthttp.RequestCtx) {
	// FIXME Use a real asset pipeline.
	if string(ctx.Path()) == "/css/codemirror.css" {
		bytes, _ := ioutil.ReadFile("css/codemirror.css")
		ctx.SetContentType("text/css")
		ctx.SetBody(bytes)
		return
	} else if string(ctx.Path()) == "/js/codemirror.js" {
		bytes, _ := ioutil.ReadFile("js/codemirror.js")
		ctx.SetContentType("application/javascript")
		ctx.SetBody(bytes)
		return
	} else if string(ctx.Path()) == "/js/codemirror-perl.js" {
		bytes, _ := ioutil.ReadFile("js/codemirror-perl.js")
		ctx.SetContentType("application/javascript")
		ctx.SetBody(bytes)
		return
	}

	vars := map[string]interface{}{"ctx": ctx}

	if len(ctx.Path()) == 1 {
		execTemplate(ctx, index, vars)
	} else if hole, ok := holes[string(ctx.Path())]; ok {
		if ctx.IsPost() {
			code := ctx.PostArgs().Peek("code")

			var stderr, stdout bytes.Buffer

			cmd := &exec.Cmd{
				Args:   []string{"perl"},
				Path:   "/usr/bin/perl",
				Stderr: &stderr,
				Stdin:  bytes.NewReader(code),
				Stdout: &stdout,
			}

			if err := cmd.Run(); err != nil {
				vars["died"] = true
			}

			vars["stderr"] = string(stderr.Bytes())
			vars["stdout"] = string(bytes.TrimSpace(stdout.Bytes()))

			// TODO Move this logic into the template?
			// vars["Pass"] = string(vars["stdout"]) == fizzBuzzAnswer
		}

		execTemplate(ctx, hole, vars)
	} else {
		ctx.NotFound()
	}
}

func execTemplate(ctx *fasthttp.RequestCtx, tmpl *template.Template, args interface{}) {
	ctx.SetContentType("text/html;charset=utf8")

	r, w := io.Pipe()

	go func() {
		defer w.Close()

		if err := tmpl.Execute(w, args); err != nil {
			panic(err)
		}
	}()

	if err := minfy.Minify("text/html", ctx, r); err != nil {
		panic(err)
	}
}

func main() {
	minfy.AddFunc("text/html", html.Minify)

	holes["/fizz-buzz"] = template.Must(
		template.ParseFiles(base, "views/fizz-buzz.html"))

	println("Listening...")

	fasthttp.ListenAndServe(":1337", handler)
}
