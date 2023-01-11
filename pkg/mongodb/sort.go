package mongodb

import (
	// Packages
	bson "go.mongodb.org/mongo-driver/bson"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
	. "github.com/mutablelogic/go-accessory"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type sort struct {
	bson.D
}

var _ Sort = (*sort)(nil)

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewSort() *sort {
	return &sort{bson.D{}}
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Add ascending sort order
func (sort *sort) Asc(fields ...string) error {
	return ErrNotImplemented
}

// Add descending sort order
func (sort *sort) Desc(fields ...string) error {
	return ErrNotImplemented
}

// Limit the number of documents returned
func (sort *sort) Limit(limit int64) error {
	return ErrNotImplemented
}
