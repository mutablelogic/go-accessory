package mongodb

import (
	"fmt"

	// Package imports
	trace "github.com/mutablelogic/go-accessory/pkg/trace"
	driver "go.mongodb.org/mongo-driver/mongo"

	// Namespace imports
	. "github.com/mutablelogic/go-accessory"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type database struct {
	*driver.Database

	metaFn  metaLookupFunc // Function to return collection metadata from prototypes
	traceFn trace.Func     // Function to trace operations
}

// Ensure *database implements the Database interface
var _ Database = (*database)(nil)

type metaLookupFunc func(...any) *meta

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewDatabase(conn *conn, name string, meta metaLookupFunc, trace trace.Func) *database {
	return &database{
		Database: conn.Client.Database(name),
		metaFn:   meta,
		traceFn:  trace,
	}
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (database *database) String() string {
	str := "<mongodb.database"
	if database.Database != nil {
		str += fmt.Sprintf(" name=%q", database.Name())
	}
	return str + ">"
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Return the name of the database
func (database *database) Name() string {
	if database.Database == nil {
		return ""
	} else {
		return database.Database.Name()
	}
}

// Return a collection
func (database *database) Collection(proto any) Collection {
	if meta := database.metaFn(proto); meta == nil {
		return nil
	} else {
		return NewCollection(database.Database, meta, database.traceFn)
	}
}

///////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (database *database) collectionForProtos(proto ...any) *collection {
	if meta := database.metaFn(proto...); meta == nil {
		return nil
	} else {
		return NewCollection(database.Database, meta, database.traceFn)
	}
}
