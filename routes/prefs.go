package routes

import (
	"net/http"

	"github.com/code-golf/code-golf/zone"
)

// Prefs serves GET /prefs
func Prefs(w http.ResponseWriter, r *http.Request) {
	render(w, r, http.StatusOK, "prefs", "Preferences", zone.List())
}
