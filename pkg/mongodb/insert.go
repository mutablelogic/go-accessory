package mongodb

import (
	"context"
	"errors"
	"reflect"
	"time"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
	"github.com/hashicorp/go-multierror"
	"github.com/mutablelogic/go-accessory/pkg/trace"
	"go.mongodb.org/mongo-driver/mongo/options"
)

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Insert one or more documents into the default database
func (client *client) Insert(ctx context.Context, doc ...any) error {
	return client.Database(defaultDatabase).Insert(ctx, doc...)
}

// Insert one or more documents into the database
func (database *database) Insert(ctx context.Context, doc ...any) error {
	if database.Database == nil {
		return ErrOutOfOrder
	} else if len(doc) == 0 {
		return ErrBadParameter
	} else if c := database.collectionForProtos(doc...); c == nil {
		t := derefType(reflect.TypeOf(doc[0]))
		return ErrBadParameter.Withf("unknown collection for document of type %q", t.Name())
	} else {
		return c.Insert(ctx, doc...)
	}
}

// Insert one or more documents into the collection
func (collection *collection) Insert(ctx context.Context, doc ...any) error {
	// Check for collection
	if collection.Collection == nil {
		return ErrOutOfOrder
	}

	// Trace
	ctx, _, modified := trace.WithCollection(ctx, trace.OpInsert, collection.Database().Name(), collection.Name())
	defer trace.Do(ctx, collection.traceFn, time.Now())

	// Call one or many
	var err error
	switch len(doc) {
	case 0:
		return ErrBadParameter
	case 1:
		*modified, err = collection.insertOne(ctx, doc[0])
		return err
	default:
		*modified, err = collection.insertMany(ctx, doc...)
		return err
	}
}

///////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (collection *collection) insertOne(ctx context.Context, doc any) (int64, error) {
	r, err := collection.Collection.InsertOne(ctx, doc, &options.InsertOneOptions{})
	if err != nil {
		return -1, err
	}
	if _, err := collection.meta.SetKey(doc, r.InsertedID); !errors.Is(err, ErrNotModified) {
		return 1, err
	} else {
		return 1, nil
	}
}

func (collection *collection) insertMany(ctx context.Context, doc ...any) (int64, error) {
	var result error
	r, err := collection.Collection.InsertMany(ctx, doc, &options.InsertManyOptions{})
	if err != nil {
		return int64(len(r.InsertedIDs)), err
	}
	for i, key := range r.InsertedIDs {
		if _, err := collection.meta.SetKey(doc[i], key); err != nil && !errors.Is(err, ErrNotModified) {
			result = multierror.Append(result, err)
		}
	}
	// Return any errors
	return int64(len(r.InsertedIDs)), result
}
