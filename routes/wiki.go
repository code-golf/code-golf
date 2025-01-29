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
	var data struct {
		HTML        template.HTML
		Name, Title string
		Nav         *config.Navigaton
	}

	// Page.
	if err := session.Database(r).Get(
		&data, "SELECT html, name FROM wiki WHERE slug = $1", param(r, "*"),
	); errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		panic(err)
	}

	// Title.
	if data.Name == "Home" {
		data.Title = "Wiki"
	} else {
		data.Title = "Wiki: " + data.Name
	}

	// Navigation.
	var pages []struct{ Name, Section, Slug string }
	if err := session.Database(r).Select(
		&pages,
		"SELECT name, section, slug FROM wiki ORDER BY name != 'Home', section, name",
	); err != nil {
		panic(err)
	}

	groups := []*config.LinkGroup{{Slug: "*"}}
	for _, page := range pages {
		if groups[len(groups)-1].Name != page.Section {
			groups = append(groups, &config.LinkGroup{Name: page.Section})
		}

		path := "/wiki"
		if page.Slug != "" {
			path += "/" + page.Slug
		}

		groups[len(groups)-1].Links = append(
			groups[len(groups)-1].Links,
			&config.NavLink{Name: page.Name, Path: path, Slug: page.Slug},
		)
	}
	data.Nav = &config.Navigaton{Groups: groups, OnePerRow: true}

	render(w, r, "wiki", data, data.Title)
}
