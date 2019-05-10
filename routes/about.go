package routes

import (
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func about(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	Render(w, r, http.StatusOK, "about", template.HTML(versionTable))
}
