package routes

import (
	"net/http"
)

// Golfer serves GET /golfers/{golfer}
func Golfer(w http.ResponseWriter, r *http.Request) {
	render(w, r, http.StatusOK, "golfer", "JRaspass", nil)
}

// GolferHoles serves GET /golfers/{golfer}/holes
func GolferHoles(w http.ResponseWriter, r *http.Request) {
	render(w, r, http.StatusOK, "golfer-holes", "JRaspass", langs)
}
