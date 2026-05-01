package routes

import (
	"slices"
	"strings"
	"testing"
	"testing/synctest"
	"time"

	"github.com/code-golf/code-golf/golfer"
)

func TestBanners(t *testing.T) {
	for _, tt := range []struct{ time, cheevo string }{
		{"2022-03-14T03:14:15Z", "<b>Pi Day</b> achievement will stop being available"},
		{"2022-10-24T13:33:37Z", "<b>Vampire Byte</b> achievement will be available"},
		{"2022-12-25T07:00:00Z", "<b>Twelvetide</b> achievement will stop being available"},
		{"2023-01-01T01:02:03Z", "<b>Twelvetide</b> achievement will stop being available"},

		// 5 mins until New Years.
		{"2025-12-31T23:55:00Z", "<b>Twelvetide</b> achievement will stop being available"},
		{"2025-12-31T23:55:00Z", "The 🍃 <b>Turn over a New Leaf</b> achievement will be available"},
	} {
		synctest.Test(t, func(t *testing.T) {
			now, _ := time.Parse(time.RFC3339, tt.time)
			time.Sleep(time.Until(now))

			if !slices.ContainsFunc(banners(&golfer.Golfer{}), func(b banner) bool {
				return strings.Contains(string(b.Body), tt.cheevo)
			}) {
				t.Errorf("banners(%v) didn't produce %v", tt.time, tt.cheevo)
			}
		})
	}
}
