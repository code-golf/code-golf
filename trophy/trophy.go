package trophy

import (
	"html/template"
	"os"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
)

type Trophy struct {
	Description     template.HTML
	Emoji, ID, Name string
}

var (
	ByID = map[string]*Trophy{}
	List = []*Trophy{}
	Tree = map[string][]*Trophy{}
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
	// Tests run from the package directory, walk upward to find trophies.toml.
	if _, err := os.Stat("trophies.toml"); os.IsNotExist(err) {
		os.Chdir("..")
	}

	if _, err := toml.DecodeFile("trophies.toml", &Tree); err != nil {
		panic(err)
	}

	for _, categories := range Tree {
		for _, trophy := range categories {
			trophy.ID = id(trophy.Name)

			ByID[trophy.ID] = trophy
			List = append(List, trophy)
		}
	}

	sort.Slice(List, func(i, j int) bool { return List[i].Name < List[j].Name })
}
