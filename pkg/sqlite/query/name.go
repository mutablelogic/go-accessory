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
	schema   string
	alias    string // For use in FROM clauses
	decltype string // For use in CREATE TABLE clauses
}

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

// N returns a new table source
func N(v string, f ...QueryFlag) Name {
	name := new(name)
	name.v = v
	for _, v := range f {
		name.f |= v
	}
	return name
}

///////////////////////////////////////////////////////////////////////////////
// METHODS

func (n *name) WithSchema(schema string) Name {
	return &name{n.query, schema, n.alias, n.decltype}
}

func (n *name) WithAlias(alias string) Name {
	return &name{n.query, n.schema, alias, n.decltype}
}

func (n *name) WithType(decltype string) Name {
	return &name{n.query, n.schema, n.alias, decltype}
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

// SchemaName returns the schema.name clause
func (name *name) SchemaName() string {
	if name.schema != "" {
		return QuoteIdentifier(name.schema) + "." + QuoteIdentifier(name.v)
	} else {
		return QuoteIdentifier(name.v)
	}
}

// Query returns the SQL query that can be executed
func (name *name) Query() string {
	var str string
	str += name.SchemaName()
	if name.alias != "" {
		str += " AS " + QuoteIdentifier(name.alias)
	}
	if name.decltype != "" {
		str += " " + name.decltype
	}

	// For column definitions, add in NOT NULL, UNIQUE or PRIMARY KEY clauses
	// For sort clause, add in ASC or DESC
	if name.f.Is(NOT_NULL) {
		for _, v := range []QueryFlag{NOT_NULL, QUERY_CONFLICT} {
			if name.f.Is(v) {
				str += " " + v.String()
			}
		}
	} else if name.f.Is(UNIQUE_KEY) {
		for _, v := range []QueryFlag{UNIQUE_KEY, QUERY_CONFLICT} {
			if name.f.Is(v) {
				str += " " + v.String()
			}
		}
	} else if name.f.Is(PRIMARY_KEY) {
		for _, v := range []QueryFlag{PRIMARY_KEY, QUERY_SORT, QUERY_CONFLICT, AUTO_INCREMENT} {
			if name.f.Is(v) {
				str += " " + v.String()
			}
		}
	} else if name.f.Is(QUERY_SORT) {
		if name.f.Is(ASC) {
			str += " " + ASC.String()
		} else if name.f.Is(DESC) {
			str += " " + DESC.String()
		}
	}

	// Return success
	return str
}
