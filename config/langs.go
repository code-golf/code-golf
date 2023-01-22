package config

import (
	"strings"

	"golang.org/x/exp/slices"
)

var (
	// Standard languages.
	LangByID = map[string]*Lang{}
	LangList []*Lang

	// Experimental languages.
	ExpLangByID = map[string]*Lang{}
	ExpLangList []*Lang

	// All languages.
	AllLangByID = map[string]*Lang{}
	AllLangList []*Lang
)

type Lang struct {
	Example    string `json:"example"`
	Experiment int    `json:"-"`
	ID         string `json:"id"`
	Name       string `json:"name"`
	Size       string `json:"size"`
	Version    string `json:"version"`
	Website    string `json:"website"`
}

func init() {
	var langs map[string]*Lang
	unmarshal("langs.toml", &langs)

	for name, lang := range langs {
		lang.Example = strings.TrimSuffix(lang.Example, "\n")
		lang.ID = ID(name)
		lang.Name = name

		AllLangByID[lang.ID] = lang
		AllLangList = append(AllLangList, lang)

		if lang.Experiment == 0 {
			LangByID[lang.ID] = lang
			LangList = append(LangList, lang)
		} else {
			ExpLangByID[lang.ID] = lang
			ExpLangList = append(ExpLangList, lang)
		}
	}

	// Case-insensitive sort.
	for _, langs := range [][]*Lang{LangList, ExpLangList, AllLangList} {
		slices.SortFunc(langs, func(a, b *Lang) bool {
			return strings.ToLower(a.Name) < strings.ToLower(b.Name)
		})
	}
}
