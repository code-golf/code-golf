package routes

import "net/http"

// NotFound serves any unknown route
func NotFound(w http.ResponseWriter, r *http.Request) {
	render(w, r, "404", "", nil)
}
