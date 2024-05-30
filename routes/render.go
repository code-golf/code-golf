package routes

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/golfer"
	"github.com/code-golf/code-golf/session"
	"github.com/code-golf/code-golf/views"
)

var assets = map[string]string{}

var dev bool

func init() {
	_, dev = os.LookupEnv("DEV")

	// HACK Tests are run from the package directory, walk a dir up.
	if _, err := os.Stat("esbuild.json"); os.IsNotExist(err) {
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

	if err := views.Render(w, name, args); err != nil {
		panic(err)
	}
}
