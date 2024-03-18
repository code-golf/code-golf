package golfer

import (
	"slices"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

func TestEarn(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	db := sqlx.NewDb(mockDB, "sqlmock")

	golfer := Golfer{ID: 123}
	cheevos := pq.StringArray{"foo", "bar", "bar", "baz"}

	for _, cheevo := range cheevos {
		mock.ExpectExec("INSERT INTO trophies").
			WithArgs(golfer.ID, cheevo).
			WillReturnResult(sqlmock.NewResult(1, 1))
	}

	for _, cheevo := range cheevos {
		golfer.Earn(db, cheevo)
	}

	want := pq.StringArray{"bar", "baz", "foo"}
	if !slices.Equal(golfer.Cheevos, want) {
		t.Errorf("golfer.Cheevos = %v; want %v", golfer.Cheevos, want)
	}
}
