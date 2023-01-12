package mongodb

import (
	"context"
	"errors"
	"reflect"
	"time"

	// Packages
	trace "github.com/mutablelogic/go-accessory/pkg/trace"
	driver "go.mongodb.org/mongo-driver/mongo"
	options "go.mongodb.org/mongo-driver/mongo/options"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
	. "github.com/mutablelogic/go-accessory"
)

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Find selects a single document based on filter and sort parameters.
// It returns ErrNotFound if no document is found
func (collection *collection) Find(ctx context.Context, sort Sort, filter ...Filter) (any, error) {
	// Check for collection
	if collection.Collection == nil {
		return nil, ErrOutOfOrder
	}

	// Trace
	ctx, matched, _ := trace.WithCollection(ctx, trace.OpFind, collection.Database().Name(), collection.Name())
	defer trace.Do(ctx, collection.traceFn, time.Now())

	// Do the find
	result := collection.Collection.FindOne(ctx, and(filter...), &options.FindOneOptions{
		Sort: sortdoc(sort),
	})

	// Check for errors
	if err := result.Err(); err != nil {
		if errors.Is(result.Err(), driver.ErrNoDocuments) {
			return nil, ErrNotFound
		} else {
			return nil, err
		}
	} else {
		*matched = 1
	}

	// Create a new document
	doc := reflect.New(collection.meta.Type).Interface()
	if err := result.Decode(doc); err != nil {
		return nil, err
	} else {
		return doc, nil
	}
}

// FindMany returns an iterable cursor based on filter and sort parameters.
// It returns ErrNotFound if no document is found
func (collection *collection) FindMany(ctx context.Context, sort Sort, filter ...Filter) (Cursor, error) {
	// Check for collection
	if collection.Collection == nil {
		return nil, ErrOutOfOrder
	}

	// Trace
	ctx, _, _ = trace.WithCollection(ctx, trace.OpFindMany, collection.Database().Name(), collection.Name())
	defer trace.Do(ctx, collection.traceFn, time.Now())

	// Do the find
	cur, err := collection.Collection.Find(ctx, and(filter...), &options.FindOptions{
		Sort:  sortdoc(sort),
		Limit: sortlimit(sort),
	})

	// Check for errors
	if errors.Is(err, driver.ErrNoDocuments) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	// Return the cursor
	return NewCursor(cur, collection.meta.Type), nil
}
