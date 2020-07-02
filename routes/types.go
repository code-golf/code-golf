package routes

import (
	"html/template"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
)

type Hole struct {
	Prev, Next, ID, Name, Category, CategoryColor, CategoryIcon string
	Preamble                                                    template.HTML
}

var holes []Hole

var holeByID = map[string]Hole{}

func init() {
	var holesTOML map[string]Hole

	if _, err := toml.DecodeFile("holes.toml", &holesTOML); err != nil {
		panic(err)
	}

	for name, hole := range holesTOML {
		hole.Name = name
		hole.ID = strings.ToLower(
			strings.ReplaceAll(strings.ReplaceAll(name, "’", ""), " ", "-"))

		holes = append(holes, hole)
	}

	sort.Slice(holes, func(i, j int) bool { return holes[i].Name < holes[j].Name })

	for i, hole := range holes {
		if i == 0 {
			holes[i].Prev = holes[len(holes)-1].ID
		} else {
			holes[i].Prev = holes[i-1].ID
		}

		if i == len(holes)-1 {
			holes[i].Next = holes[0].ID
		} else {
			holes[i].Next = holes[i+1].ID
		}

		switch hole.Category {
		case "Art":
			holes[i].CategoryColor = "red"
			holes[i].CategoryIcon = "\uf53f"
		case "Computing":
			holes[i].CategoryColor = "orange"
			holes[i].CategoryIcon = "\uf544"
		case "Gaming":
			holes[i].CategoryColor = "yellow"
			holes[i].CategoryIcon = "\uf11b"
		case "Mathematics":
			holes[i].CategoryColor = "green"
			holes[i].CategoryIcon = "\uf698"
		case "Sequence":
			holes[i].CategoryColor = "blue"
			holes[i].CategoryIcon = "\uf162"
		case "Transform":
			holes[i].CategoryColor = "purple"
			holes[i].CategoryIcon = "\uf074"
		}

		holeByID[hole.ID] = holes[i]
	}
}
