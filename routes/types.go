package routes

import (
	"html/template"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
)

type Lang struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Hole struct {
	Prev, Next, ID, Name, Category, CategoryColor, CategoryIcon string
	Preamble                                                    template.HTML
}

type Trophy struct {
	ID, Name    string
	Description template.HTML
}

var langs = []Lang{
	{"bash", "Bash"},
	{"brainfuck", "Brainfuck"},
	{"c", "C"},
	{"fortran", "Fortran"},
	{"go", "Go"},
	{"haskell", "Haskell"},
	{"j", "J"},
	{"javascript", "JavaScript"},
	{"julia", "Julia"},
	{"lisp", "Lisp"},
	{"lua", "Lua"},
	{"nim", "Nim"},
	{"perl", "Perl"},
	{"php", "PHP"},
	{"python", "Python"},
	{"raku", "Raku"},
	{"ruby", "Ruby"},
	{"rust", "Rust"},
	{"swift", "Swift"},
}

var trophies = []Trophy{
	{
		"elephpant-in-the-room",
		"ElePHPant in the Room",
		"Solve any hole in PHP.",
	},
	{
		"happy-birthday-code-golf",
		"Happy Birthday, Code Golf",
		"Solve any hole in any language on <a href=//github.com/code-golf/code-golf/commit/4b44>2 Oct</a>.",
	},
	{
		"hello-world",
		"Hello, World!",
		"Solve any hole in any language.",
	},
	{
		"inception",
		"Inception",
		"Solve <a href=/brainfuck#brainfuck>Brainfuck in Brainfuck</a>.",
	},
	{
		"interview-ready",
		"Interview Ready",
		"Solve <a href=/fizz-buzz>Fizz Buzz</a> in any language.",
	},
	{
		"its-over-9000",
		"It’s Over 9000!",
		"Earn over 9,000 points.",
	},
	{
		"my-god-its-full-of-stars",
		"My God, It’s Full of Stars",
		"Star <a href=//github.com/code-golf/code-golf>the Code Golf repository</a>.",
	},
	{
		"ouroboros",
		"Ouroboros",
		"Solve <a href=/quine#python>Quine in Python</a>.",
	},
	{
		"patches-welcome",
		"Patches Welcome",
		"Contribute a merged PR to <a href=//github.com/code-golf/code-golf>the Code Golf repository</a>.",
	},
	{
		"polyglot",
		"Polyglot",
		"Solve at least one hole in every language.",
	},
	{
		"slowcoach",
		"Slowcoach",
		"Fail an attempt by exceeding the time limit.",
	},
	{
		"tim-toady",
		"Tim Toady",
		"Solve the same hole in both Perl and Raku.",
	},
	{
		"the-watering-hole",
		"The Watering Hole",
		"Solve your nineteenth hole.",
	},
}

var holes []Hole
var langByID = map[string]Lang{}
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

	for _, lang := range langs {
		langByID[lang.ID] = lang
	}

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
