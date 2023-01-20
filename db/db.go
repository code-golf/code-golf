package db

import (
	"context"
	"database/sql"
	"regexp"
	"strings"

	"github.com/jmoiron/sqlx"
)

// mixedCase boundary, e.g. fo(o)(B)ar or SQ(L)(B)ar.
// TODO Handle digits if the need arises but digits are rare in columns names.
var boundary = regexp.MustCompile(`([a-z])([A-Z])|([A-Z])([A-Z][a-z])`)

func mixedCaseToSnakeCase(s string) string {
	return strings.ToLower(boundary.ReplaceAllString(s, "${1}${3}_${2}${4}"))
}

// Open a connection pool to the Database.
func Open() (db *sqlx.DB) {
	db = sqlx.MustOpen("postgres", "")
	db.MapperFunc(mixedCaseToSnakeCase)
	return
}

// Queryable includes all methods shared by sqlx.DB & sqlx.Tx, allowing either
// type to be used interchangeably. https://github.com/jmoiron/sqlx/pull/809
type Queryable interface {
	sqlx.ExecerContext
	sqlx.Ext
	sqlx.Preparer
	sqlx.PreparerContext
	sqlx.QueryerContext

	Get(any, string, ...any) error
	GetContext(context.Context, any, string, ...any) error
	MustExec(string, ...any) sql.Result
	MustExecContext(context.Context, string, ...any) sql.Result
	NamedExec(string, any) (sql.Result, error)
	NamedExecContext(context.Context, string, any) (sql.Result, error)
	NamedQuery(string, any) (*sqlx.Rows, error)
	PrepareNamed(string) (*sqlx.NamedStmt, error)
	PrepareNamedContext(context.Context, string) (*sqlx.NamedStmt, error)
	Preparex(string) (*sqlx.Stmt, error)
	PreparexContext(context.Context, string) (*sqlx.Stmt, error)
	QueryRow(string, ...any) *sql.Row
	QueryRowContext(context.Context, string, ...any) *sql.Row
	Select(any, string, ...any) error
	SelectContext(context.Context, any, string, ...any) error
}

var _ Queryable = (*sqlx.DB)(nil)
var _ Queryable = (*sqlx.Tx)(nil)
