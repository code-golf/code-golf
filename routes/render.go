package routes

import (
	"crypto/rand"
	"encoding/base64"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/code-golf/code-golf/cheevo"
	"github.com/code-golf/code-golf/country"
	"github.com/code-golf/code-golf/golfer"
	"github.com/code-golf/code-golf/hole"
	"github.com/code-golf/code-golf/lang"
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
	css = map[string]template.CSS{}
	js  = map[string]template.JS{}
	svg = map[string]template.HTML{}

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
	"timeShort":  pretty.TimeShort,
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

func init() {
	_, dev = syscall.Getenv("DEV")
	uppercaseProps := regexp.MustCompile(`{{.+?[A-Z].*?}}`)

	if err := filepath.Walk("views", func(file string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		b, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		data := string(b)
		ext := path.Ext(file)
		name := file[len("views/") : len(file)-len(ext)]

		switch ext {
		case ".css":
			data = strings.ReplaceAll(data, "twemojiWoff2", twemojiWoff2Path)

			if data, err = min.CSS(data); err != nil {
				return err
			}

			css[name[len("css/"):]] = template.CSS(data)
		case ".html":
			// Minify templates without uppsercase properties.
			// The real fix is https://github.com/tdewolff/minify/issues/35
			if !uppercaseProps.MatchString(data) {
				if data, err = min.HTML(data); err != nil {
					return err
				}
			}

			tmpl = template.Must(tmpl.New(name).Parse(data))
		case ".js":
			if data, err = min.JS(data); err != nil {
				return err
			}

			js[name[len("js/"):]] = template.JS(data)
		case ".svg":
			// Trim namespace as it's not needed for inline SVG under HTML 5.
			data = strings.ReplaceAll(data, ` xmlns="http://www.w3.org/2000/svg"`, "")

			svg[name[len("svg/"):]] = template.HTML(data)
		}

		return nil
	}); err != nil {
		panic(err)
	}
}

func render(w http.ResponseWriter, r *http.Request, name string, data interface{}, meta ...string) {
	// The generated value SHOULD be at least 128 bits long (before encoding),
	// and SHOULD be generated via a cryptographically secure random number
	// generator - https://w3c.github.io/webappsec-csp/#security-nonces
	nonce := make([]byte, 16)
	if _, err := rand.Read(nonce); err != nil {
		panic(err)
	}

	type CheevoBanner struct {
		Cheevo     *cheevo.Cheevo
		During     bool
		Start, End time.Time
	}

	theme := "auto"
	theGolfer := session.Golfer(r)
	if theGolfer != nil {
		theme = theGolfer.Theme
	}

	args := struct {
		CSS                                                                  template.CSS
		CheevoBanner                                                         *CheevoBanner
		Countries                                                            map[string]*country.Country
		Data                                                                 interface{}
		DarkModeMediaQuery, Description, JSExt, LogInURL, Nonce, Path, Title string
		Golfer                                                               *golfer.Golfer
		GolferInfo                                                           *golfer.GolferInfo
		Holes                                                                map[string]hole.Hole
		JS                                                                   template.JS
		Langs                                                                map[string]lang.Lang
		Location                                                             *time.Location
		Request                                                              *http.Request
	}{
		Countries:          country.ByID,
		CSS:                getThemeCSS(theme) + css["base"] + css[path.Dir(name)] + css[name],
		Data:               data,
		DarkModeMediaQuery: getDarkModeMediaQuery(theme),
		Description:        "Code Golf is a game designed to let you show off your code-fu by solving problems in the least number of characters.",
		Golfer:             theGolfer,
		GolferInfo:         session.GolferInfo(r),
		Holes:              hole.ByID,
		Langs:              lang.ByID,
		JS:                 js["base"] + js[path.Dir(name)] + js[name],
		Nonce:              base64.StdEncoding.EncodeToString(nonce),
		Path:               r.URL.Path,
		Request:            r,
		Title:              "Code Golf",
	}

	if len(meta) > 0 {
		args.Title = meta[0]
	}

	if len(meta) > 1 {
		args.Description = meta[1]
	}

	if args.Golfer != nil && args.Golfer.TimeZone != nil {
		args.Location = args.Golfer.TimeZone
	} else {
		args.Location = time.UTC
	}

	// Star Wars cheevo banner. TODO Generalise.
	if args.Golfer != nil && !args.Golfer.Earnt("may-the-4ᵗʰ-be-with-you") {
		var (
			now   = time.Now().UTC()
			year  = now.Year()
			start = time.Date(year, time.May, 4, 0, 0, 0, 0, time.UTC)
			end   = time.Date(year, time.May, 5, 0, 0, 0, 0, time.UTC)
		)

		if now.Before(end) {
			args.CheevoBanner = &CheevoBanner{
				cheevo.ByID["may-the-4ᵗʰ-be-with-you"],
				start.Before(now), start, end,
			}
		}
	}

	if name == "hole" {
		args.JSExt = holeJsPath
		args.CSS = css["vendor/codemirror"] + css["vendor/codemirror-dialog"] + css["vendor/codemirror-dark"] + args.CSS
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
		config := config

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
