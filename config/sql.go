package config

import (
	"database/sql/driver"
	"strings"
)

type (
	Cheevos []*Cheevo
	Holes   []*Hole
	Langs   []*Lang
)

type NullCountry struct {
	Country *Country
	Valid   bool
}

type NullHole struct {
	Hole  *Hole
	Valid bool
}

type NullLang struct {
	Lang  *Lang
	Valid bool
}

func (n *NullCountry) Scan(id any) error {
	n.Country, n.Valid = CountryByID[asString(id)]
	return nil
}

func (n *NullHole) Scan(id any) error {
	n.Hole, n.Valid = HoleByID[asString(id)]
	return nil
}

func (n *NullLang) Scan(id any) error {
	n.Lang, n.Valid = LangByID[asString(id)]
	return nil
}

func (c *Cheevo) Scan(src any) error { return scanID(c, src, CheevoByID) }
func (h *Hole) Scan(src any) error   { return scanID(h, src, AllHoleByID) }
func (l *Lang) Scan(src any) error   { return scanID(l, src, AllLangByID) }

func (c *Cheevos) Scan(src any) error { return scanIDs(c, src, CheevoByID) }
func (h *Holes) Scan(src any) error   { return scanIDs(h, src, AllHoleByID) }
func (l *Langs) Scan(src any) error   { return scanIDs(l, src, AllLangByID) }

func (c Cheevo) Value() (driver.Value, error) { return c.ID, nil }
func (h Hole) Value() (driver.Value, error)   { return h.ID, nil }
func (l Lang) Value() (driver.Value, error)   { return l.ID, nil }

func asString(src any) (s string) {
	switch v := src.(type) {
	case []byte:
		s = string(v)
	case string:
		s = v
	}
	return
}

func scanID[T any](thing *T, src any, lookup map[string]*T) error {
	*thing = *lookup[asString(src)]
	return nil
}

// Very crude pg array parsing, only one-dimensional arrays.
// See https://github.com/lib/pq/blob/master/array.go for proper parsing.
func scanIDs[E any, S ~[]*E](things *S, src any, lookup map[string]*E) error {
	if ids := asString(src); len(ids) > 2 {
		for id := range strings.SplitSeq(ids[1:len(ids)-1], ",") {
			*things = append(*things, lookup[id])
		}
	}
	return nil
}
