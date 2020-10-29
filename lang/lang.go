package lang

import (
	"os"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
)

type Lang struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Website string
	Version string
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

	// Sort by ID rather than Name so the "sharps" are next to each other.
	sort.Slice(List, func(i, j int) bool { return List[i].ID < List[j].ID })
}
