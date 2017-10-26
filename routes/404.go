package routes

import "net/http"

func print404(w http.ResponseWriter, r *http.Request) {
	printHeader(w, r, 404)
	w.Write([]byte("<article><h1>404 Not Found"))
}
