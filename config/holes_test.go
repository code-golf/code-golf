package config

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHoleSynopses(t *testing.T) {
	for _, hole := range AllHoleList {
		t.Run(hole.Name, func(t *testing.T) {
			assert.NotEmpty(t, hole.Synopsis, "Has a synopsis")

			assert.True(t, strings.HasSuffix(hole.Synopsis, "."),
				"Synopsis ends in a period")

			assert.NotContains(t, hole.Synopsis, "\n",
				"Synopsis spans exactly one line")

			assert.False(t, strings.ContainsAny(hole.Synopsis, `"'`),
				"Synopsis doesn't use ASCII quotes/apostrophes")
		})
	}
}

func TestNextHolesReleased(t *testing.T) {
	for i := 1; i < len(NextHoles); i++ {
		if NextHoles[0].Released != NextHoles[i].Released {
			t.Error("All of NextHoles must have the same released date")
		}
	}
}

func TestHoleReleasedStable(t *testing.T) {
	for _, hole := range HoleList {
		t.Run(hole.Name, func(t *testing.T) {
			assert.NotEmpty(t, hole.Released,
				"Stable holes must have a released date")
		})
	}
}
