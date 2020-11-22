package hole

import (
	"html/template"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
	min "github.com/tdewolff/minify/v2/minify"
)

type Hole struct {
	Experiment                                                  int
	Prev, Next, ID, Name, Category, CategoryColor, CategoryIcon string
	Preamble                                                    template.HTML
	Links                                                       []struct{ Name, URL string }
}

var (
	ByID = map[string]Hole{}
	List []Hole

	ExperimentalByID = map[string]Hole{}
	ExperimentalList []Hole
)

func init() {
	var holesTOML map[string]Hole

	if _, err := toml.DecodeFile("holes.toml", &holesTOML); err != nil {
		panic(err)
	}

	for name, hole := range holesTOML {
		hole.Name = name
		hole.ID = strings.ToLower(
			strings.ReplaceAll(strings.ReplaceAll(name, "’", ""), " ", "-"))

		// Minify HTML
		if html, err := min.HTML(string(hole.Preamble)); err != nil {
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
			List = append(List, hole)
		} else {
			ExperimentalList = append(ExperimentalList, hole)
		}
	}

	// Case-insensitive sort.
	sort.Slice(List, func(i, j int) bool {
		return strings.ToLower(List[i].Name) < strings.ToLower(List[j].Name)
	})

	for _, hole := range ExperimentalList {
		ExperimentalByID[hole.ID] = hole
	}

	for i, hole := range List {
		if i == 0 {
			List[i].Prev = List[len(List)-1].ID
		} else {
			List[i].Prev = List[i-1].ID
		}

		if i == len(List)-1 {
			List[i].Next = List[0].ID
		} else {
			List[i].Next = List[i+1].ID
		}

		ByID[hole.ID] = List[i]
	}
}
