package config

import (
	"embed"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

//go:embed *.toml
var tomls embed.FS

func id(name string) string {
	name = strings.ReplaceAll(name, " ", "-")
	name = strings.ReplaceAll(name, "!", "")
	name = strings.ReplaceAll(name, "#", "-sharp")
	name = strings.ReplaceAll(name, "(", "")
	name = strings.ReplaceAll(name, ")", "")
	name = strings.ReplaceAll(name, ",", "")
	name = strings.ReplaceAll(name, ";", "-")
	name = strings.ReplaceAll(name, "><>", "fish")
	name = strings.ReplaceAll(name, "’", "")

	return strings.ToLower(name)
}

func unmarshal(file string, value interface{}) {
	if data, err := tomls.ReadFile(file); err != nil {
		panic(err)
	} else if err := toml.Unmarshal(data, value); err != nil {
		panic(err)
	}
}
