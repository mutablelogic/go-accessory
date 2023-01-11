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

///////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// and together a set of filters
func and(f ...Filter) any {
	var elems bson.A
	for _, f := range f {
		if f != nil {
			elems = append(elems, f.(*filter).M)
		}
	}
	if len(elems) == 0 {
		return bson.D{}
	}
	if len(elems) == 1 {
		return elems[0]
	}
	return bson.M{"$and": elems}
}
