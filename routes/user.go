package routes

import "net/http"

// GET /users/{name}
func userGET(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/golfers/"+param(r, "name"), http.StatusPermanentRedirect)
}
