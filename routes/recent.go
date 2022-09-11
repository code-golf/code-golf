package routes

import "net/http"

// GET /recent/{lang}
func recentGET(w http.ResponseWriter, r *http.Request) {
	langID := param(r, "lang")
	if langID == "all-langs" {
		langID = "all"
	}

	url := "/recent/solutions/all/" + langID + "/bytes"
	http.Redirect(w, r, url, http.StatusPermanentRedirect)
}
