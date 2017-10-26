package routes

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func logout(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w.Header()["Set-Cookie"] = []string{"__Host-user=;MaxAge=0;Path=/;Secure"}
	w.Header().Set("Location", "/")

	w.WriteHeader(302)
}
