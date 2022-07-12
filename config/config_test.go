package config

import "testing"

func TestID(t *testing.T) {
	for _, tt := range []struct{ name, id string }{
		{"><>", "fish"},
		{"C#", "c-sharp"},
		{"C++", "cpp"},
		{"DON’T PANIC!", "dont-panic"},
		{"Happy Birthday, Code Golf", "happy-birthday-code-golf"},
		{"May the 4ᵗʰ Be with You", "may-the-4ᵗʰ-be-with-you"},
		{"Morse (Decoder)", "morse-decoder"},
		{"My God, It’s Full of Stars", "my-god-its-full-of-stars"},
		{"Off-the-grid", "off-the-grid"},
		{"tl;dr", "tl-dr"},
		{"λ", "λ"},
		{"π", "π"},
		{"τ", "τ"},
		{"φ", "φ"},
		{"√2", "√2"},
		{"𝑒", "𝑒"},
	} {
		if got := ID(tt.name); got != tt.id {
			t.Errorf("ID(%v) = %v; want %v", tt.name, got, tt.id)
		}
	}
}
