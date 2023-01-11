package mongodb

import (
	"context"
	"errors"
	"time"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
	"github.com/hashicorp/go-multierror"
	"github.com/mutablelogic/go-accessory/pkg/trace"
	"go.mongodb.org/mongo-driver/mongo/options"
)

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Insert one or more documents into the database
func (database *database) Insert(ctx context.Context, doc ...any) error {
	if database.Database == nil {
		return ErrOutOfOrder
	} else if len(doc) == 0 {
		return ErrBadParameter
	} else {
		return database.Collection(doc[0]).(*collection).Insert(ctx, doc...)
	}
}

// Insert one or more documents into the collection
func (collection *collection) Insert(ctx context.Context, doc ...any) error {
	// Check for collection
	if collection.Collection == nil {
		return ErrOutOfOrder
	}

	// Ensure ctx is not nil
	ctx = c(ctx)

	// Trace
	trace.Do(trace.WithCollection(ctx, trace.OpInsert, collection.Database().Name(), collection.Name()), collection.traceFn, time.Now())

	// Call one or many
	switch len(doc) {
	case 0:
		return ErrBadParameter
	case 1:
		return collection.insertOne(ctx, doc[0])
	default:
		return collection.insertMany(ctx, doc...)
	}
}

///////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (collection *collection) insertOne(ctx context.Context, doc any) error {
	r, err := collection.Collection.InsertOne(ctx, doc, &options.InsertOneOptions{})
	if err != nil {
		return err
	}
	// Update document key
	if _, err := collection.meta.SetKey(doc, r.InsertedID); !errors.Is(err, ErrNotModified) {
		return err
	} else {
		return nil
	}
}

func (collection *collection) insertMany(ctx context.Context, doc ...any) error {
	var result error
	r, err := collection.Collection.InsertMany(ctx, doc, &options.InsertManyOptions{})
	if err != nil {
		return err
	}
	for i, key := range r.InsertedIDs {
		if _, err := collection.meta.SetKey(doc[i], key); err != nil && !errors.Is(err, ErrNotModified) {
			result = multierror.Append(result, err)
		}
	}
	// Return any errors
	return result
}
