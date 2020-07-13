package pretty

import (
	"html/template"
	"testing"
	"time"
)

func TestBytes(t *testing.T) {
	for _, tt := range []struct {
		want string
		b    int
	}{
		{"1.0 B", 1},
		{"1.0 KiB", 1024},
		{"1.0 MiB", 1024 * 1024},
	} {
		if got := Bytes(tt.b); got != tt.want {
			t.Errorf("Bytes(%v) = %v; want %v", tt.b, got, tt.want)
		}
	}
}

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

func TestTime(t *testing.T) {
	for _, tt := range []struct {
		want template.HTML
		t    time.Time
	}{
		{"a min ago", time.Now().UTC()},
		{"7 Jun 1989", time.Date(1989, time.June, 7, 0, 0, 0, 0, time.UTC)},
	} {
		got := Time(tt.t)
		if got = got[63 : len(got)-len("</time>")]; got != tt.want {
			t.Errorf("Ordinal(%v) = %v; want %v", tt.t, got, tt.want)
		}
	}
}
