package views

import (
	"embed"
	"encoding/xml"
	"html/template"
	"io/fs"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/golfer"
	"github.com/code-golf/code-golf/pager"
	"github.com/code-golf/code-golf/pretty"
)

//go:embed html/*
var views embed.FS

var tmpl = template.New("").Funcs(template.FuncMap{
	"atoi":      func(s string) int { i, _ := strconv.Atoi(s); return i },
	"bytes":     pretty.Bytes,
	"comma":     pretty.Comma,
	"dec":       func(i int) int { return i - 1 },
	"duration":  pretty.Duration,
	"hasPrefix": strings.HasPrefix,
	"hasSuffix": strings.HasSuffix,
	"html":      func(html string) template.HTML { return template.HTML(html) },
	"id":        config.ID,
	"inc":       func(i int) int { return i + 1 },
	"ord":       pretty.Ordinal,
	"page":      pager.Page,
	"title":     pretty.Title,
	"time":      pretty.Time,

	"amount": func(i int, term string) string {
		if i != 1 {
			term += "s"
		}
		return pretty.Comma(i) + " " + term
	},

	// Backend version, keep in sync with js/util.js.
	"avatar": func(rawURL string, size int) (*url.URL, error) {
		u, err := url.Parse(rawURL)
		if err != nil {
			return nil, err
		}

		// Set the avatar size. Most support "s" and "size", but i.sstatic.net
		// only supports "s" and cdn.discordapp.com only supports "size".
		q := u.Query()
		if u.Host == "cdn.discordapp.com" {
			q.Set("size", strconv.Itoa(size))
		} else {
			q.Set("s", strconv.Itoa(size))
		}
		u.RawQuery = q.Encode()

		return u, nil
	},

	"map": func(kv ...any) map[any]any {
		m := make(map[any]any, len(kv)/2)
		for i := 0; i < len(kv); i += 2 {
			m[kv[i]] = kv[i+1]
		}
		return m
	},

	"param": func(r *http.Request, key string) string {
		value, _ := url.QueryUnescape(r.PathValue(key))
		return value
	},

	"setting": func(golfer *golfer.Golfer, page, id string) any {
		if golfer != nil {
			return golfer.Settings[page][id]
		}

		for _, setting := range config.Settings[page] {
			if setting.ID == id {
				return setting.Default
			}
		}

		return nil
	},

	"svg": func(name string, attrs ...string) (template.HTML, error) {
		type Use struct {
			Href string `xml:"href,attr"`
		}

		type SVG struct {
			XMLName xml.Name   `xml:"svg"`
			Attrs   []xml.Attr `xml:",attr"`
			Title   string     `xml:"title,omitempty"`
			Use     Use        `xml:"use"`
		}

		path := config.Assets["svg/"+name+".svg"]
		svg := SVG{Use: Use{Href: path + "#a"}}

		for i := 0; i < len(attrs); i += 2 {
			if attrs[i] == "title" {
				svg.Title = attrs[i+1]
			} else {
				svg.Attrs = append(svg.Attrs, xml.Attr{
					Name:  xml.Name{Local: attrs[i]},
					Value: attrs[i+1],
				})
			}
		}

		b, err := xml.Marshal(svg)
		return template.HTML(b), err
	},
})

func init() {
	// Manually walk & parse HTML templates so we can support directories.
	if err := fs.WalkDir(views, "html", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			data, err := views.ReadFile(path)
			if err != nil {
				return err
			}

			// Trim the directory prefix and extension suffix.
			name := path[len("html/") : len(path)-len(".html")]
			tmpl = template.Must(tmpl.New(name).Parse(string(data)))
		}

		return err
	}); err != nil {
		panic(err)
	}
}

func Render(w http.ResponseWriter, name string, data any) error {
	return tmpl.ExecuteTemplate(w, name, data)
}
