package routes

import (
	"net/http"

	"github.com/code-golf/code-golf/session"
)

// GET /golfers/{golfer}/langs
func golferLangsGET(w http.ResponseWriter, r *http.Request) {
	golfer := session.GolferInfo(r).Golfer

	render(w, r, "golfer/langs", nil, golfer.Name)
}
