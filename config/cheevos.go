package config

import (
	"html/template"
	"sort"
	"strings"
	"time"

	"github.com/pelletier/go-toml/v2"
)

var (
	CheevoByID = map[string]*Cheevo{}
	CheevoList []*Cheevo
	CheevoTree map[string][]*Cheevo
)

type Cheevo struct {
	Dates       []toml.LocalDate `json:"-"`
	Description template.HTML    `json:"description"`
	Emoji       string           `json:"emoji"`
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Target      int              `json:"-"`
	Times       []time.Time      `json:"-"`
}

func init() {
	unmarshal("cheevos.toml", &CheevoTree)

	for _, categories := range CheevoTree {
		for _, cheevo := range categories {
			cheevo.ID = ID(cheevo.Name)

			for _, date := range cheevo.Dates {
				cheevo.Times = append(cheevo.Times, date.AsTime(time.UTC))
			}

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
