/*
N(...).T(string) is a factory method that returns a new column, or a specified type:

	N("col").T("type") => "col TEXT"
	N("col").T("type").NotNull() => "col TEXT NOT NULL"
	N("col").T("type").Unique() => "col TEXT UNIQUE"
	N("col").T("type").PrimaryKey() => "col TEXT PRIMARY KEY"
*/
package querybuilder

import "strings"

///////////////////////////////////////////////////////////////////////////////
// TYPES

type column struct {
	name
	decltype string
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

// Indicate that the column must be unique
func (q column) Unique() column {
	q.flags |= unique
	return q
}

// Indicate that the column is part of the primary key
func (q column) PrimaryKey() column {
	q.flags |= primarykey
	return q
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (q column) String() string {
	return join(q.name.SchemaName(), q.decltype, (q.flags & notnull), (q.flags & unique), (q.flags & primarykey))
}
