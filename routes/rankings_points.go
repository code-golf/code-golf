package routes

import (
	"net/http"

	"github.com/code-golf/code-golf/hole"
	"github.com/code-golf/code-golf/lang"
)

// RankingsPoints serves GET /rankings/points
func RankingsPoints(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Holes []hole.Hole
		Langs []lang.Lang
	}{
		Holes: hole.List,
		Langs: lang.List,
	}

	render(w, r, "rankings/points", "Rankings: Points", data)
}
