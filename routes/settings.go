package routes

import (
	"context"
	"net/http"

	"github.com/code-golf/code-golf/golfer"
	"github.com/code-golf/code-golf/zone"
)

// Settings serves GET /settings
func Settings(w http.ResponseWriter, r *http.Request) {
	info := golfer.GetInfo(db(r), "JRaspass")
	r = r.WithContext(context.WithValue(r.Context(), "golferInfo", info))

	render(w, r, http.StatusOK, "settings", "Settings", struct {
		Langs     []Lang
		TimeZones []zone.Zone
	}{langs, zone.List()})
}
