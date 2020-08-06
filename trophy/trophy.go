package trophy

import "html/template"

type Trophy struct {
	Emoji, ID, Name string
	Description     template.HTML
}

var ByID = map[string]Trophy{}

func init() {
	for _, trophy := range List {
		ByID[trophy.ID] = trophy
	}
}

var List = []Trophy{
	{
		"☕", "caffeinated", "Caffeinated",
		"Solve any hole in both Java and JavaScript.",
	},
	{
		"🐘", "elephpant-in-the-room", "ElePHPant in the Room",
		"Solve any hole in PHP.",
	},
	{
		"🎂", "happy-birthday-code-golf", "Happy Birthday, Code Golf",
		"Solve any hole on <a href=//github.com/code-golf/code-golf/commit/4b44>2 Oct</a>.",
	},
	{
		"👋", "hello-world", "Hello, World!",
		"Solve your first hole.",
	},
	{
		"🧠", "inception", "Inception",
		"Solve <a href=/brainfuck#brainfuck>Brainfuck in Brainfuck</a>.",
	},
	{
		"🇺🇸", "independence-day", "Independence Day",
		"Solve <a href=/united-states>United States</a> on <a href=//www.wikipedia.org/wiki/Independence_Day_(United_States)>4 Jul</a>.",
	},
	{
		"💼", "interview-ready", "Interview Ready",
		"Solve <a href=/fizz-buzz>Fizz Buzz</a>.",
	},
	{
		"🐉", "its-over-9000", "It’s Over 9000!",
		"Earn over 9,000 points.",
	},
	{
		"⭐", "my-god-its-full-of-stars", "My God, It’s Full of Stars",
		"Star <a href=//github.com/code-golf/code-golf>the Code Golf repository</a>.",
	},
	{
		"🐍", "ouroboros", "Ouroboros",
		"Solve <a href=/quine#python>Quine in Python</a>.",
	},
	{
		"💾", "patches-welcome", "Patches Welcome",
		"Contribute a merged PR to <a href=//github.com/code-golf/code-golf>the Code Golf repository</a>.",
	},
	{
		"🥧", "pi-day", "Pi Day",
		"Solve <a href=/π>π</a> on <a href=//www.wikipedia.org/wiki/Pi_Day>14 Mar</a>.",
	},
	{
		"🔣", "polyglot", "Polyglot",
		"Solve any hole in every language.",
	},
	{
		"🦥", "slowcoach", "Slowcoach",
		"Fail any hole by exceeding the time limit.",
	},
	{
		"🐪", "tim-toady", "Tim Toady",
		"Solve any hole in both Perl and Raku.",
	},
	{
		"🍺", "the-watering-hole", "The Watering Hole",
		"Solve your nineteenth hole.",
	},
}
