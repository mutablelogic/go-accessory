/*
N(...).CreateTable() is a factory method that returns a new create table struct:

	N("t").CreateTable() => "CREATE TABLE t ()"
*/
package querybuilder

import (
	"fmt"
	"reflect"

	// Package imports
	quote "github.com/mutablelogic/go-accessory/pkg/querybuilder/quote"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type createTable struct {
	flags
	name
	column []any
}

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	defaultType = "TEXT"
)

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (n name) CreateTable(v ...any) createTable {
	q := createTable{name: n}
	for _, v := range v {
		switch v := v.(type) {
		case column:
			q.column = append(q.column, v)
		case string:
			q.column = append(q.column, N(v).T(defaultType))
		case name:
			q.column = append(q.column, v.T(defaultType))
		default:
			panic("Invalid type: " + fmt.Sprint(reflect.TypeOf(v)))
		}
	}
	return q
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Indicate that the table is temporary
func (q createTable) Temporary() createTable {
	q.flags |= temporary
	return q
}

// Indicate that the table should not be created if it already exists
func (q createTable) IfNotExists() createTable {
	q.flags |= ifNotExists
	return q
}

// Indicate that the table should be created as unlogged
func (q createTable) Unlogged() createTable {
	q.flags |= unlogged
	return q
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (q createTable) String() string {
	return quote.Join("CREATE", (q.flags & (temporary | unlogged)), "TABLE", (q.flags & ifNotExists), q.name.SchemaName(), "("+quote.JoinSep(",", q.column...)+")")
}
