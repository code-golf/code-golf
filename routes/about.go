package routes

import (
	"net/http"

	"github.com/code-golf/code-golf/lang"
)

// About serves GET /about
func About(w http.ResponseWriter, r *http.Request) {
	render(w, r, "about", "About", lang.List)
}
