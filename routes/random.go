package routes

import (
	"math/rand"
	"net/http"

	"github.com/code-golf/code-golf/hole"
)

// Random serves GET /random
func Random(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, hole.List[rand.Intn(len(hole.List))].ID, http.StatusFound)
}

// NGRandom serves GET /ng/random
func NGRandom(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/ng/"+hole.List[rand.Intn(len(hole.List))].ID, http.StatusFound)
}
