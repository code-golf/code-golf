package lang

type Lang struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var ByID = map[string]Lang{}

func init() {
	for _, lang := range List {
		ByID[lang.ID] = lang
	}
}

var List = []Lang{
	{"bash", "Bash"},
	{"brainfuck", "Brainfuck"},
	{"c", "C"},
	{"c-sharp", "C#"},
	{"f-sharp", "F#"},
	{"fortran", "Fortran"},
	{"go", "Go"},
	{"haskell", "Haskell"},
	{"j", "J"},
	{"java", "Java"},
	{"javascript", "JavaScript"},
	{"julia", "Julia"},
	{"lisp", "Lisp"},
	{"lua", "Lua"},
	{"nim", "Nim"},
	{"perl", "Perl"},
	{"php", "PHP"},
	{"powershell", "PowerShell"},
	{"python", "Python"},
	{"raku", "Raku"},
	{"ruby", "Ruby"},
	{"rust", "Rust"},
	{"swift", "Swift"},
}
