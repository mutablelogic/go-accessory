package mongodb

import (
	// Packages
	bson "go.mongodb.org/mongo-driver/bson"
	primitive "go.mongodb.org/mongo-driver/bson/primitive"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type Filter struct {
	bson.M
}

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func F() *Filter {
	return &Filter{bson.M{}}
}

// EqualsId filter for a document with a specific id. It will
// return the same filter or nil if it was an invalid value
func (filter *Filter) EqualsId(id string) *Filter {
	if id, err := primitive.ObjectIDFromHex(id); err == nil {
		filter.M["_id"] = bson.M{"$eq": id}
		return filter
	} else {
		return nil
	}
}

///////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func or(filter ...*Filter) any {
	var elems bson.A
	for _, f := range filter {
		if f != nil {
			elems = append(elems, f.M)
		}
	}
	if len(elems) == 0 {
		return bson.D{}
	}
	if len(elems) == 1 {
		return elems[0]
	}
	return bson.M{"$or": elems}
}

func and(filter ...*Filter) any {
	var elems bson.A
	for _, f := range filter {
		if f != nil {
			elems = append(elems, f.M)
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
