package config

import (
	"html/template"
	"strings"
	"time"

	"github.com/pelletier/go-toml/v2"
	"golang.org/x/exp/slices"
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
	slices.SortFunc(CheevoList, func(a, b *Cheevo) bool {
		return strings.ToLower(a.Name) < strings.ToLower(b.Name)
	})
}
