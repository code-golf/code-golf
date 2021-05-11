package cheevo

import (
	"html/template"
	"os"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
)

type Cheevo struct {
	Description     template.HTML
	Emoji, ID, Name string
}

var (
	ByID = map[string]*Cheevo{}
	List = []*Cheevo{}
	Tree = map[string][]*Cheevo{}
)

// TODO Share with other packages
func id(name string) string {
	name = strings.ReplaceAll(name, " ", "-")
	name = strings.ReplaceAll(name, "!", "")
	name = strings.ReplaceAll(name, ",", "")
	name = strings.ReplaceAll(name, ";", "-")
	name = strings.ReplaceAll(name, "’", "")

	return strings.ToLower(name)
}

func init() {
	// Tests run from the package directory, walk upward to find cheevos.toml.
	if _, err := os.Stat("cheevos.toml"); os.IsNotExist(err) {
		os.Chdir("..")
	}

	if _, err := toml.DecodeFile("cheevos.toml", &Tree); err != nil {
		panic(err)
	}

	for _, categories := range Tree {
		for _, cheevo := range categories {
			cheevo.ID = id(cheevo.Name)

			ByID[cheevo.ID] = cheevo
			List = append(List, cheevo)
		}
	}

	sort.Slice(List, func(i, j int) bool { return List[i].Name < List[j].Name })
}
