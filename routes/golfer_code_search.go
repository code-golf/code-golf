package routes

import (
	"net/http"
)

// GET /golfer/search
func golferSearchGET(w http.ResponseWriter, r *http.Request) {
	render(w, r, "golfer/code-search", nil, "Solution search")
}
