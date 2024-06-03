package null

import "testing"

func TestNullIfEmpty(t *testing.T) {
	// int
	if got := NullIfZero(0); got.Valid {
		t.Errorf(`NullIfZero(0).Valid = %v; want false`, got.Valid)
	}
	if got := NullIfZero(123); !got.Valid {
		t.Errorf(`NullIfZero(123).Valid = %v; want true`, got.Valid)
	}

	// string
	if got := NullIfZero(""); got.Valid {
		t.Errorf(`NullIfZero("").Valid = %v; want false`, got.Valid)
	}
	if got := NullIfZero("foo"); !got.Valid {
		t.Errorf(`NullIfZero("foo").Valid = %v; want true`, got.Valid)
	}
}
