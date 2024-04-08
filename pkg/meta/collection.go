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

	// Fields
	Fields []*Field

	// Primary Key Fields
	PrimaryKey []*Field
}

// Extra is an interface which provides extra metadata for a collection
type extra_name interface {
	// Return the collection name
	CollectionName() string
}

// Extra is an interface which provides extra metadata for a collection
type extra_ref interface {
	// Return the struct field reference to another struct
	CollectionRef(string) any
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
	if extra, ok := v.Interface().(extra_name); ok {
		if extra.CollectionName() != "" {
			r.Name = extra.CollectionName()
		}
	}

	// Enumerate atrributes
	if fields := attrForStruct(nil, r.Type, tag); fields == nil {
		return nil
	} else {
		r.Fields = fields
	}

	// Create primary keys
	for _, field := range r.Fields {
		if field.IsPrimaryKey() {
			r.PrimaryKey = append(r.PrimaryKey, field)
		}
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
	if len(c.Fields) > 0 {
		str += fmt.Sprintf(" fields=%v", c.Fields)
	}
	if len(c.PrimaryKey) > 0 {
		str += fmt.Sprintf(" pk=%v", c.PrimaryKey)
	}
	return str + ">"
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Value returns the value of a field in a data object, and whether it is zero-valued
func (c *Collection) Value(field *Field, data any) (any, bool, error) {
	// data object needs to be the same
	v := reflect.ValueOf(data)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct || v.Type() != c.Type {
		return nil, false, fmt.Errorf("invalid data type: %v", v.Kind())
	}

	// Get the field by index
	v = v.FieldByIndex(field.Index)

	// Return the value
	return v.Interface(), v.IsZero(), nil
}
