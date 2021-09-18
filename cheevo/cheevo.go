package cheevo

import (
	"html/template"
	"os"
	"sort"

	"github.com/BurntSushi/toml"
	"github.com/code-golf/code-golf/config"
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
			cheevo.ID = config.ID(cheevo.Name)

			ByID[cheevo.ID] = cheevo
			List = append(List, cheevo)
		}
	}

	sort.Slice(List, func(i, j int) bool { return List[i].Name < List[j].Name })
}
