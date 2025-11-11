package ordered

// Based on https://pkg.go.dev/encoding/json/v2#example-package-OrderedObject

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
)

type Map []item

type item struct{ Key, Value any }

func (m *Map) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	if k := dec.PeekKind(); k != '{' {
		return fmt.Errorf("expected object start, but encountered %v", k)
	}
	if _, err := dec.ReadToken(); err != nil {
		return err
	}

	for dec.PeekKind() != '}' {
		var i item
		if err := json.UnmarshalDecode(dec, &i.Key); err != nil {
			return err
		}
		if err := json.UnmarshalDecode(dec, &i.Value); err != nil {
			return err
		}
		*m = append(*m, i)
	}

	_, err := dec.ReadToken()
	return err
}
