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
	limit *int64
}

var _ Sort = (*sort)(nil)

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewSort() *sort {
	return &sort{bson.D{}, nil}
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Add ascending sort order
func (sort *sort) Asc(fields ...string) error {
	for _, field := range fields {
		sort.D = append(sort.D, bson.E{field, 1})
	}
	return nil
}

// Add descending sort order
func (sort *sort) Desc(fields ...string) error {
	for _, field := range fields {
		sort.D = append(sort.D, bson.E{field, -1})
	}
	return nil
}

// Limit the number of documents returned
func (sort *sort) Limit(limit int64) error {
	if limit < 0 {
		return ErrBadParameter.With("limit")
	}
	sort.limit = &limit
	return nil
}

///////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func sortdoc(s Sort) any {
	if s == nil {
		return bson.D{}
	}
	return s.(*sort).D
}

func sortlimit(s Sort) *int64 {
	if s == nil {
		return nil
	} else {
		return s.(*sort).limit
	}
}
