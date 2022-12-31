package config

import (
	"strings"
	"testing"

	"github.com/SeerUK/assert"
)

func TestHoleSynopses(t *testing.T) {
	for _, hole := range append(HoleList, ExpHoleList...) {
		t.Run(hole.Name, func(t *testing.T) {
			assert.True(t, hole.Synopsis != "", "Has a synopsis")

			assert.True(t, strings.HasSuffix(hole.Synopsis, "."),
				"Synopsis ends in a period")

			assert.False(t, strings.ContainsRune(hole.Synopsis, '\n'),
				"Synopsis spans exactly one line")

			assert.False(t, strings.ContainsAny(hole.Synopsis, `"'`),
				"Synopsis doesn't use ASCII quotes/apostrophes")
		})
	}
}
