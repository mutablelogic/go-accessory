package query

import (
	// Namespace imports
	. "github.com/mutablelogic/go-accessory"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type query struct {
	v string
	f QueryFlag
}

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

// Q returns a new ad-hoc query
func Q(v string) *query {
	return &query{v, NONE}
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Query returns the SQL query that can be executed
func (q *query) Query() string {
	return q.v
}

// Append flags for any query
func (q *query) With(f QueryFlag) Query {
	return &query{q.v, q.f | f}
}
