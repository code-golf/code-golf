package routes

import (
	"math/rand"
	"net/http"
)

// Random serves GET /random
func Random(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, holes[rand.Intn(len(holes))].ID, http.StatusFound)
}
