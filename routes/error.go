package routes

import "net/http"

// Forbidden serves a 403.
func Forbidden(w http.ResponseWriter, r *http.Request) {
	render(w, r, "403", "403 Forbidden", nil)
}

// NotFound serves a 404.
func NotFound(w http.ResponseWriter, r *http.Request) {
	render(w, r, "404", "404 Not Found", nil)
}
