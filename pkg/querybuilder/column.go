/*
N(...).T(string) is a factory method that returns a new column, or a specified type:

	N("col").T("type") => "col TEXT"
	N("col").T("type").NotNull() => "col TEXT NOT NULL"
	N("col").T("type").Unique() => "col TEXT UNIQUE"
	N("col").T("type").PrimaryKey() => "col TEXT PRIMARY KEY"
	N("col").T("type").Foreign("other_table") => "col TEXT REFERENCES other_table"
	N("col").T("type").Foreign("other_table","a") => "col TEXT REFERENCES other_table (a)"
*/
package querybuilder

import (
	"strings"

	// Package imports
	quote "github.com/mutablelogic/go-accessory/pkg/querybuilder/quote"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type column struct {
	name
	decltype string
	def      string
	key      any
	flags
}

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (n name) T(decltype string) column {
	return column{name: n, decltype: strings.ToUpper(decltype)}
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Indicate that the coliumn cannot contain a NULL value
func (q column) NotNull() column {
	q.flags |= notnull
	return q
}

// Indicate that the column has a unique constraint
func (q column) Unique() column {
	q.key = Key().Unique()
	return q
}

// Indicate that the column has a primary key constraint
func (q column) PrimaryKey() column {
	q.key = Key()
	q.flags |= primarykey
	return q
}

// Indicate that the column has a foreign key constraint
func (q column) Foreign(name string, column ...string) column {
	q.key = Key().Foreign(name, column...)
	return q
}

// The default value on insert
func (q column) Default(v any) column {
	q.def = quote.Join("DEFAULT", v)
	return q
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (q column) String() string {
	// Remove NOT NULL if PRIMARY KEY (implied)
	if q.flags.Is(primarykey) {
		q.flags &= ^notnull
	}
	return quote.Join(q.name.SchemaName(), q.decltype, q.key, q.def)
}
