package config

import (
	"net/http"
	"net/url"
	"regexp"
	"slices"
	"strings"

	"github.com/go-chi/chi/v5"
)

type NavLink struct{ Emoji, Name, Slug, Path string }

type LinkGroup struct {
	Links      []*NavLink
	Name, Slug string
}

type Navigaton struct {
	Path      string
	Groups    []*LinkGroup
	OnePerRow bool
}

var Nav map[string]*Navigaton

// TODO OnePerRow needs a better name because it also means a single dropdown,
// maybe something like unified/single namespace.
func init() {
	Nav = map[string]*Navigaton{
		"rankings/cheevos": {
			OnePerRow: true,
			Path:      "/rankings/cheevos/{cheevo}",
			Groups: append(
				[]*LinkGroup{group("", "cheevo", "All")},
				groupsCheevos()...,
			),
		},

		"rankings/holes": {
			Path: "/rankings/holes/{hole}/{lang}/{scoring}",
			Groups: []*LinkGroup{
				group("Scoring", "scoring", "Bytes", "Chars"),
				groupLangs(),
				groupHoles(),
			},
		},

		"rankings/langs": {
			Path: "/rankings/langs/{lang}/{scoring}",
			Groups: []*LinkGroup{
				group("Scoring", "scoring", "Bytes", "Chars"),
				groupLangs(),
			},
		},

		"rankings/medals": {
			Path: "/rankings/medals/{hole}/{lang}/{scoring}",
			Groups: []*LinkGroup{
				group("Scoring", "scoring", "All", "Bytes", "Chars"),
				groupLangs(),
				groupHoles(),
			},
		},

		"rankings/misc": {
			OnePerRow: true,
			Path:      "/rankings/misc/{type}",
			Groups: []*LinkGroup{
				group("", "type", "Diamond Deltas", "Followers",
					"Holes Authored", "Oldest Diamonds", "Referrals",
					"Solutions"),
			},
		},

		"rankings/recent-holes": {
			Path: "/rankings/recent-holes/{lang}/{scoring}",
			Groups: []*LinkGroup{
				group("Scoring", "scoring", "Bytes", "Chars"),
				groupLangs(),
			},
		},

		"recent/solutions": {
			Path: "/recent/solutions/{hole}/{lang}/{scoring}",
			Groups: []*LinkGroup{
				group("Scoring", "scoring", "Bytes", "Chars"),
				groupLangs(),
				groupHoles(),
			},
		},
	}

	// Set the link Path template and pre-fill the one slug we can.
	for _, nav := range Nav {
		for _, group := range nav.Groups {
			for _, link := range group.Links {
				link.Path = strings.ReplaceAll(
					nav.Path,
					"{"+group.Slug+"}",
					link.Slug,
				)
			}
		}
	}
}

var PathSlug = regexp.MustCompile("{[a-z]+}")

func (l NavLink) PopulatePath(r *http.Request) string {
	if Path := PathSlug.ReplaceAllStringFunc(l.Path, func(s string) string {
		value, _ := url.QueryUnescape(chi.URLParam(r, s[1:len(s)-1]))
		return value
	}); Path != r.URL.Path {
		return Path
	}

	return ""
}

func (n *Navigaton) ReverseGroups() []*LinkGroup {
	len := len(n.Groups)
	groups := make([]*LinkGroup, len)

	for i, group := range n.Groups {
		groups[len-i-1] = group
	}

	return groups
}

func group(name, slug string, linkNames ...string) *LinkGroup {
	group := LinkGroup{
		Links: make([]*NavLink, len(linkNames)),
		Name:  name,
		Slug:  slug,
	}

	for i, name := range linkNames {
		group.Links[i] = &NavLink{Name: name, Slug: ID(name)}
	}

	return &group
}

func groupsCheevos() (groups []*LinkGroup) {
	for category, cheevos := range CheevoTree {
		group := group(category, "cheevo")

		for _, cheevo := range cheevos {
			group.Links = append(group.Links, &NavLink{
				Emoji: cheevo.Emoji,
				Name:  cheevo.Name,
				Slug:  cheevo.ID,
			})
		}

		groups = append(groups, group)
	}

	slices.SortFunc(groups, func(a, b *LinkGroup) int {
		return strings.Compare(a.Name, b.Name)
	})

	return
}

func groupHoles() *LinkGroup {
	group := group("Hole", "hole", "All")

	for _, hole := range HoleList {
		group.Links = append(group.Links, &NavLink{
			Name: hole.Name,
			Slug: hole.ID,
		})
	}

	return group
}

func groupLangs() *LinkGroup {
	group := group("Language", "lang", "All")

	for _, lang := range LangList {
		group.Links = append(group.Links, &NavLink{
			Name: lang.Name,
			Slug: lang.ID,
		})
	}

	return group
}
