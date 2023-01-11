package mongodb

import (
	// Packages
	bson "go.mongodb.org/mongo-driver/bson"
	primitive "go.mongodb.org/mongo-driver/bson/primitive"

	// Namespace imports
	. "github.com/mutablelogic/go-accessory"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type filter struct {
	bson.M
}

var _ Filter = (*filter)(nil)

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewFilter() *filter {
	return &filter{bson.M{}}
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (filter *filter) Key(v string) error {
	if id, err := primitive.ObjectIDFromHex(v); err != nil {
		return err
	} else {
		filter.M["_id"] = bson.M{"$eq": id}
	}
	// Return success
	return nil
}
