package routes

import (
	"net/http"

	"github.com/code-golf/code-golf/config"
)

// GET /golfer/search
func golferSearchGET(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Langs map[string]string
		Holes map[string]string
	}{
		Langs: make(map[string]string),
		Holes: make(map[string]string),
	}

	for k, v := range config.AllLangByID {
		data.Langs[k] = v.Name
	}
	for k, v := range config.AllHoleByID {
		data.Holes[k] = v.Name
	}

	render(w, r, "golfer/code-search", data, "Solution search")
}
