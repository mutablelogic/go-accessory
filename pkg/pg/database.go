package pg

import (
	"context"
	"fmt"
	"reflect"

	// Package imports
	// Namespace imports
	. "github.com/djthorpe/go-errors"
	. "github.com/mutablelogic/go-accessory"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type database struct {
	name, schema string
}

// Ensure *database implements the Database interface
var _ Database = (*database)(nil)

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

// Create a new database object from a connection - the database should already
// exist
func (conn *conn) new_database(name, schema string) *database {
	database := new(database)
	database.name = name
	database.schema = schema
	return database
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (database *database) String() string {
	str := "<pg.database"
	str += fmt.Sprintf(" name=%q", database.name)
	return str + ">"
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (database *database) Name() string {
	return database.name
}

func (database *database) Collection(doc any) Collection {
	// TODO: Cache the collection
	if doc == nil {
		return nil
	} else if collection := database.new_collection(reflect.ValueOf(doc), database.schema); collection == nil {
		return nil
	} else {
		return collection
	}
}

func (database *database) Insert(context.Context, ...any) error {
	return ErrNotImplemented
}
