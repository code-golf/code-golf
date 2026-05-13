package routes

import (
	"cmp"
	"fmt"
	"maps"
	"net/http"
	"slices"
	"strings"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/session"
)

// GET /stats/{page:sessions}
func statsSessionsGET(w http.ResponseWriter, r *http.Request) {
	var userAgents []ua

	if err := session.Database(r).Select(
		&userAgents, "SELECT browser FROM sessions",
	); err != nil {
		panic(err)
	}

	type browser struct {
		Count             int
		ID, Name, Percent string
		Versions          map[string]int
	}
	browsers := map[string]*browser{}

	for _, userAgent := range userAgents {
		name, version := userAgent.Name, ""

		// Whitelist browser names to avoid displaying inappropriate text.
		switch name {
		case "Chrome", "Edge", "Firefox", "Opera":
			version = fmt.Sprint(userAgent.VersionNo.Major)
		case "Safari":
			version = fmt.Sprintf("%d.%d",
				userAgent.VersionNo.Major, userAgent.VersionNo.Minor)
		case "Samsung Browser", "curl":
		default:
			name = "Other"
		}

		b, ok := browsers[name]
		if !ok {
			b = &browser{
				ID:       config.ID(name),
				Name:     name,
				Versions: map[string]int{},
			}
			browsers[name] = b
		}

		b.Count++

		if version != "" {
			b.Versions[version]++
		}
	}

	data := slices.Collect(maps.Values(browsers))

	for _, browser := range data {
		browser.Percent = fmt.Sprintf("%.2f",
			float32(browser.Count)/float32(len(userAgents))*100)
	}

	slices.SortFunc(data, func(a, b *browser) int {
		return cmp.Or(
			cmp.Compare(b.Count, a.Count),
			cmp.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name)),
		)
	})

	render(w, r, "stats", data, "Statistics: Sessions")
}
