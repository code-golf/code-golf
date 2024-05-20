package routes

import (
	"encoding/json"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/golfer"
	"github.com/code-golf/code-golf/pager"
	"github.com/code-golf/code-golf/pretty"
	"github.com/code-golf/code-golf/session"
)

var (
	assets = map[string]string{}
	svg    = map[string]template.HTML{}

	dev bool
)

var tmpl = template.New("").Funcs(template.FuncMap{
	"bytes":     pretty.Bytes,
	"comma":     pretty.Comma,
	"dec":       func(i int) int { return i - 1 },
	"hasPrefix": strings.HasPrefix,
	"hasSuffix": strings.HasSuffix,
	"html":      func(html string) template.HTML { return template.HTML(html) },
	"inc":       func(i int) int { return i + 1 },
	"ord":       pretty.Ordinal,
	"page":      func(i int) int { return i/pager.PerPage + 1 },
	"param":     param,
	"svg":       func(name string) template.HTML { return svg[name] },
	"symbol": func(name string) template.HTML {
		return template.HTML(strings.ReplaceAll(string(svg[name]), "svg", "symbol"))
	},
	"title":      pretty.Title,
	"time":       pretty.Time,
	"trimPrefix": strings.TrimPrefix,
})

func slurp(dir string) map[string]string {
	files := map[string]string{}

	if err := filepath.Walk(dir, func(file string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			data, err := os.ReadFile(file)
			if err != nil {
				return err
			}

			name := file[len(dir)+1 : len(file)-len(path.Ext(file))]
			files[name] = string(data)
		}

		return err
	}); err != nil {
		panic(err)
	}

	return files
}

func init() {
	_, dev = os.LookupEnv("DEV")

	// HACK Tests are run from the package directory, walk a dir up.
	if _, err := os.Stat("css"); os.IsNotExist(err) {
		os.Chdir("..")
	}

	// Assets.
	if file, err := os.Open("esbuild.json"); err == nil {
		defer file.Close()

		var esbuild struct {
			Outputs map[string]struct{ EntryPoint string }
		}

		if err := json.NewDecoder(file).Decode(&esbuild); err != nil {
			panic(err)
		}

		for dist, src := range esbuild.Outputs {
			if src.EntryPoint != "" {
				assets[src.EntryPoint] = "/" + dist
			}
		}
	}

	// SVG.
	for name, data := range slurp("svg") {
		// Trim namespace as it's not needed for inline SVG under HTML 5.
		data = strings.ReplaceAll(data, ` xmlns="http://www.w3.org/2000/svg"`, "")

		svg[name] = template.HTML(data)
	}

	// Views.
	for name, data := range slurp("views") {
		tmpl = template.Must(tmpl.New(name).Parse(data))
	}
}

func render(w http.ResponseWriter, r *http.Request, name string, data ...any) {
	theme := "auto"
	theGolfer := session.Golfer(r)
	if theGolfer != nil {
		theme = theGolfer.Theme
	}

	type cssLink struct{ Path, Media string }

	args := struct {
		Banners                            []banner
		CSS                                []cssLink
		AutoDarkCSSLink                    string
		Cheevos                            map[string][]*config.Cheevo
		Data, Description, Title           any
		LogInURL, Name, Nonce, Path, Theme string
		Golfer                             *golfer.Golfer
		GolferInfo                         *golfer.GolferInfo
		Holes                              map[string]*config.Hole
		JS                                 []string
		Langs                              map[string]*config.Lang
		Location                           *time.Location
		Nav                                *config.Navigaton
		Request                            *http.Request
		Settings                           []*config.Setting
	}{
		Banners:     banners(theGolfer, time.Now().UTC()),
		Cheevos:     config.CheevoTree,
		CSS:         []cssLink{{assets["css/common/base.css"], ""}},
		Data:        data[0],
		Description: "Code Golf is a game designed to let you show off your code-fu by solving problems in the least number of characters.",
		Golfer:      theGolfer,
		GolferInfo:  session.GolferInfo(r),
		Holes:       config.HoleByID,
		JS:          []string{assets["js/base.tsx"]},
		Langs:       config.LangByID,
		Name:        name,
		Nonce:       nonce(),
		Path:        r.URL.Path,
		Request:     r,
		Settings:    config.Settings[strings.TrimSuffix(name, "-tabs")],
		Theme:       theme,
		Title:       "Code Golf",
	}

	if len(data) > 1 {
		args.Title = data[1]
	}

	if len(data) > 2 {
		args.Description = data[2]
	}

	if args.Golfer != nil {
		args.Location = args.Golfer.Location()
	} else {
		args.Location = time.UTC
	}

	// Get route specific CSS, JS, and navigation by splitting the name.
	// e.g. foo/bar/baz → foo, foo/bar, foo/bar/baz.
	subName := ""
	for _, part := range strings.Split(name, "/") {
		subName = path.Join(subName, part)

		if nav, ok := config.Nav[subName]; ok {
			args.Nav = nav
		}

		if url, ok := assets["css/"+subName+".css"]; ok {
			args.CSS = append(args.CSS, cssLink{Path: url})
		}

		for _, ext := range []string{"ts", "tsx"} {
			if url, ok := assets["js/"+subName+"."+ext]; ok {
				args.JS = append(args.JS, url)
			}
		}
	}

	// If theme is "dark" or "light", only that CSS file is loaded as you expect
	// If theme is auto, then the light theme is loaded, and the dark theme might
	// be loaded based on a <link media="..."> query
	if theme == "auto" {
		args.CSS = append(
			args.CSS,
			cssLink{assets["css/common/light.css"], ""},
			cssLink{assets["css/common/dark.css"], "(prefers-color-scheme:dark)"},
		)
	} else {
		args.CSS = append(args.CSS, cssLink{assets["css/common/"+theme+".css"], ""})
	}

	header := w.Header()

	header.Set("Content-Language", "en")
	header.Set("Content-Type", "text/html; charset=utf-8")
	header.Set("Referrer-Policy", "no-referrer")
	header.Set("X-Content-Type-Options", "nosniff")
	header.Set("X-Frame-Options", "DENY")
	header.Set("Content-Security-Policy",
		"base-uri 'none';"+
			"connect-src 'self';"+
			"default-src 'none';"+
			"form-action 'self';"+
			"font-src 'self';"+
			"frame-ancestors 'none';"+
			"img-src 'self' data: avatars.githubusercontent.com;"+
			"script-src 'self' 'nonce-"+args.Nonce+"';"+
			"style-src 'self' 'unsafe-inline'",
	)

	if args.Golfer == nil {
		// Shallow copy because we want to modify a string.
		config := oauthConfig

		config.RedirectURL = "https://code.golf/callback"

		if dev {
			config.RedirectURL += "/dev"
		}

		config.RedirectURL += "?redirect_uri=" + url.QueryEscape(r.RequestURI)

		// TODO State is a token to protect the user from CSRF attacks.
		args.LogInURL = config.AuthCodeURL("")
	}

	if err := tmpl.ExecuteTemplate(w, name, args); err != nil {
		panic(err)
	}
}
