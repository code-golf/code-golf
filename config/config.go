package config

import (
	"embed"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

//go:embed data/*
var config embed.FS

var id = strings.NewReplacer(
	// Emoji
	"🦄 ", "", "💎 ", "", "🥇 ", "",

	// Remove.
	"!", "", "(", "", ")", "", ",", "", "’", "", "?", "",

	// Hyphenate.
	" -", "-", " ", "-", ";", "-", "–", "-",

	// Special.
	"#", "-sharp",
	"+", "p",
	"><>", "fish",
)

// ID is a lowercase simplified version of the name used in URLs and the DB.
// Some characters (like "#") are changed to be safe in URLs.
func ID(name string) string { return strings.ToLower(id.Replace(name)) }

func unmarshal(file string, value any) {
	if data, err := config.ReadFile(file); err != nil {
		panic(err)
	} else if err := toml.Unmarshal(data, value); err != nil {
		panic(err)
	}
}
