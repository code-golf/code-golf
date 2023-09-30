package routes

import (
	"encoding/json"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/golfer"
	"github.com/code-golf/code-golf/pretty"
	"github.com/code-golf/code-golf/session"
	"github.com/tdewolff/minify/v2/minify"
)

var (
	assets = map[string]string{}
	svg    = map[string]template.HTML{}

	dev bool
)

var tmpl = template.New("").Funcs(template.FuncMap{
	"bytes":      pretty.Bytes,
	"colour":     colour,
	"colourRank": colourRank,
	"comma":      pretty.Comma,
	"dec":        func(i int) int { return i - 1 },
	"hasPrefix":  strings.HasPrefix,
	"hasSuffix":  strings.HasSuffix,
	"html":       func(html string) template.HTML { return template.HTML(html) },
	"inc":        func(i int) int { return i + 1 },
	"ord":        pretty.Ordinal,
	"param":      param,
	"svg":        func(name string) template.HTML { return svg[name] },
	"symbol": func(name string) template.HTML {
		return template.HTML(strings.ReplaceAll(string(svg[name]), "svg", "symbol"))
	},
	"title":      pretty.Title,
	"time":       pretty.Time,
	"trimPrefix": strings.TrimPrefix,
})

type cssLink struct{ Path, Media string }

func getCSSLinks(name string, theme string) []cssLink {
	var links = []cssLink{{assets["css/common/base.css"], ""}}

	if path, ok := assets["css/"+name+".css"]; ok {
		links = append(links, cssLink{Path: path})
	}

	// If theme is "dark" or "light", only that CSS file is loaded as you expect
	// If theme is auto, then the light theme is loaded, and the dark theme might
	// be loaded based on a <link media="..."> query
	if theme == "dark" || theme == "light" {
		return append(links, cssLink{assets["css/common/"+theme+".css"], ""})
	}

	return append(
		links,
		cssLink{assets["css/common/light.css"], ""},
		cssLink{assets["css/common/dark.css"], "(prefers-color-scheme:dark)"},
	)
}

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
	uppercaseProps := regexp.MustCompile(`{{.+?[A-Z].*?}}`)
	for name, data := range slurp("views") {
		// Minify templates without uppercase properties.
		// The real fix is https://github.com/tdewolff/minify/issues/35
		if !uppercaseProps.MatchString(data) {
			var err error
			if data, err = minify.HTML(data); err != nil {
				panic(err)
			}
		}

		tmpl = template.Must(tmpl.New(name).Parse(data))
	}
}

func render(w http.ResponseWriter, r *http.Request, name string, data ...any) {
	theme := "auto"
	theGolfer := session.Golfer(r)
	if theGolfer != nil {
		theme = theGolfer.Theme
	}

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
	}{
		Banners:     banners(theGolfer, time.Now().UTC()),
		Cheevos:     config.CheevoTree,
		CSS:         getCSSLinks(name, theme),
		Data:        data[0],
		Description: "Code Golf is a game designed to let you show off your code-fu by solving problems in the least number of characters.",
		Golfer:      theGolfer,
		GolferInfo:  session.GolferInfo(r),
		Holes:       config.HoleByID,
		JS:          []string{assets["js/base.tsx"]},
		Langs:       config.LangByID,
		Name:        name,
		Nav:         config.Nav[name],
		Nonce:       nonce(),
		Path:        r.URL.Path,
		Request:     r,
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

	// Append route specific JS.
	// e.g. GET /foo/bar might add js/foo.ts and/or js/foo/bar.tsx.
	for _, path := range []string{path.Dir(name), name} {
		for _, ext := range []string{"ts", "tsx"} {
			if url, ok := assets["js/"+path+"."+ext]; ok {
				args.JS = append(args.JS, url)
			}
		}
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
