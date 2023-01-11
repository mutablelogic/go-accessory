package mongodb

import (
	"context"

	// Package imports
	driver "go.mongodb.org/mongo-driver/mongo"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
	. "github.com/mutablelogic/go-accessory"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type collection struct {
	*driver.Collection
}

// Ensure *collection implements the Collection interface
var _ Collection = (*collection)(nil)

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewCollection(database *driver.Database, name string) *collection {
	return &collection{
		Collection: database.Collection(name),
	}
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Return the name of the collection
func (collection *collection) Name() string {
	if collection.Collection == nil {
		return ""
	}
	return collection.Name()
}

// Delete zero or one documents and returns the number of deleted documents (which should be
// zero or one. The filter argument is used to determine a document to delete. If there is more than
// one filter, they are ANDed together
func (collection *collection) Delete(context.Context, ...Filter) (int64, error) {
	return 0, ErrNotImplemented
}

// DeleteMany deletes zero or more documents and returns the number of deleted documents.
func (collection *collection) DeleteMany(context.Context, ...Filter) (int64, error) {
	return 0, ErrNotImplemented
}

// Find selects a single document based on filter and sort parameters.
// It returns ErrNotFound if no document is found
func (collection *collection) Find(context.Context, Sort, ...Filter) (any, error) {
	return nil, ErrNotImplemented
}

// FindMany returns an iterable cursor based on filter and sort parameters.
// It returns ErrNotFound if no document is found
func (collection *collection) FindMany(context.Context, Sort, ...Filter) (Cursor, error) {
	return nil, ErrNotImplemented
}

// Update zero or one document with given values and return the number
// of documents matched and modified, neither of which should be more than one.
func (collection *collection) Update(context.Context, any, ...Filter) (int64, int64, error) {
	return 0, 0, ErrNotImplemented
}

// Update zero or more document with given values and return the number
// of documents matched and modified, neither of which should be more than one.
func (collection *collection) UpdateMany(context.Context, any, ...Filter) (int64, int64, error) {
	return 0, 0, ErrNotImplemented
}
