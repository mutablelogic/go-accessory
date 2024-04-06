/*
N(...).CreateSchema() is a factory method that returns a new create schema struct:

	N("t").CreateSchema() => "CREATE SCHEMA t"
	N("t").CreateSchema().IfNotExists() => "CREATE SCHEMA IF NOT EXISTS t"
*/
package querybuilder

///////////////////////////////////////////////////////////////////////////////
// TYPES

type createSchema struct {
	flags
	name
}

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (n name) CreateSchema() createSchema {
	return createSchema{name: n}
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Indicate that the scheme should not be created if it already exists
func (q createSchema) IfNotExists() createSchema {
	q.flags |= ifNotExists
	return q
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (q createSchema) String() string {
	return join("CREATE SCHEMA", (q.flags & ifNotExists), q.name.Name())
}
