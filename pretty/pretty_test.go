package pretty

import (
	"html/template"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBytes(t *testing.T) {
	for _, tt := range []struct {
		want string
		b    int
	}{
		{"1.0 B", 1},
		{"1.0 KiB", 1024},
		{"1.0 MiB", 1024 * 1024},
		{"1.0 GiB", 1024 * 1024 * 1024},
	} {
		assert.Equal(t, tt.want, Bytes(tt.b))
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
		{"1,234,567", 1234567},
	} {
		assert.Equal(t, tt.want, Comma(tt.i))
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
		assert.Equal(t, tt.want, Ordinal(tt.i))
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
		assert.Contains(t, Time(tt.t), ">"+tt.want+"<")
	}
}
