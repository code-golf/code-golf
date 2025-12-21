package routes

import (
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/code-golf/code-golf/golfer"
)

func TestBanners(t *testing.T) {
	for _, tt := range []struct{ time, cheevo string }{
		{"2022-03-14T03:14:15Z", "<b>Pi Day</b> achievement will stop being available"},
		{"2022-10-24T13:33:37Z", "<b>Vampire Byte</b> achievement will be available"},
		{"2022-12-25T07:00:00Z", "<b>Twelvetide</b> achievement will stop being available"},
		{"2023-01-01T01:02:03Z", "<b>Twelvetide</b> achievement will stop being available"},
	} {
		time, _ := time.Parse(time.RFC3339, tt.time)

		if !slices.ContainsFunc(banners(&golfer.Golfer{}, time), func(b banner) bool {
			return strings.Contains(string(b.Body), tt.cheevo)
		}) {
			t.Errorf("banners(%v) didn't produce %v", tt.time, tt.cheevo)
		}
	}
}
