package mongodb

import (

	// Package imports
	trace "github.com/mutablelogic/go-accessory/pkg/trace"
	driver "go.mongodb.org/mongo-driver/mongo"

	// Namespace imports
	//. "github.com/djthorpe/go-errors"
	. "github.com/mutablelogic/go-accessory"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type collection struct {
	*driver.Collection

	meta    *meta
	traceFn trace.Func
}

// Ensure *collection implements the Collection interface
var _ Collection = (*collection)(nil)

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewCollection(database *driver.Database, meta *meta, fn trace.Func) *collection {
	// Check arguments
	if database == nil || meta == nil {
		return nil
	}
	// Return collection
	return &collection{
		Collection: database.Collection(meta.Name),
		meta:       meta,
		traceFn:    fn,
	}
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Return the name of the collection
func (collection *collection) Name() string {
	if collection.Collection == nil {
		return ""
	}
	return collection.Collection.Name()
}
