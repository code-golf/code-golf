package routes

import (
	"net/http"

	"github.com/code-golf/code-golf/zone"
)

// Settings serves GET /settings
func Settings(w http.ResponseWriter, r *http.Request) {
	render(w, r, http.StatusOK, "settings", "Settings", struct {
		Langs     []Lang
		TimeZones []zone.Zone
	}{langs, zone.List()})
}
