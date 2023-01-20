package db

import "testing"

func TestMixedCaseToSnakeCase(t *testing.T) {
	for _, tt := range []struct{ mixed, snake string }{
		{"foo", "foo"},
		{"Foo", "foo"},
		{"fooBarBaz", "foo_bar_baz"},
		{"FooBarBaz", "foo_bar_baz"},

		{"UserID", "user_id"},
		{"UserIDForeign", "user_id_foreign"},
	} {
		if got := mixedCaseToSnakeCase(tt.mixed); got != tt.snake {
			t.Errorf("mixedCaseToSnakeCase(%v) = %v; want %v", tt.mixed, got, tt.snake)
		}
	}
}
