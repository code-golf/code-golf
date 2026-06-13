package config

import (
	"cmp"
	"html/template"
	"slices"
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
	Holes       []*Hole          `json:"-"`
	ID          string           `json:"id"`
	Langs       []*Lang          `json:"-"`
	Name        string           `json:"name"`
	Target      int              `json:"-"`
	Times       []time.Time      `json:"-"`
}

func initCheevos() {
	unmarshal("data/cheevos.toml", &CheevoTree)

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

	// Dyanmic cheevo langs.
	var (
		fourByFour     = CheevoByID["4×4"]
		flagThoseMines = CheevoByID["flag-those-mines"]
		ringToss       = CheevoByID["ring-toss"]
	)
	for _, lang := range LangList {
		if len(lang.Name) == 4 {
			fourByFour.Langs = append(fourByFour.Langs, lang)
		}

		if strings.HasPrefix(lang.Name, "F") {
			flagThoseMines.Langs = append(flagThoseMines.Langs, lang)
		}

		if strings.ContainsRune(strings.ToLower(lang.Name), 'o') {
			ringToss.Langs = append(ringToss.Langs, lang)
		}
	}

	// Case-insensitive sort.
	slices.SortFunc(CheevoList, func(a, b *Cheevo) int {
		return cmp.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name))
	})
}
