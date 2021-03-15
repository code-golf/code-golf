package routes

import (
	"net/http"

	"github.com/code-golf/code-golf/hole"
	"github.com/code-golf/code-golf/lang"
)

// RankingsHoles serves GET /rankings/holes/{hole}/{lang}/{scoring}
func RankingsHoles(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Holes []hole.Hole
		Langs []lang.Lang
	}{
		Holes: hole.List,
		Langs: lang.List,
	}

	render(w, r, "rankings/holes", "Rankings: Holes", data)
}
