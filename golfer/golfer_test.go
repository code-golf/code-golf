package golfer

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestEarn(t *testing.T) {
	db, mock, _ := sqlmock.New()

	golfer := Golfer{ID: 123}
	cheevos := []string{"foo", "bar", "bar", "baz"}

	for _, cheevo := range cheevos {
		mock.ExpectExec("INSERT INTO trophies").
			WithArgs(golfer.ID, cheevo).
			WillReturnResult(sqlmock.NewResult(1, 1))
	}

	for _, cheevo := range cheevos {
		golfer.Earn(db, cheevo)
	}

	want := []string{"bar", "baz", "foo"}
	if !reflect.DeepEqual(golfer.Cheevos, want) {
		t.Errorf("golfer.Cheevos = %v; want %v", golfer.Cheevos, want)
	}
}
