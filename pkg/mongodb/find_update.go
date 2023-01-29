package mongodb

import (
	"context"
	"errors"
	"reflect"
	"time"

	// Packages
	trace "github.com/mutablelogic/go-accessory/pkg/trace"
	bson "go.mongodb.org/mongo-driver/bson"
	driver "go.mongodb.org/mongo-driver/mongo"
	options "go.mongodb.org/mongo-driver/mongo/options"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
	. "github.com/mutablelogic/go-accessory"
)

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// FindUpdate selects a single document based on filter and sort parameters,
// updates the document with the given values and returns the updated document.
func (collection *collection) FindUpdate(ctx context.Context, patch any, sort Sort, filter ...Filter) (any, error) {
	// Check for collection
	if collection.Collection == nil {
		return nil, ErrOutOfOrder
	}

	// Trace
	ctx, matched, _ := trace.WithCollection(ctx, trace.OpFindUpdate, collection.Database().Name(), collection.Name())
	defer trace.Do(ctx, collection.traceFn, time.Now())

	// Execute operation
	result := collection.Collection.FindOneAndUpdate(ctx, and(filter...), bson.D{{"$set", patch}}, &options.FindOneAndUpdateOptions{
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
