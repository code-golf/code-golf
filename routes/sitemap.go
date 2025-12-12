package routes

import (
	"encoding/xml"
	"net/http"
	"net/url"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/session"
)

// GET /sitemap.xml
func sitemapGET(w http.ResponseWriter, r *http.Request) {
	type URL struct {
		Loc string `xml:"loc"`
	}

	sitemap := struct {
		XMLName xml.Name `xml:"http://www.sitemaps.org/schemas/sitemap/0.9 urlset"`
		URLs    []URL    `xml:"url"`
	}{
		URLs: []URL{
			{"https://code.golf"},
			{"https://code.golf/about"},
			{"https://code.golf/ideas"},
			{"https://code.golf/stats"},
		},
	}

	var urls []URL
	if err := session.Database(r).Select(
		&urls,
		` SELECT concat_ws('/', 'https://code.golf/wiki', nullif(slug, '')) loc
		    FROM wiki
		ORDER BY slug`,
	); err != nil {
		panic(err)
	}
	sitemap.URLs = append(sitemap.URLs, urls...)

	for _, hole := range config.HoleList {
		sitemap.URLs = append(
			sitemap.URLs, URL{"https://code.golf/" + url.PathEscape(hole.ID)})
	}

	w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	w.Write([]byte(xml.Header))

	if err := xml.NewEncoder(w).Encode(sitemap); err != nil {
		panic(err)
	}
}
