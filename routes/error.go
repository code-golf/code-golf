package routes

import "net/http"

// Forbidden serves a 403.
func Forbidden(w http.ResponseWriter, r *http.Request) {
	render(w, r, "403", nil, "403 Forbidden")
}

// NotFound serves a 404.
func NotFound(w http.ResponseWriter, r *http.Request) {
	render(w, r, "404", nil, "404 Not Found")
}
