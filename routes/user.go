package routes

import (
	"net/http"
)

// User serves GET /users/{name}
func User(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/golfers/"+param(r, "name"), http.StatusPermanentRedirect)
}
