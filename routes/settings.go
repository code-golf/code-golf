package routes

import (
	"context"
	"net/http"

	"github.com/code-golf/code-golf/golfer"
	"github.com/code-golf/code-golf/lang"
	"github.com/code-golf/code-golf/zone"
)

// Settings serves GET /settings
func Settings(w http.ResponseWriter, r *http.Request) {
	info := golfer.GetInfo(db(r), "JRaspass")
	r = r.WithContext(context.WithValue(r.Context(), "golferInfo", info))

	render(w, r, "settings", "Settings", struct {
		Langs     []lang.Lang
		TimeZones []zone.Zone
	}{lang.List, zone.List()})
}
