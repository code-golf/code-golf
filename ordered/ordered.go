package ordered

import (
	"encoding/json"
	"sort"
)

type Map []item
type item struct {
	Key   string
	Value interface{}
	pos   int
}

func (m *Map) UnmarshalJSON(b []byte) error {
	var items map[string]item
	if err := json.Unmarshal(b, &items); err != nil {
		return err
	}

	for key, item := range items {
		item.Key = key

		*m = append(*m, item)
	}

	sort.Slice(*m, func(i, j int) bool { return (*m)[i].pos < (*m)[j].pos })

	return nil
}

var pos int

func (i *item) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &i.Value); err != nil {
		return err
	}

	i.pos = pos
	pos++

	return nil
}
