package config

import (
	"bytes"
	"cmp"
	"encoding/json"
	"html/template"
	"slices"
	"strings"
	templateTxt "text/template"

	"github.com/code-golf/code-golf/ordered"
	"github.com/pelletier/go-toml/v2"
	"github.com/tdewolff/minify/v2/minify"
)

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
)

type Link struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	V    []int  `json:"-"`
}

type Hole struct {
	Categories                              []string         `json:"-"`
	Category                                string           `json:"category"`
	CategoryColor, CategoryIcon, Prev, Next string           `json:"-"`
	Data                                    template.JS      `json:"-"`
	Experiment                              int              `json:"-"`
	ID                                      string           `json:"id"`
	Links                                   []Link           `json:"links"`
	Name                                    string           `json:"name"`
	Preamble                                template.HTML    `json:"preamble"`
	Released                                toml.LocalDate   `json:"released"`
	Releases                                []toml.LocalDate `json:"-"`
	Synopsis                                string           `json:"synopsis"`
	Variants                                []*Hole          `json:"-"`
}

func init() {
	var holes map[string]*struct {
		Hole
		Variants []string
	}
	unmarshal("holes.toml", &holes)

	funcs := template.FuncMap{
		"hasPrefix": strings.HasPrefix,
		"hasSuffix": strings.HasSuffix,
	}

	// Expand variants.
	for name, hole := range holes {
		// Don't process holes without variants or already processed variants.
		if len(hole.Variants) == 0 || len(hole.Hole.Variants) != 0 {
			continue
		}

		// Delete the original meta hole.
		delete(holes, name)

		// Parse the templated preamble.
		preamble, err := template.New("").Funcs(funcs).Parse(string(hole.Preamble))
		if err != nil {
			panic(err)
		}

		// Parse the templated synopsis.
		synopsis, err := templateTxt.New("").Funcs(funcs).Parse(hole.Synopsis)
		if err != nil {
			panic(err)
		}

		var variants []*Hole
		for i, variant := range hole.Variants {
			hole := *hole

			// Process the templated preamble with the current variant.
			var b bytes.Buffer
			if err := preamble.Execute(&b, variant); err != nil {
				panic(err)
			}
			hole.Preamble = template.HTML(b.String())

			// Process the templated synopsis with the current variant.
			b.Reset()
			if err := synopsis.Execute(&b, variant); err != nil {
				panic(err)
			}
			hole.Synopsis = b.String()

			holes[variant] = &hole
			variants = append(variants, &hole.Hole)

			if len(hole.Categories) > 0 {
				hole.Category = hole.Categories[i]
			}

			if len(hole.Releases) > 0 {
				hole.Released = hole.Releases[i]
			}
		}

		// Reference the variants from each variant.
		for _, variant := range variants {
			variant.Variants = variants
		}
	}

	for name, hole := range holes {
		hole.ID = ID(name)
		hole.Name = name

		// Process the templated preamble with the data.
		if hole.Data != "" {
			t, err := template.New("").Parse(string(hole.Preamble))
			if err != nil {
				panic(err)
			}

			var data ordered.Map
			if err := json.Unmarshal([]byte(hole.Data), &data); err != nil {
				panic(err)
			}

			var b bytes.Buffer
			if err := t.Execute(&b, data); err != nil {
				panic(err)
			}
			hole.Preamble = template.HTML(b.String())
		}

		// Minify preamble.
		if html, err := minify.HTML(string(hole.Preamble)); err != nil {
			panic(err)
		} else {
			hole.Preamble = template.HTML(html)
		}

		// Filter out links that don't match this variant.
		links := make([]Link, 0, len(hole.Links))
		for _, link := range hole.Links {
			if len(link.V) == 0 || slices.ContainsFunc(
				link.V, func(i int) bool { return hole.Variants[i] == name },
			) {
				links = append(links, link)
			}
		}
		hole.Links = links

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

		AllHoleByID[hole.ID] = &hole.Hole
		AllHoleList = append(AllHoleList, &hole.Hole)

		if hole.Experiment == 0 {
			HoleByID[hole.ID] = &hole.Hole
			HoleList = append(HoleList, &hole.Hole)
		} else {
			ExpHoleByID[hole.ID] = &hole.Hole
			ExpHoleList = append(ExpHoleList, &hole.Hole)
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

	for i, holes := range [][]*Hole{HoleList, ExpHoleList, AllHoleList} {
		// Case-insensitive sort.
		slices.SortFunc(holes, func(a, b *Hole) int {
			return cmp.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name))
		})

		// Set Prev, Next. Not for "AllHoleList" as it would overwrite.
		if i < 2 {
			for j, hole := range holes {
				if j == 0 {
					hole.Prev = holes[len(holes)-1].ID
				} else {
					hole.Prev = holes[j-1].ID
				}

				if j == len(holes)-1 {
					hole.Next = holes[0].ID
				} else {
					hole.Next = holes[j+1].ID
				}
			}
		}
	}
}
