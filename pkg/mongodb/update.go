package mongodb

import (
	"context"
	"time"

	// Packages
	trace "github.com/mutablelogic/go-accessory/pkg/trace"
	bson "go.mongodb.org/mongo-driver/bson"
	options "go.mongodb.org/mongo-driver/mongo/options"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
	. "github.com/mutablelogic/go-accessory"
)

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Update zero or one document with given values and return the number
// of documents matched and modified, neither of which should be more than one.
func (collection *collection) Update(ctx context.Context, patch any, filter ...Filter) (int64, int64, error) {
	// Check for collection
	if collection.Collection == nil {
		return -1, -1, ErrOutOfOrder
	}

	// Return error it no filter is provided
	if len(filter) == 0 {
		return -1, -1, ErrBadParameter.With("no filter argument provided")
	}

	// Trace
	ctx, matched, modified := trace.WithCollection(ctx, trace.OpUpdate, collection.Database().Name(), collection.Name())
	defer trace.Do(ctx, collection.traceFn, time.Now())

	// Do the find
	result, err := collection.Collection.UpdateOne(ctx, and(filter...), bson.D{{"$set", patch}}, &options.UpdateOptions{})
	if err != nil {
		return -1, -1, err
	} else {
		*matched = result.MatchedCount
		*modified = result.ModifiedCount
		return result.MatchedCount, result.ModifiedCount, nil
	}
}

// Update zero or more document with given values and return the number
// of documents matched and modified, neither of which should be more than one.
func (collection *collection) UpdateMany(ctx context.Context, patch any, filter ...Filter) (int64, int64, error) {
	// Check for collection
	if collection.Collection == nil {
		return -1, -1, ErrOutOfOrder
	}

	// Return error it no filter is provided
	if len(filter) == 0 {
		return -1, -1, ErrBadParameter.With("no filter argument provided")
	}

	// Trace
	ctx, matched, modified := trace.WithCollection(ctx, trace.OpUpdateMany, collection.Database().Name(), collection.Name())
	defer trace.Do(ctx, collection.traceFn, time.Now())

	// Do the find
	result, err := collection.Collection.UpdateMany(ctx, and(filter...), bson.D{{"$set", patch}}, &options.UpdateOptions{})
	if err != nil {
		return -1, -1, err
	} else {
		*matched = result.MatchedCount
		*modified = result.ModifiedCount
		return result.MatchedCount, result.ModifiedCount, nil
	}
}
