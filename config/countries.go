package config

var (
	CountryByID = map[string]*Country{}
	CountryTree map[string][]*Country
)

type Country struct{ ID, Flag, Name string }

type NullCountry struct {
	Country *Country
	Valid   bool
}

func (n *NullCountry) Scan(id any) error {
	n.Country, n.Valid = CountryByID[string(id.([]byte))]
	return nil
}

func init() {
	unmarshal("countries.toml", &CountryTree)

	for _, countries := range CountryTree {
		for _, country := range countries {
			for _, letter := range country.ID {
				country.Flag += string('🇦' - 'A' + letter)
			}

			CountryByID[country.ID] = country
		}
	}
}
