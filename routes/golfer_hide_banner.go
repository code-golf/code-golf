package routes

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/session"
)

var (
	cheevoBannerRegex  = regexp.MustCompile(`^cheevo-(?:before|until)-\d{4}-\d{2}-\d{2}-`)
	holeOfTheWeekRegex = regexp.MustCompile(`^hole-of-the-week-\d{4}-\d{2}-\d{2}$`)
)

// POST /golfer/{action:hide|restore}-banner
func golferHideRestoreBannerPOST(w http.ResponseWriter, r *http.Request) {
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
	} else {
		valid = holeOfTheWeekRegex.MatchString(banner)
	}

	if valid {
		golfer := session.Golfer(r)

		const key = "hide-banner"
		if param(r, "action") == "hide" {
			if _, ok := golfer.Settings[key]; !ok {
				golfer.Settings[key] = map[string]any{}
			}
			golfer.Settings[key][banner] = true
		} else {
			delete(golfer.Settings[key], banner)
			if len(golfer.Settings[key]) == 0 {
				delete(golfer.Settings, key)
			}
		}

		golfer.SaveSettings(session.Database(r))
	}

	http.Redirect(w, r, r.FormValue("path"), http.StatusFound)
}
