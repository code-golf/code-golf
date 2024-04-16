package null

import (
	"database/sql"
	"encoding/json"
	"time"
)

type (
	Bool   = Null[bool]
	Int    = Null[int]
	String = Null[string]
	Time   = Null[time.Time]

	Null[T any] struct{ sql.Null[T] }
)

func New[T any](value T, valid bool) Null[T] {
	return Null[T]{sql.Null[T]{V: value, Valid: valid}}
}

func (n Null[T]) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(n.V)
}
