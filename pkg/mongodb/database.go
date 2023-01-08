package mongodb

import (
	"context"
	"errors"
	"fmt"

	// Package imports
	bson "go.mongodb.org/mongo-driver/bson"
	driver "go.mongodb.org/mongo-driver/mongo"
	options "go.mongodb.org/mongo-driver/mongo/options"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
	. "github.com/djthorpe/go-mongodb"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type database struct {
	*driver.Database

	// Function to return collection name from prototype
	colFn CollectionNameFunc
	docFn DocumentUpdateFunc
}

// CollectionNameFunc returns the name of the collection for a given set
// of documents, or returns empty string otherwise
type CollectionNameFunc func(...any) string

// DocumentUpdateFunc updates a document with a key (if the document is
// settable and returns the key in string representation. The key
// can either be a string or a primitive.ObjectID
// and will return ErrNotModified if the document is not settable
type DocumentUpdateFunc func(doc, key any) (string, error)

// Ensure *database implements the Database interface
var _ Database = (*database)(nil)

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewDatabase(client *client, name string, colFn CollectionNameFunc, docFn DocumentUpdateFunc) *database {
	return &database{
		Database: client.Client.Database(name),
		colFn:    colFn,
		docFn:    docFn,
	}
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (database *database) String() string {
	str := "<mongodb.database"
	if database.Database != nil {
		str += fmt.Sprintf(" name=%q", database.Name())
	}
	return str + ">"
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Return the name of the database
func (database *database) Name() string {
	if database.Database == nil {
		return ""
	} else {
		return database.Database.Name()
	}
}

// List the collections in the database
func (database *database) Collections(ctx context.Context) ([]string, error) {
	if database.Database == nil {
		return nil, ErrOutOfOrder
	}
	return database.ListCollectionNames(c(ctx), bson.D{})
}

// Insert a single document to the database and return key for the document
func (database *database) Insert(ctx context.Context, document any) (string, error) {
	// Check for database
	if database.Database == nil {
		return "", ErrOutOfOrder
	}
	// Get collection, insert, update the document, then return the id
	if collection, err := database.collection(document); err != nil {
		return "", err
	} else if result, err := collection.InsertOne(c(ctx), document, &options.InsertOneOptions{}); err != nil {
		return "", err
	} else if key, err := database.updateKey(document, result.InsertedID); err != nil {
		return "", err
	} else {
		return key, nil
	}
}

///////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// Return collection from the given document prototypes. All the prototypes
// should be of the same type.
func (database *database) collection(proto ...any) (*driver.Collection, error) {
	if name := database.colFn(proto...); name == "" {
		return nil, ErrBadParameter.With("Unable to determine collection name from prototype")
	} else {
		return database.Database.Collection(name), nil
	}
}

// Update a document with the given key, then return the key as a string. Accepted
// keys are primitive.ObjectID or string. ErrNotModified errors are ignored as
// documents which are not settable are not updated, but that's OK.
func (database *database) updateKey(doc, key any) (string, error) {
	if key, err := database.docFn(doc, key); err == nil {
		return key, nil
	} else if errors.Is(err, ErrNotModified) {
		return key, nil
	} else {
		return "", err
	}
}
