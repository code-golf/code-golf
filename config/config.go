package config

import "strings"

// ID converts a name into a URL-safe ID.
func ID(name string) string {
	name = strings.ReplaceAll(name, " ", "-")
	name = strings.ReplaceAll(name, "!", "")
	name = strings.ReplaceAll(name, "#", "-sharp")
	name = strings.ReplaceAll(name, ",", "")
	name = strings.ReplaceAll(name, ";", "-")
	name = strings.ReplaceAll(name, "><>", "fish")
	name = strings.ReplaceAll(name, "’", "")

	return strings.ToLower(name)
}
