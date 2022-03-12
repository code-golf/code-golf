package routes

import "net/http"

// GET /log-out
func logOutGET(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		HttpOnly: true,
		MaxAge:   -1,
		Name:     "__Host-session",
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		Secure:   true,
	})

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
