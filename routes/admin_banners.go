package routes

import (
	"net/http"
	"time"

	"github.com/code-golf/code-golf/golfer"
	"github.com/code-golf/code-golf/session"
)

// GET /admin/banners
func adminBannersGET(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Banners []struct {
			Banner string
			Count  int
		}
		LiveBanners map[string]bool
	}{LiveBanners: map[string]bool{}}

	if err := session.Database(r).Select(
		&data.Banners,
		` SELECT jsonb_object_keys(settings->'hide-banner') banner, COUNT(*)
		    FROM users
		GROUP BY banner
		ORDER BY banner`,
	); err != nil {
		panic(err)
	}

	// Use a dummy golfer to generate a list of live banners.
	for _, banner := range banners(&golfer.Golfer{}, time.Now().UTC()) {
		data.LiveBanners[banner.HideKey] = true
	}

	render(w, r, "admin/banners", data, "Admin Banners")
}

// POST /banners/{banner}
func adminBannerPOST(w http.ResponseWriter, r *http.Request) {
	// Delete banner from everyone's settings['hide-banner'][banner].
	// Where it's the only banner also delete settings['hide-banner'].
	session.Database(r).MustExec(
		`UPDATE users
		    SET settings = CASE
		   WHEN settings->'hide-banner' = jsonb_build_object($1, true)
		   THEN settings - 'hide-banner'
		   ELSE settings #- ARRAY['hide-banner', $1]
		    END
		  WHERE settings->'hide-banner' ? $1`,
		param(r, "banner"),
	)

	http.Redirect(w, r, "/admin/banners", http.StatusSeeOther)
}
