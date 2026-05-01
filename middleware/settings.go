package middleware

import (
	"maps"
	"net/http"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/session"
)

// Settings adds the page settings into the context.
func Settings(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := session.Get(r)

		// Build the settings map from config.Settings default values.
		s.Settings = make(map[string]map[string]any, len(config.Settings))
		for page, settings := range config.Settings {
			s.Settings[page] = make(map[string]any, len(settings))

			for _, setting := range settings {
				s.Settings[page][setting.ID] = setting.Default
			}
		}

		// If we have a golfer, merge their settings on top.
		if s.Golfer != nil {
			for page, settings := range s.Golfer.Settings {
				maps.Copy(s.Settings[page], settings)
			}
		}

		next.ServeHTTP(w, r)
	})
}
