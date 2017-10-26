package routes

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func about(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	printHeader(w, r, 200)

	const html = "<article><h1>About</h1>" +
		"<p>Code Golf is written in <a href=//golang.org>Go</a> " +
		"and is <a href=//github.com/JRaspass/code-golf>open source</a>, " +
		"patches welcome!</p><h2>Versions</h2><table>" + versionTable

	w.Write([]byte(html))
}
