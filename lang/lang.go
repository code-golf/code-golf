package lang

import (
	"os"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
)

type Lang struct {
	Example string `json:"example"`
	ID      string `json:"id"`
	Name    string `json:"name"`
	Size    string `json:"size"`
	Version string `json:"version"`
	Website string `json:"website"`
}

var (
	ByID = map[string]Lang{}
	List []Lang
)

func init() {
	var langsTOML map[string]Lang

	// Tests run from the package directory, walk upward to find langs.toml.
	if _, err := os.Stat("langs.toml"); os.IsNotExist(err) {
		os.Chdir("..")
	}

	if _, err := toml.DecodeFile("langs.toml", &langsTOML); err != nil {
		panic(err)
	}

	for name, lang := range langsTOML {
		lang.Name = name
		lang.ID = strings.ReplaceAll(strings.ToLower(name), "#", "-sharp")
		lang.ID = strings.ReplaceAll(strings.ToLower(lang.ID), "><>", "fish")

		ByID[lang.ID] = lang
		List = append(List, lang)
	}

	// Case-insensitive sort.
	sort.Slice(List, func(i, j int) bool {
		return strings.ToLower(List[i].Name) < strings.ToLower(List[j].Name)
	})
}
