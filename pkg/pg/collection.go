package pg

import (
	"context"
	"fmt"
	"reflect"

	// Package imports
	meta "github.com/mutablelogic/go-accessory/pkg/meta"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
	. "github.com/mutablelogic/go-accessory"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type collection struct {
	Meta   *meta.Collection
	Schema string
}

// Ensure *collection implements the Collection interface
var _ Collection = (*collection)(nil)

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	tagName = "bson"
)

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

// Create a new database object from a connection - the database should already
// exist
func (database *database) new_collection(r reflect.Value, schema string) *collection {
	if meta := meta.New(r, tagName); meta == nil {
		return nil
	} else {
		return &collection{
			Meta:   meta,
			Schema: schema,
		}
	}
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (c *collection) String() string {
	str := "<pg.collection"
	str += fmt.Sprintf(" name=%q", c.Meta.Name)
	str += fmt.Sprintf(" schema=%q", c.Schema)
	return str + ">"
}

// /////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Return the name of the collection
func (c *collection) Name() string {
	return c.Meta.Name
}

// Delete zero or one documents and returns the number of deleted documents (which should be
// zero or one. The filter argument is used to determine a document to delete. If there is more than
// one filter, they are ANDed together
func (c *collection) Delete(context.Context, ...Filter) (int64, error) {
	return 0, ErrNotImplemented
}

// DeleteMany deletes zero or more documents and returns the number of deleted documents.
func (c *collection) DeleteMany(context.Context, ...Filter) (int64, error) {
	return 0, ErrNotImplemented
}

// Find selects a single document based on filter and sort parameters.
// It returns ErrNotFound if no document is found
func (c *collection) Find(context.Context, Sort, ...Filter) (any, error) {
	return 0, ErrNotImplemented
}

// FindMany returns an iterable cursor based on filter and sort parameters.
// It returns ErrNotFound if no document is found
func (c *collection) FindMany(context.Context, Sort, ...Filter) (Cursor, error) {
	return nil, ErrNotImplemented
}

// Update zero or one document with given values and return the number
// of documents matched and modified, neither of which should be more than one.
func (c *collection) Update(context.Context, any, ...Filter) (int64, int64, error) {
	return 0, 0, ErrNotImplemented
}

// Update zero or more document with given values and return the number
// of documents matched and modified, neither of which should be more than one.
func (c *collection) UpdateMany(context.Context, any, ...Filter) (int64, int64, error) {
	return 0, 0, ErrNotImplemented
}

// FindUpdate selects a single document based on filter and sort parameters,
// updates the document with the given values and returns the document as it appeared
// before updating, or ErrNotFound if no document is found and updated.
func (c *collection) FindUpdate(context.Context, any, Sort, ...Filter) (any, error) {
	return 0, ErrNotImplemented
}
