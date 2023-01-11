package mongodb

import (
	"context"
	"time"

	// Packages
	trace "github.com/mutablelogic/go-accessory/pkg/trace"
	options "go.mongodb.org/mongo-driver/mongo/options"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
	. "github.com/mutablelogic/go-accessory"
)

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Delete zero or one documents and returns the number of deleted documents (which should be
// zero or one. The filter argument is used to determine a document to delete. If there is more than
// one filter, they are ANDed together
func (collection *collection) Delete(ctx context.Context, filter ...Filter) (int64, error) {
	// Check for collection
	if collection.Collection == nil {
		return -1, ErrOutOfOrder
	}

	// Return error it no filter is provided
	if len(filter) == 0 {
		return -1, ErrBadParameter.With("no filter argument provided")
	}

	// Trace
	ctx, matched, modified := trace.WithCollection(ctx, trace.OpDelete, collection.Database().Name(), collection.Name())
	defer trace.Do(ctx, collection.traceFn, time.Now())

	// Perform the delete
	result, err := collection.Collection.DeleteOne(ctx, and(filter...), &options.DeleteOptions{})
	if err != nil {
		return -1, err
	} else {
		*matched, *modified = result.DeletedCount, result.DeletedCount
		return result.DeletedCount, nil
	}
}

// DeleteMany deletes zero or more documents and returns the number of deleted documents.
func (collection *collection) DeleteMany(context.Context, ...Filter) (int64, error) {
	return 0, ErrNotImplemented
}
