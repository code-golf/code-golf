package routes

import (
	"net/http"

	"github.com/code-golf/code-golf/config"
)

// GET /sandbox
func sandboxGET(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Authors      []string
		HideDetails  bool
		Hole         *config.Hole
		Langs        map[string]*config.Lang
		RankingsView string
		Solutions    []map[string]string
	}{
		Langs:       config.AllLangByID,
		Hole:        &config.Hole {
			Name: "Sandbox",
			ID: "sandbox",
		},
	}

	render(w, r, "hole-tabs", data, "Sandbox", "XYZ123")
}
