package config

import "database/sql/driver"

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

func (c *Cheevo) Scan(id any) error {
	*c = *CheevoByID[asString(id)]
	return nil
}

func (h *Hole) Scan(id any) error {
	*h = *AllHoleByID[asString(id)]
	return nil
}

func (h *Hole) Value() (driver.Value, error) { return h.ID, nil }

func (l *Lang) Scan(id any) error {
	*l = *LangByID[asString(id)]
	return nil
}

func asString(src any) (s string) {
	switch v := src.(type) {
	case []byte:
		s = string(v)
	case string:
		s = v
	}
	return
}
