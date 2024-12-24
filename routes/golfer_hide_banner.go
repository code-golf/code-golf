package routes

import (
	"net/http"
	"strings"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/session"
)

// POST /golfer/hide-banner
func golferHideBannerPOST(w http.ResponseWriter, r *http.Request) {
	// banner = hide_key
	banner := r.FormValue("banner")
	valid := false

	// Validate the banner is of a known form.
	// cf. routes/banner.go for where the keys are created.
	if holeID, ok := strings.CutPrefix(banner, "latest-hole-"); ok {
		_, valid = config.HoleByID[holeID]
	} else if holeID, ok := strings.CutPrefix(banner, "upcoming-hole-"); ok {
		_, valid = config.ExpHoleByID[holeID]
	}

	if valid {
		golfer := session.Golfer(r)

		if _, err := session.Database(r).Exec(
			`INSERT INTO hidden_banners (hide_key, user_id)
			      VALUES                ($1,       $2)
			 ON CONFLICT DO NOTHING`,
			banner,
			golfer.ID,
		); err != nil {
			panic(err)
		}
	}

	http.Redirect(w, r, r.FormValue("path"), http.StatusFound)
}
