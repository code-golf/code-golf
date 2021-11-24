package config

import (
	"html/template"
	"sort"
	"strings"
)

var (
	CheevoByID = map[string]*Cheevo{}
	CheevoList []*Cheevo
	CheevoTree map[string][]*Cheevo
)

type Cheevo struct {
	Description template.HTML `json:"-"`
	Emoji       string        `json:"emoji"`
	ID          string        `json:"-"`
	Name        string        `json:"name"`
}

func init() {
	unmarshal("cheevos.toml", &CheevoTree)

	for _, categories := range CheevoTree {
		for _, cheevo := range categories {
			cheevo.ID = id(cheevo.Name)

			CheevoByID[cheevo.ID] = cheevo
			CheevoList = append(CheevoList, cheevo)
		}
	}

	// Case-insensitive sort.
	sort.Slice(CheevoList, func(i, j int) bool {
		return strings.ToLower(CheevoList[i].Name) <
			strings.ToLower(CheevoList[j].Name)
	})
}
