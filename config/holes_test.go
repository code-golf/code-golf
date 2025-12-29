package config

import (
	"strings"
	"testing"

	"github.com/pelletier/go-toml/v2"
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

func TestHoleReleasedExperimental(t *testing.T) {
	var holes []string
	var emptyLocalDate toml.LocalDate
	for _, hole := range ExpHoleList {
		if hole.Released != emptyLocalDate {
			holes = append(holes, hole.Name)
		}
	}

	assert.Truef(t, len(holes) == 0 || len(holes) == 1,
		"Up to one exp hole should have a released date, got = %v", holes)
}

func TestHoleReleasedStable(t *testing.T) {
	for _, hole := range HoleList {
		t.Run(hole.Name, func(t *testing.T) {
			assert.NotEmpty(t, hole.Released,
				"Stable holes must have a released date")
		})
	}
}
