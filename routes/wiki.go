package routes

import (
	"database/sql"
	"errors"
	"html/template"
	"net/http"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/session"
)

// GET /wiki/*
func wikiGET(w http.ResponseWriter, r *http.Request) {
	data := struct {
		HTML  template.HTML
		Nav   *config.Navigaton
		Title string
	}{Title: "Wiki"}

	// Page (if we have a slug).
	if slug := param(r, "*"); slug != "" {
		if err := session.Database(r).QueryRow(
			"SELECT html, 'Wiki: ' || name FROM wiki WHERE slug = $1", slug,
		).Scan(&data.HTML, &data.Title); errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			panic(err)
		}
	}

	// Navigation.
	var pages []struct{ Name, Section, Slug string }
	if err := session.Database(r).Select(
		&pages, "SELECT name, section, slug FROM wiki ORDER BY section, name",
	); err != nil {
		panic(err)
	}

	groups := []*config.LinkGroup{
		{Links: []*config.NavLink{{Name: "Home", Path: "/wiki"}}, Slug: "*"},
	}
	for _, page := range pages {
		if groups[len(groups)-1].Name != page.Section {
			groups = append(groups, &config.LinkGroup{Name: page.Section})
		}

		groups[len(groups)-1].Links = append(
			groups[len(groups)-1].Links,
			&config.NavLink{Name: page.Name, Path: "/wiki/" + page.Slug, Slug: page.Slug},
		)
	}
	data.Nav = &config.Navigaton{Groups: groups, OnePerRow: true}

	render(w, r, "wiki", data, data.Title)
}
