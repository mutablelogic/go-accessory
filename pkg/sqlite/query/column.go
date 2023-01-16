package query

import (
	// Namespace imports
	. "github.com/mutablelogic/go-accessory"
	. "github.com/mutablelogic/go-accessory/pkg/sqlite/quote"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type column struct {
	name
	decltype string
	notnull  bool
	primary  bool
	auto     bool
}

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	defaultColumnDecltype = "TEXT"
)

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func newcolumn(n *name, decltype string) *column {
	column := new(column)
	column.name = *n
	column.decltype = decltype
	return column
}

///////////////////////////////////////////////////////////////////////////////
// METHODS

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

// Query returns the SQL query that can be executed
func (column *column) Query() string {
	str := QuoteIdentifier(column.name.v)
	if column.decltype != "" {
		str += " " + QuoteDeclType(column.decltype)
	} else {
		str += " " + defaultColumnDecltype
	}
	return str
}
