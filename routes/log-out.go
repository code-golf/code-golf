package routes

import "net/http"

// LogOut serves GET /log-out
func LogOut(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Set-Cookie", "__Host-user=;MaxAge=0;Path=/;Secure")

	http.Redirect(w, r, "/", http.StatusFound)
}
