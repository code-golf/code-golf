package routes

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/session"
)

var cheevoBannerRegex = regexp.MustCompile(`^cheevo-(?:before|until)-\d{4}-\d{2}-\d{2}-`)

// POST /golfer/hide-banner
func golferHideBannerPOST(w http.ResponseWriter, r *http.Request) {
	banner := r.FormValue("banner")
	valid := false

	// Validate the banner is of a known form.
	if holeID, ok := strings.CutPrefix(banner, "latest-hole-"); ok {
		_, valid = config.HoleByID[holeID]
	} else if langID, ok := strings.CutPrefix(banner, "latest-lang-"); ok {
		_, valid = config.LangByID[langID]
	} else if holeID, ok := strings.CutPrefix(banner, "upcoming-hole-"); ok {
		_, valid = config.ExpHoleByID[holeID]
	} else if prefix := cheevoBannerRegex.FindString(banner); prefix != "" {
		_, valid = config.CheevoByID[banner[len(prefix):]]
	}

	if valid {
		golfer := session.Golfer(r)

		if _, ok := golfer.Settings["hide-banner"]; !ok {
			golfer.Settings["hide-banner"] = map[string]any{}
		}
		golfer.Settings["hide-banner"][banner] = true

		golfer.SaveSettings(session.Database(r))
	}

	http.Redirect(w, r, r.FormValue("path"), http.StatusFound)
}
