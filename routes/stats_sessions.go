package routes

import (
	"cmp"
	"fmt"
	"net/http"
	"slices"
	"strings"

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

	counts := map[string]int{}

	// Whitelist browser names to avoid displaying inappropriate text.
	for _, userAgent := range userAgents {
		switch name := userAgent.Name; name {
		case "Chrome", "curl", "Edge", "Firefox", "Opera", "Safari",
			"Samsung Browser":
			counts[name]++
		default:
			counts["Other"]++
		}
	}

	type browser struct {
		Count         int
		Name, Percent string
	}

	data := make([]browser, 0, len(counts))
	for name, count := range counts {
		data = append(data, browser{
			count, name,
			fmt.Sprintf("%.2f", float32(count)/float32(len(userAgents))*100),
		})
	}

	slices.SortFunc(data, func(a, b browser) int {
		return cmp.Or(
			cmp.Compare(b.Count, a.Count),
			cmp.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name)),
		)
	})

	render(w, r, "stats", data, "Statistics: Sessions")
}
