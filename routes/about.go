package routes

import (
	"html/template"
	"net/http"
)

// About serves GET /about
func About(w http.ResponseWriter, r *http.Request) {
	render(w, r, "about", "About", template.HTML(versionTable))
}
