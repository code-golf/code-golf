package golfer

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func TestEarn(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	db := sqlx.NewDb(mockDB, "sqlmock")

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
