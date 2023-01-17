package query

import (
	// Namespace imports
	. "github.com/mutablelogic/go-accessory"
)

//. "github.com/mutablelogic/go-accessory/pkg/sqlite/quote"

///////////////////////////////////////////////////////////////////////////////
// TYPES

type createTable struct {
	name
	col []*name
}

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (n *name) CreateTable(c ...Name) CreateTable {
	cols := make([]*name, 0, len(c))
	for _, v := range c {
		if v, ok := v.(*name); ok {
			cols = append(cols, v)
		}
	}
	return &createTable{name{query: n.query, schema: n.schema}, cols}
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

// Query returns the SQL query that can be executed
func (table *createTable) Query() string {
	var str string
	str += "CREATE TABLE " + table.name.SchemaName()
	return str
}
