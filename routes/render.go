package routes

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/JRaspass/code-golf/cookie"
	"github.com/JRaspass/code-golf/pretty"
)

const (
	day   = 24 * time.Hour
	week  = 7 * day
	month = 4 * week
)

func ord(i int) string {
	switch i % 10 {
	case 1:
		if i%100 != 11 {
			return "st"
		}
	case 2:
		if i%100 != 12 {
			return "nd"
		}
	case 3:
		if i%100 != 13 {
			return "rd"
		}
	}
	return "th"
}

var tmpl = template.New("").Funcs(template.FuncMap{
	"comma":     pretty.Comma,
	"hasPrefix": strings.HasPrefix,
	"hasSuffix": strings.HasSuffix,
	"ord":       ord,
	"symbol": func(name string) template.HTML {
		if svg, err := ioutil.ReadFile("views/" + name + ".svg"); err != nil {
			panic(err)
		} else {
			return template.HTML(bytes.ReplaceAll(svg, []byte("svg"), []byte("symbol")))
		}
	},
	"title": strings.Title,
	"time": func(t time.Time) template.HTML {
		var sb strings.Builder

		rfc := t.Format(time.RFC3339)

		sb.WriteString("<time datetime=")
		sb.WriteString(rfc)
		sb.WriteString(" title=")
		sb.WriteString(rfc)
		sb.WriteRune('>')

		switch diff := time.Now().Sub(t); true {
		case diff < 2*time.Minute:
			sb.WriteString("1 min ago")
		case diff < 2*time.Hour:
			fmt.Fprintf(&sb, "%d mins ago", diff/time.Minute)
		case diff < 2*day:
			fmt.Fprintf(&sb, "%d hours ago", diff/time.Hour)
		case diff < 2*week:
			fmt.Fprintf(&sb, "%d days ago", diff/day)
		case diff < 2*month:
			fmt.Fprintf(&sb, "%d weeks ago", diff/week)
		default:
			fmt.Fprintf(&sb, "%d months ago", diff/month)
		}

		sb.WriteString("</time>")

		return template.HTML(sb.String())
	},
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
		args.LogInURL = "//github.com/login/oauth/authorize?client_id=7f6709819023e9215205&scope=user:email&redirect_uri=https://code-golf.io/callback?redirect_uri%3D" + url.QueryEscape(url.QueryEscape(r.RequestURI))
	}

	w.WriteHeader(code)

	if err := tmpl.ExecuteTemplate(w, name, args); err != nil {
		panic(err)
	}
}
