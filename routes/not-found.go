package routes

import "net/http"

// NotFound serves any unknown route
func NotFound(w http.ResponseWriter, r *http.Request) {
	render(w, r, http.StatusNotFound, "404", "", nil)
}
