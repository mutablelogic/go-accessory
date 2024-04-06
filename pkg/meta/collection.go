package meta

import (
	"fmt"
	"reflect"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

// Collection is the mapping between a reflect.Type and a collection
type Collection struct {
	// Go type (which is always a struct)
	Type reflect.Type

	// Collection name
	Name string

	// Attributes
	Attr []*Attr
}

// Extra is an interface which provides extra metadata for a collection
type extra interface {
	Name() string
}

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

// Create a new Reflect type, or return nil if the type is not valid
func New(v reflect.Value, tag string) *Collection {
	r := new(Collection)

	// Determine the type
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil
	} else {
		r.Type = v.Type()
		r.Name = r.Type.Name()
	}

	// If v is a collection, then set the name and schema
	if extra, ok := v.Interface().(extra); ok {
		if extra.Name() != "" {
			r.Name = extra.Name()
		}
	}

	// Enumerate atrributes
	if attr := attrForStruct(nil, r.Type, tag); attr == nil {
		return nil
	} else {
		r.Attr = attr
	}

	// Return success
	return r
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (c *Collection) String() string {
	str := "<collection"
	str += fmt.Sprintf(" type=%q", c.Type)
	str += fmt.Sprintf(" name=%q", c.Name)
	if len(c.Attr) > 0 {
		str += fmt.Sprintf(" attrs=%v", c.Attr)
	}
	return str + ">"
}
