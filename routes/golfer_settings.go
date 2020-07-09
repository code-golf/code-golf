package routes

import (
	"net/http"

	"github.com/code-golf/code-golf/lang"
	"github.com/code-golf/code-golf/zone"
)

// GolferSettings serves GET /golfer/settings
func GolferSettings(w http.ResponseWriter, r *http.Request) {
	render(w, r, "golfer/settings", "Settings", struct {
		Langs     []lang.Lang
		TimeZones []zone.Zone
	}{lang.List, zone.List()})
}
