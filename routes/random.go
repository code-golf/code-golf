package routes

import (
	"math/rand"
	"net/http"

	"github.com/code-golf/code-golf/config"
)

// Random serves GET /random
func Random(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r,
		config.HoleList[rand.Intn(len(config.HoleList))].ID, http.StatusFound)
}

// NGRandom serves GET /ng/random
func NGRandom(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r,
		"/ng/"+config.HoleList[rand.Intn(len(config.HoleList))].ID, http.StatusFound)
}
