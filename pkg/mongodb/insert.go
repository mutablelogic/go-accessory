package mongodb

import (
	"context"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
)

func (database *database) Insert(ctx context.Context, doc ...any) error {
	// Check for database
	if database.Database == nil {
		return ErrOutOfOrder
	}
	return ErrNotImplemented
}

/*
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

// Insert a single document to the database and return key for the document
func (database *database) InsertMany(ctx context.Context, document ...any) ([]string, error) {
	return nil, ErrNotImplemented
}
*/
