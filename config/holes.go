package config

import (
	"bytes"
	"cmp"
	"database/sql"
	"embed"
	"encoding/json"
	"html/template"
	"reflect"
	"slices"
	"strings"
	templateTxt "text/template"

	"github.com/code-golf/code-golf/ordered"
	"github.com/lib/pq/hstore"
	"github.com/pelletier/go-toml/v2"
)

//go:embed hole-answers
var answers embed.FS

var (
	// Standard holes.
	HoleByID = map[string]*Hole{}
	HoleList []*Hole

	// Experimental holes.
	ExpHoleByID = map[string]*Hole{}
	ExpHoleList []*Hole

	// All holes.
	AllHoleByID = map[string]*Hole{}
	AllHoleList []*Hole

	// Ten most recent holes, used for /rankings/recent-holes.
	RecentHoles []*Hole

	// A map of hole ID to category for passing to SQL queries.
	HoleCategoryHstore = hstore.Hstore{Map: map[string]sql.NullString{}}

	// Aliases & Redirects
	HoleAliases   = map[string]string{}
	HoleRedirects = map[string]string{}
)

type Link struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Hole struct {
	Aliases, Redirects          []string       `json:"-"`
	Answer                      string         `json:"-"`
	AnswerFunc                  HoleAnswerFunc `json:"-"`
	CaseFold                    bool           `json:"-" toml:"case-fold"`
	Category                    string         `json:"category"`
	CategoryColor, CategoryIcon string         `json:"-"`
	Data                        template.JS    `json:"-"`
	DataMap                     ordered.Map    `json:"-"`
	Experiment                  int            `json:"experiment,omitzero"`
	ID                          string         `json:"id"`
	MultisetItemDelimiter       string         `json:"-" toml:"multiset-item-delimiter"`
	OutputDelimiter             string         `json:"-" toml:"output-delimiter"`
	Links                       []Link         `json:"links,omitempty"`
	Name                        string         `json:"name"`
	Preamble                    template.HTML  `json:"preamble"`
	Released                    toml.LocalDate `json:"released"`
	Synopsis                    string         `json:"synopsis"`
	Variants                    []*Hole        `json:"-"`
}

type HoleAnswer struct {
	Args   []string
	Answer string
}

type HoleAnswerFunc func() []HoleAnswer

func initHoles() {
	var holes map[string]*Hole
	unmarshal("data/holes.toml", &holes)

	funcs := template.FuncMap{
		"hasPrefix": strings.HasPrefix,
		"hasSuffix": strings.HasSuffix,
	}

	// Expand variants.
	copyFields := []string{
		"Category", "Data", "Links", "Preamble", "Released", "Synopsis", "Variants"}
	for name, hole := range holes {
		hole.Name = name

		if len(hole.Variants) == 0 {
			continue
		}

		// Prepend the current hole as a variant.
		hole.Variants = slices.Insert(hole.Variants, 0, hole)

		src := reflect.ValueOf(hole).Elem()
		for _, variant := range hole.Variants {
			// Already been processed.
			if _, ok := holes[variant.Name]; ok {
				continue
			}

			// Copy any fields missing in the variant from the hole.
			dst := reflect.ValueOf(variant).Elem()
			for _, field := range copyFields {
				if dst.FieldByName(field).IsZero() {
					dst.FieldByName(field).Set(src.FieldByName(field))
				}
			}

			holes[variant.Name] = variant
		}
	}

	for name, hole := range holes {
		hole.ID = ID(name)

		// Aliases.
		for _, alias := range hole.Aliases {
			HoleAliases[alias] = hole.ID
		}

		// Redirects.
		for _, redirect := range hole.Redirects {
			HoleRedirects[redirect] = hole.ID
		}

		// Answers.
		// ¯\_(ツ)_/¯ cannot embed file hole-answers/√2.txt: invalid name √2.txt
		if b, err := answers.ReadFile(
			"hole-answers/" + strings.Replace(hole.ID, "√2", "root-2", 1) + ".txt",
		); err == nil {
			hole.Answer = string(bytes.TrimSuffix(b, []byte{'\n'}))
		}

		// Unmarshall any data into an ordered map.
		if hole.Data != "" {
			if err := json.Unmarshal([]byte(hole.Data), &hole.DataMap); err != nil {
				panic(err)
			}
		}

		// Process the templated preamble & synopsis.
		preamble := template.Must(template.New("").Funcs(funcs).Parse(string(hole.Preamble)))
		synopsis := templateTxt.Must(templateTxt.New("").Funcs(funcs).Parse(hole.Synopsis))

		var b bytes.Buffer
		if err := preamble.Execute(&b, hole); err != nil {
			panic(err)
		}
		hole.Preamble = template.HTML(b.String())

		b.Reset()
		if err := synopsis.Execute(&b, hole); err != nil {
			panic(err)
		}
		hole.Synopsis = b.String()

		switch hole.Category {
		case "Art":
			hole.CategoryColor = "red"
			hole.CategoryIcon = "brush"
		case "Computing":
			hole.CategoryColor = "orange"
			hole.CategoryIcon = "cpu"
		case "Gaming":
			hole.CategoryColor = "yellow"
			hole.CategoryIcon = "joystick"
		case "Mathematics":
			hole.CategoryColor = "green"
			hole.CategoryIcon = "calculator"
		case "Sequence":
			hole.CategoryColor = "blue"
			hole.CategoryIcon = "sort-numeric-down"
		case "Transform":
			hole.CategoryColor = "purple"
			hole.CategoryIcon = "shuffle"
		}

		AllHoleByID[hole.ID] = hole
		AllHoleList = append(AllHoleList, hole)

		if hole.Experiment == 0 {
			HoleByID[hole.ID] = hole
			HoleList = append(HoleList, hole)

			HoleCategoryHstore.Map[hole.ID] =
				sql.NullString{String: hole.Category, Valid: true}
		} else {
			ExpHoleByID[hole.ID] = hole
			ExpHoleList = append(ExpHoleList, hole)
		}
	}

	// Ten most recent holes.
	RecentHoles = make([]*Hole, len(HoleList))
	copy(RecentHoles, HoleList)
	slices.SortFunc(RecentHoles, func(a, b *Hole) int {
		if c := cmp.Compare(b.Released.String(), a.Released.String()); c != 0 {
			return c
		}
		return cmp.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name))
	})
	RecentHoles = RecentHoles[:10]

	for _, holes := range [][]*Hole{HoleList, ExpHoleList, AllHoleList} {
		// Case-insensitive sort.
		slices.SortFunc(holes, func(a, b *Hole) int {
			return cmp.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name))
		})
	}
}
