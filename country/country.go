package country

import (
	"os"

	"github.com/BurntSushi/toml"
)

type Country struct{ ID, Flag, Name string }

var (
	ByID = map[string]*Country{}
	Tree = map[string][]*Country{}
)

func init() {
	// Tests run from the package directory, walk upward to find countries.toml.
	if _, err := os.Stat("countries.toml"); os.IsNotExist(err) {
		os.Chdir("..")
	}

	if _, err := toml.DecodeFile("countries.toml", &Tree); err != nil {
		panic(err)
	}

	for _, countries := range Tree {
		for _, country := range countries {
			ByID[country.ID] = country

			for _, letter := range country.ID {
				country.Flag += string('🇦' - 'A' + letter)
			}
		}
	}
}
