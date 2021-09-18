package config

import (
	"sort"
	"strings"
)

var (
	LangByID = map[string]*Lang{}
	LangList []*Lang
)

type Lang struct {
	Example string `json:"example"`
	ID      string `json:"id"`
	Name    string `json:"name"`
	Size    string `json:"size"`
	Version string `json:"version"`
	Website string `json:"website"`
}

func init() {
	var langs map[string]*Lang
	unmarshal("langs.toml", &langs)

	for name, lang := range langs {
		lang.Example = strings.TrimSpace(lang.Example)
		lang.ID = id(name)
		lang.Name = name

		LangByID[lang.ID] = lang
		LangList = append(LangList, lang)
	}

	// Case-insensitive sort.
	sort.Slice(LangList, func(i, j int) bool {
		return strings.ToLower(LangList[i].Name) <
			strings.ToLower(LangList[j].Name)
	})
}
