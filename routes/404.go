package routes

import "net/http"

func print404(w http.ResponseWriter, r *http.Request) {
	printHeader(w, r, 404)
	w.Write([]byte("<main><h1>404 Not Found"))
}
