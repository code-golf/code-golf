package routes

import (
	"net/http"

	"github.com/code-golf/code-golf/session"
)

// Healthz serves GET /healthz.
func Healthz(w http.ResponseWriter, r *http.Request) {
	if err := session.Database(r).Ping(); err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusNoContent)
}
