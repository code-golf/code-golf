package routes

import "net/http"

// GET /
func homeGET(w http.ResponseWriter, r *http.Request) {
	render(w, r, "home", getHomeCards(r))
}
