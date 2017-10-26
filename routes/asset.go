package routes

import (
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func asset(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	switch r.URL.Path {
	case "/favicon.ico":
		w.Write(favicon)
		return
	case cssPath:
		w.Header()["Cache-Control"] = []string{"max-age=9999999,public"}
		w.Header()["Content-Type"] = []string{"text/css"}

		if strings.Contains(r.Header.Get("Accept-Encoding"), "br") {
			w.Header()["Content-Encoding"] = []string{"br"}
			w.Write(cssBr)
		} else {
			w.Header()["Content-Encoding"] = []string{"gzip"}
			w.Write(cssGz)
		}
		return
	case jsHolePath:
		w.Header()["Cache-Control"] = []string{"max-age=9999999,public"}
		w.Header()["Content-Type"] = []string{"application/javascript"}

		if strings.Contains(r.Header.Get("Accept-Encoding"), "br") {
			w.Header()["Content-Encoding"] = []string{"br"}
			w.Write(jsHoleBr)
		} else {
			w.Header()["Content-Encoding"] = []string{"gzip"}
			w.Write(jsHoleGz)
		}
		return
	case jsScoresPath:
		w.Header()["Cache-Control"] = []string{"max-age=9999999,public"}
		w.Header()["Content-Type"] = []string{"application/javascript"}

		if strings.Contains(r.Header.Get("Accept-Encoding"), "br") {
			w.Header()["Content-Encoding"] = []string{"br"}
			w.Write(jsScoresBr)
		} else {
			w.Header()["Content-Encoding"] = []string{"gzip"}
			w.Write(jsScoresGz)
		}
		return
	}
}
