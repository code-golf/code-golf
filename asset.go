package main

import (
	"io/ioutil"
	"net/http"
)

func asset(w http.ResponseWriter, r *http.Request) bool {
	// TODO Use a real asset pipeline.
	switch r.URL.Path {
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
		return false
	}

	return true
}
