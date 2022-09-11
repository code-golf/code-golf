package routes

import (
	"math/rand"
	"net/http"

	"github.com/code-golf/code-golf/config"
)

// GET /random
func randomGET(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r,
		config.HoleList[rand.Intn(len(config.HoleList))].ID, http.StatusFound)
}
