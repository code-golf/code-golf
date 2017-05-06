package main

import (
	"io/ioutil"
	"net/http"
)

func raspass(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		bytes, _ := ioutil.ReadFile("raspass.me/index.html")
		w.Header().Set("Content-Type", "text/html;charset=utf8")
		w.Write(bytes)
	case "/bg.png":
		bytes, _ := ioutil.ReadFile("raspass.me/bg.png")
		w.Header().Set("Content-Type", "image/png")
		w.Write(bytes)
	case "/favicon.ico":
		bytes, _ := ioutil.ReadFile("raspass.me/favicon.ico")
		w.Header().Set("Content-Type", "image/x-icon")
		w.Write(bytes)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}
