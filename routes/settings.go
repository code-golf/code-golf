package routes

import (
	"net/http"

	"github.com/code-golf/code-golf/lang"
	"github.com/code-golf/code-golf/zone"
)

// Settings serves GET /settings
func Settings(w http.ResponseWriter, r *http.Request) {
	render(w, r, "settings", "Settings", struct {
		Langs     []lang.Lang
		TimeZones []zone.Zone
	}{lang.List, zone.List()})
}
