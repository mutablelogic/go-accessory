package query

import (
	// Namespace imports
	. "github.com/mutablelogic/go-accessory"
	. "github.com/mutablelogic/go-accessory/pkg/sqlite/quote"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type name struct {
	query
	schema string
	alias  string // For use in FROM clauses
	desc   bool   // For use in ORDER BY clauses
}

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

// N returns a new table source
func N(v string) Name {
	name := new(name)
	name.v = v
	return name
}

///////////////////////////////////////////////////////////////////////////////
// METHODS

func (n *name) WithSchema(schema string) Name {
	return &name{query{v: n.v}, schema, n.alias, n.desc}
}

func (n *name) WithAlias(alias string) Name {
	return &name{query{v: n.v}, n.schema, alias, n.desc}
}

func (n *name) WithDesc() Name {
	return &name{query{v: n.v}, n.schema, n.alias, true}
}

func (n *name) Column(v string) Column {
	return newcolumn(n, v)
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

// Query returns the SQL query that can be executed
func (name *name) Query() string {
	var str string
	if name.schema != "" {
		str += QuoteIdentifier(name.schema) + "." + QuoteIdentifier(name.v)
	} else {
		str += QuoteIdentifier(name.v)
	}
	if name.alias != "" {
		str += " AS " + QuoteIdentifier(name.alias)
	}
	if name.desc {
		str += " DESC"
	}
	return str
}
