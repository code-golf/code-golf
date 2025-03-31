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
