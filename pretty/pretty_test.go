package pretty

import "testing"

func TestComma(t *testing.T) {
	for _, tt := range []struct {
		want string
		i    int
	}{
		{"1", 1},
		{"12", 12},
		{"123", 123},
		{"1,234", 1234},
		{"12,345", 12345},
		{"123,456", 123456},
	} {
		if got := Comma(tt.i); got != tt.want {
			t.Errorf("Comma(%v) = %v; want %v", tt.i, got, tt.want)
		}
	}
}

func TestOrdinal(t *testing.T) {
	for _, tt := range []struct {
		want string
		i    int
	}{
		{"st", 1},
		{"nd", 2},
		{"rd", 3},
		{"th", 4},
		{"th", 11},
		{"th", 12},
		{"th", 13},
		{"st", 101},
		{"nd", 102},
		{"rd", 103},
		{"th", 111},
		{"th", 112},
		{"th", 113},
	} {
		if got := Ordinal(tt.i); got != tt.want {
			t.Errorf("Ordinal(%v) = %v; want %v", tt.i, got, tt.want)
		}
	}
}
