package db

import (
	"regexp"
	"strings"

	"github.com/vinovest/sqlx"
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
