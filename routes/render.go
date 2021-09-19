package routes

import (
	"crypto/rand"
	"encoding/base64"
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
	min "github.com/tdewolff/minify/v2/minify"
)

func colour(i int) string {
	if i <= 1 {
		return "yellow"
	}
	if i <= 2 {
		return "orange"
	}
	if i <= 3 {
		return "red"
	}
	if i <= 10 {
		return "purple"
	}
	if i <= 100 {
		return "blue"
	}
	return "green"
}

var (
	assets = map[string]string{}
	css    = map[string]template.CSS{}
	svg    = map[string]template.HTML{}

	dev bool
)

var tmpl = template.New("").Funcs(template.FuncMap{
	"bytes":     pretty.Bytes,
	"colour":    colour,
	"comma":     pretty.Comma,
	"hasPrefix": strings.HasPrefix,
	"hasSuffix": strings.HasSuffix,
	"html":      func(html string) template.HTML { return template.HTML(html) },
	"inc":       func(i int) int { return i + 1 },
	"ord":       pretty.Ordinal,
	"svg":       func(name string) template.HTML { return svg[name] },
	"symbol": func(name string) template.HTML {
		return template.HTML(strings.ReplaceAll(string(svg[name]), "svg", "symbol"))
	},
	"title":      strings.Title,
	"time":       pretty.Time,
	"trimPrefix": strings.TrimPrefix,
})

func getDarkModeMediaQuery(theme string) string {
	switch theme {
	case "dark":
		return "all"
	case "light":
		return "not all"
	}
	return "(prefers-color-scheme:dark)"
}

func getThemeCSS(theme string) template.CSS {
	switch theme {
	case "dark":
		return css["dark"]
	case "light":
		return css["light"]
	}

	return css["light"] + "@media(prefers-color-scheme:dark){" + css["dark"] + "}"
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

	// CSS.
	for name, data := range slurp("css") {
		var err error
		if data, err = min.CSS(data); err != nil {
			panic(err)
		}

		css[name] = template.CSS(data)
	}

	// HACK Prepend font.css (with font URL) onto base.css.
	// TODO Use esbuild for all CSS? Still serve inline?
	if fontCSS, ok := assets["css/font.css"]; ok {
		cssBytes, err := os.ReadFile(fontCSS[1:])
		if err != nil {
			panic(err)
		}
		css["base"] = template.CSS(cssBytes) + css["base"]
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
			if data, err = min.HTML(data); err != nil {
				panic(err)
			}
		}

		tmpl = template.Must(tmpl.New(name).Parse(data))
	}
}

func render(w http.ResponseWriter, r *http.Request, name string, data ...interface{}) {
	// The generated value SHOULD be at least 128 bits long (before encoding),
	// and SHOULD be generated via a cryptographically secure random number
	// generator - https://w3c.github.io/webappsec-csp/#security-nonces
	nonce := make([]byte, 16)
	if _, err := rand.Read(nonce); err != nil {
		panic(err)
	}

	type CheevoBanner struct {
		Cheevo     *config.Cheevo
		During     bool
		Start, End time.Time
	}

	theme := "auto"
	theGolfer := session.Golfer(r)
	if theGolfer != nil {
		theme = theGolfer.Theme
	}

	args := struct {
		CSS                                       template.CSS
		CheevoBanner                              *CheevoBanner
		Countries                                 map[string]*config.Country
		Data, Description, Title                  interface{}
		DarkModeMediaQuery, LogInURL, Nonce, Path string
		Golfer                                    *golfer.Golfer
		GolferInfo                                *golfer.GolferInfo
		Holes                                     map[string]*config.Hole
		JS                                        []string
		Langs                                     map[string]*config.Lang
		Location                                  *time.Location
		Request                                   *http.Request
	}{
		Countries:          config.CountryByID,
		CSS:                getThemeCSS(theme) + css["base"] + css[path.Dir(name)] + css[name],
		Data:               data[0],
		DarkModeMediaQuery: getDarkModeMediaQuery(theme),
		Description:        "Code Golf is a game designed to let you show off your code-fu by solving problems in the least number of characters.",
		Golfer:             theGolfer,
		GolferInfo:         session.GolferInfo(r),
		Holes:              config.HoleByID,
		JS:                 []string{assets["js/base.js"]},
		Langs:              config.LangByID,
		Nonce:              base64.StdEncoding.EncodeToString(nonce),
		Path:               r.URL.Path,
		Request:            r,
		Title:              "Code Golf",
	}

	if len(data) > 1 {
		args.Title = data[1]
	}

	if len(data) > 2 {
		args.Description = data[2]
	}

	if args.Golfer != nil && args.Golfer.TimeZone != nil {
		args.Location = args.Golfer.TimeZone
	} else {
		args.Location = time.UTC
	}

	// Independence Day cheevo banner. TODO Generalise.
	if args.Golfer != nil && !args.Golfer.Earnt("independence-day") {
		var (
			now   = time.Now().UTC()
			year  = now.Year()
			start = time.Date(year, time.July, 4, 0, 0, 0, 0, time.UTC)
			end   = time.Date(year, time.July, 5, 0, 0, 0, 0, time.UTC)
		)

		if now.Before(end) {
			args.CheevoBanner = &CheevoBanner{
				config.CheevoByID["independence-day"],
				start.Before(now), start, end,
			}
		}
	}

	// TODO CSS imports?
	if name == "hole" {
		args.CSS = css["vendor/codemirror"] + css["vendor/codemirror-dialog"] +
			css["vendor/codemirror-dark"] + args.CSS
	}

	// Append route specific JS.
	// e.g. GET /foo/bar might add js/foo.js and/or js/foo/bar.js.
	for _, path := range []string{path.Dir(name), name} {
		if url, ok := assets["js/"+path+".js"]; ok {
			args.JS = append(args.JS, url)
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
			"style-src 'unsafe-inline'",
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

	switch name {
	case "403":
		w.WriteHeader(http.StatusForbidden)
	case "404":
		w.WriteHeader(http.StatusNotFound)
	}

	if err := tmpl.ExecuteTemplate(w, name, args); err != nil {
		panic(err)
	}
}
