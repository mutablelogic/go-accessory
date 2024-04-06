package pg

import (
	"context"
	"fmt"
	"reflect"

	// Package imports
	trace "github.com/mutablelogic/go-accessory/pkg/trace"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
	. "github.com/mutablelogic/go-accessory"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type database struct {
	conn       *conn
	name       string
	schema     string
	collection map[reflect.Type]Collection
}

// Ensure *database implements the Database interface
var _ Database = (*database)(nil)

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

// Create a new database object from a connection - the database should already exist
func (conn *conn) new_database(name, schema string) *database {
	// Create a new database object
	database := new(database)
	database.conn = conn
	database.name = name
	database.schema = schema
	database.collection = make(map[reflect.Type]Collection)

	// Create the schema if it doesn't exist
	if schema != "" {
		ctx := context.Background()
		if err := conn.CreateSchema(ctx, schema, true); err != nil {
			trace.Err(ctx, conn.tracefn, err)
			return nil
		}
	}

	// Return the database
	return database
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (database *database) String() string {
	str := "<pg.database"
	str += fmt.Sprintf(" name=%q", database.name)
	if len(database.collection) > 0 {
		str += fmt.Sprintf(" collection=%q", database.collection)
	}
	return str + ">"
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Return the name of the database
func (database *database) Name() string {
	return database.name
}

// Return a collection object for a document type
func (database *database) Collection(doc any) Collection {
	// Dereference
	v := reflect.ValueOf(doc)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.IsNil() || v.Kind() != reflect.Struct {
		return nil
	}

	// Return cached collection
	if collection, exists := database.collection[v.Type()]; exists {
		return collection
	}

	// Create a new collection
	ctx := context.Background()
	if collection := database.new_collection(ctx, v, database.schema); collection == nil {
		return nil
	} else {
		database.collection[v.Type()] = collection
		return collection
	}
}

// Insert documents of the same type to the database within a transaction.
// The document keys are updated if the document is writable.
func (database *database) Insert(ctx context.Context, doc ...any) error {
	return ErrNotImplemented
	/*
	   return conn.Do(ctx, func(ctx context.Context) error {

	   })
	*/
}
