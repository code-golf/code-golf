package routes

import "net/http"

// Robots serves GET /robots.txt
func Robots(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
