package main

import (
	"bytes"
	"crypto/tls"
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
		if r.Method == http.MethodPost {
			code := r.FormValue("code")

			var out bytes.Buffer

			// TODO Need better IDs
			// TODO Use github.com/opencontainers/runc/libcontainer
			cmd := exec.Cmd{
				Args:   []string{"runc", "start", "id"},
				Dir:    "containers/perl",
				Path:   "/usr/bin/runc",
				Stdin:  strings.NewReader(code),
				Stdout: &out,
			}

			if err := cmd.Run(); err != nil {
				println(err.Error())
			}

			vars["code"] = code
			vars["pass"] = fizzBuzzAnswer == string(bytes.TrimSpace(out.Bytes()))
		} else {
			vars["code"] = fizzBuzzExample
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

func mustLoadX509KeyPair(certFile, keyFile string) tls.Certificate {
	if cert, err := tls.LoadX509KeyPair(certFile, keyFile); err != nil {
		panic(err)
	} else {
		return cert
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

	server := &http.Server{
		Handler: &handler{},
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{
				mustLoadX509KeyPair(
					"/home/jraspass/dehydrated/certs/code-golf.io/fullchain.pem",
					"/home/jraspass/dehydrated/certs/code-golf.io/privkey.pem",
				),
				mustLoadX509KeyPair(
					"/home/jraspass/dehydrated/certs/raspass.me/fullchain.pem",
					"/home/jraspass/dehydrated/certs/raspass.me/privkey.pem",
				),
			},
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256, // HTTP/2-required.
			},
			CurvePreferences:         []tls.CurveID{tls.CurveP256, tls.X25519},
			MinVersion:               tls.VersionTLS12,
			PreferServerCipherSuites: true,
		},
	}

	server.TLSConfig.BuildNameToCertificate()

	println("Listening…")

	// Redirect HTTP to HTTPS.
	go func() {
		panic(http.ListenAndServe(
			":80",
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(
					w,
					r,
					"https://"+r.Host+r.URL.String(),
					http.StatusMovedPermanently,
				)
			}),
		))
	}()

	// Serve HTTPS.
	panic(server.ListenAndServeTLS("", ""))
}
