package filter

import (
	"fmt"

	// Packages
	"go.mongodb.org/mongo-driver/bson"

	// Namespace imports
	. "github.com/mutablelogic/go-accessory"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type expr struct {
	v   bson.M
	err error
}

type op string

var _ Filter = (*expr)(nil)

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	OpEq  op = "$eq"
	OpNe  op = "$ne"
	OpGt  op = "$gt"
	OpGte op = "$gte"
	OpLt  op = "$lt"
	OpLte op = "$lte"
	OpIn  op = "$in"
	OpNin op = "$nin"
	OpAnd op = "$and"
	OpNot op = "$not"
	OpOr op = "$or"
	OpNor op = "$nor"
)

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

// Create a new filter with AND expressions
func New(op op,expr ...*expr) *expr {
	expr := new(expr)
	if len(expr) == 0 {
		return expr
	}
	switch op {
		case OpAnd:
			

	filter.v = bson.M{op: make([]bson.M, len(and))}
	for i, v := range and {
		if v.err != nil {
			filter.err = v.err
	return filter
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (expr *expr) String() string {
	str := "<expr"
	if filter.err != nil {
		str += fmt.Sprintf(" error=%q", filter.err.Error())
	}
	if filter.v != nil {
		str += fmt.Sprint(" expr=", filter.v)
	}
	return str + ">"
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Returns any error that has occurred
func (filter *filter) Err() error {
	return filter.err
}

// Return a BSON representation of the filter
func (filter *filter) BSON() bson.M {
	return filter.v
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

/*
func (filter *filter) Key(v string) error {
	if id, err := primitive.ObjectIDFromHex(v); err != nil {
		return err
	} else {
		filter.M["_id"] = bson.M{"$eq": id}
	}
	// Return success
	return nil
}
*/
