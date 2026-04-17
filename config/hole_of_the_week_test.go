package config

import (
	"testing"
	"testing/synctest"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestThisWeek(t *testing.T) {
	for _, tt := range []struct{ now, thisWeek, msg string }{
		{"2026-01-12T00:00:00Z", "2026-01-12T00:00:00Z", "Earliest"},
		{"2026-01-12T13:33:37Z", "2026-01-12T00:00:00Z", "Monday"},
		{"2026-01-13T13:33:37Z", "2026-01-12T00:00:00Z", "Tuesday"},
		{"2026-01-14T13:33:37Z", "2026-01-12T00:00:00Z", "Wednesday"},
		{"2026-01-15T13:33:37Z", "2026-01-12T00:00:00Z", "Thursday"},
		{"2026-01-16T13:33:37Z", "2026-01-12T00:00:00Z", "Friday"},
		{"2026-01-17T13:33:37Z", "2026-01-12T00:00:00Z", "Saturday"},
		{"2026-01-18T13:33:37Z", "2026-01-12T00:00:00Z", "Sunday"},
		{"2026-01-18T23:59:59Z", "2026-01-12T00:00:00Z", "Latest"},
	} {
		synctest.Test(t, func(t *testing.T) {
			time.Sleep(time.Until(parseTime(tt.now)))

			assert.Equal(t, parseTime(tt.thisWeek), ThisWeek(), tt.msg)
		})
	}
}

func parseTime(s string) time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return t
}
