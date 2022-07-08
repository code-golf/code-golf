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
	Description template.HTML `json:"description"`
	Emoji       string        `json:"emoji"`
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Target      int           `json:"-"`
}

func init() {
	unmarshal("cheevos.toml", &CheevoTree)

	for _, categories := range CheevoTree {
		for _, cheevo := range categories {
			cheevo.ID = ID(cheevo.Name)

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
