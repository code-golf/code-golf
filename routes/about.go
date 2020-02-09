package routes

import (
	"html/template"
	"net/http"
)

// About serves GET /about
func About(w http.ResponseWriter, r *http.Request) {
	render(w, r, http.StatusOK, "about", "About", template.HTML(versionTable))
}
