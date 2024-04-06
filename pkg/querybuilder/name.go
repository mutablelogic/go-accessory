/*
N() is a factory method that returns a new name struct. The name struct can have schema, alias
and decltype fields:

	N("table").String() => "table"
	N("table").WithSchema("public").String() => "schema.table"
	N("table").As("a").String() => "table AS a"
	N("table").WithType("int").String() => "table INT"
*/
package querybuilder

import (
	quote "github.com/mutablelogic/go-accessory/pkg/querybuilder/quote"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type name struct {
	name   string
	schema string
	alias  string // For use in FROM clauses
}

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

// N returns a new table source, or tabel column source
func N(v string) name {
	return name{name: v}
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// WithSchema returns a new name with the schema set
func (n name) WithSchema(v string) name {
	return name{name: n.name, schema: v, alias: n.alias}
}

// As returns a new name with the alias set
func (n name) As(v string) name {
	return name{name: n.name, schema: n.schema, alias: v}
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

// String returns all elements quoted
func (n *name) String() string {
	str := n.SchemaName()
	if n.alias != "" {
		str += " AS " + quote.Identifier(n.alias)
	}
	return str
}

// SchemaName returns the quoted name of the table with optional schema component
func (n *name) SchemaName() string {
	str := quote.Identifier(n.name)
	if n.schema != "" {
		str = quote.Identifier(n.schema) + "." + str
	}
	return str
}

// Name returns the quoted name component only
func (n *name) Name() string {
	return quote.Identifier(n.name)
}
