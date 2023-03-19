package routes

import (
	"log"
	"net/http"

	"github.com/code-golf/code-golf/session"
)

// GET /healthz.
func healthzGET(w http.ResponseWriter, r *http.Request) {
	if err := session.Database(r).Ping(); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}
