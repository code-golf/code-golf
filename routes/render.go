package routes

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/code-golf/code-golf/cookie"
	"github.com/code-golf/code-golf/pretty"
)

var tmpl = template.New("").Funcs(template.FuncMap{
	"comma":     pretty.Comma,
	"hasPrefix": strings.HasPrefix,
	"hasSuffix": strings.HasSuffix,
	"ord":       pretty.Ordinal,
	"symbol": func(name string) template.HTML {
		if svg, err := ioutil.ReadFile("views/" + name + ".svg"); err != nil {
			panic(err)
		} else {
			return template.HTML(bytes.ReplaceAll(svg, []byte("svg"), []byte("symbol")))
		}
	},
	"title": strings.Title,
	"time":  pretty.Time,
})

func init() {
	// Tests run from the package directory, walk upwards until we find views.
	for {
		if _, err := os.Stat("views"); os.IsNotExist(err) {
			os.Chdir("..")
		} else {
			break
		}
	}

	if err := filepath.Walk("views", func(file string, _ os.FileInfo, err error) error {
		if ext := path.Ext(file); ext == ".html" || ext == ".svg" {
			if b, err := ioutil.ReadFile(file); err != nil {
				return err
			} else {
				name := file[len("views/") : len(file)-len(ext)]
				tmpl = template.Must(tmpl.New(name).Parse(string(b)))
			}
		}

		return err
	}); err != nil {
		panic(err)
	}
}

// Render wraps common logic required for rendering a view to the user.
func Render(
	w http.ResponseWriter,
	r *http.Request,
	code int,
	name, title string,
	data interface{},
) {
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
			"form-action 'none';"+
			"font-src 'self';"+
			"frame-ancestors 'none';"+
			"img-src 'self' data: avatars.githubusercontent.com;"+
			"script-src 'self';"+
			"style-src 'self'",
	)

	args := struct {
		CommonCssPath, Login, LogInURL, Path, Title string
		Data                                        interface{}
	}{
		CommonCssPath: commonCssPath,
		Data:          data,
		Path:          r.URL.Path,
		Title:         title,
	}

	if _, args.Login = cookie.Read(r); args.Login == "" {
		// Shallow copy because we want to modify a string.
		config := config

		config.RedirectURL = "https://code-golf.io/callback?redirect_uri=" +
			url.QueryEscape(r.RequestURI)

		// TODO State is a token to protect the user from CSRF attacks.
		args.LogInURL = config.AuthCodeURL("")
	}

	w.WriteHeader(code)

	if err := tmpl.ExecuteTemplate(w, name, args); err != nil {
		panic(err)
	}
}
