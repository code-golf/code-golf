package config

import (
	"html/template"
	"sort"
	"strings"

	"github.com/tdewolff/minify/v2/minify"
)

var (
	HoleByID = map[string]*Hole{}
	HoleList []*Hole

	ExpHoleByID = map[string]*Hole{}
	ExpHoleList []*Hole
)

type Hole struct {
	Category, CategoryColor, CategoryIcon string
	Experiment                            int
	ID, Name, Prev, Next                  string
	Preamble                              template.HTML
	Links                                 []struct{ Name, URL string }
}

func init() {
	var holes map[string]*Hole
	unmarshal("holes.toml", &holes)

	for name, hole := range holes {
		hole.ID = id(name)
		hole.Name = name

		// Minify preamble.
		if html, err := minify.HTML(string(hole.Preamble)); err != nil {
			panic(err)
		} else {
			hole.Preamble = template.HTML(html)
		}

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

		if hole.Experiment == 0 {
			HoleByID[hole.ID] = hole
			HoleList = append(HoleList, hole)
		} else {
			ExpHoleByID[hole.ID] = hole
			ExpHoleList = append(ExpHoleList, hole)
		}
	}

	// Case-insensitive sort.
	sort.Slice(HoleList, func(i, j int) bool {
		return strings.ToLower(HoleList[i].Name) <
			strings.ToLower(HoleList[j].Name)
	})
	sort.Slice(ExpHoleList, func(i, j int) bool {
		return strings.ToLower(ExpHoleList[i].Name) <
			strings.ToLower(ExpHoleList[j].Name)
	})

	// Set Prev, Next.
	for i, hole := range HoleList {
		if i == 0 {
			hole.Prev = HoleList[len(HoleList)-1].ID
		} else {
			hole.Prev = HoleList[i-1].ID
		}

		if i == len(HoleList)-1 {
			hole.Next = HoleList[0].ID
		} else {
			hole.Next = HoleList[i+1].ID
		}
	}
	for i, hole := range ExpHoleList {
		if i == 0 {
			hole.Prev = ExpHoleList[len(ExpHoleList)-1].ID
		} else {
			hole.Prev = ExpHoleList[i-1].ID
		}

		if i == len(ExpHoleList)-1 {
			hole.Next = ExpHoleList[0].ID
		} else {
			hole.Next = ExpHoleList[i+1].ID
		}
	}
}
