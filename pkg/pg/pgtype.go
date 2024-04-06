package pg

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	// Packages
	meta "github.com/mutablelogic/go-accessory/pkg/meta"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
	. "github.com/mutablelogic/go-accessory/pkg/querybuilder"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type SQType struct {
	Name string
	Size string
	Ptr  bool
}

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	primaryKeyName = "_id"
)

var (
	intType     = reflect.TypeOf(int(0))
	uintType    = reflect.TypeOf(uint(0))
	stringType  = reflect.TypeOf("")
	boolType    = reflect.TypeOf(false)
	float32Type = reflect.TypeOf(float32(0))
	float64Type = reflect.TypeOf(float64(0))
	timeType    = reflect.TypeOf(time.Time{})
)

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Returns an array of columns for a given meta object. Returns an error
// if any column is not supported
func PGColumns(meta *meta.Collection) ([]any, error) {
	var result error

	// Return error if meta is nil
	if meta == nil {
		return nil, ErrBadParameter
	}

	cols := make([]any, 0, len(meta.Fields))
	for _, field := range meta.Fields {
		if col, err := PGColumn(field); err != nil {
			result = errors.Join(result, err)
		} else {
			cols = append(cols, col)
		}
	}

	// Return result and any errors
	return cols, result
}

// Returns a column for a given meta attribute
func PGColumn(field *meta.Field) (any, error) {
	var decltype string

	// Set decltype from type tag if it exists
	if t := field.Get("type"); t != "" {
		decltype = t
	}

	ty := PGType(field.Type)
	if decltype == "" {
		if ty == nil {
			return nil, ErrNotImplemented.Withf("Field %q of type %v not supported", field.Name, field.Type)
		} else {
			decltype = ty.Name + ty.Size
		}
	}
	if field.Name == primaryKeyName {
		if field.Is("omitempty") {
			return nil, ErrBadParameter.Withf("Primary key field %q cannot be omitempty", field.Name)
		}
		if field.Type != stringType {
			return nil, ErrBadParameter.Withf("Primary key field %q must be a string", field.Name)
		}
		decltype = "UUID"
	}

	col := N(field.Name).T(decltype)
	if !field.Is("omitempty") && (ty != nil && !ty.Ptr) {
		col = col.NotNull()
	}
	if field.Is("unique") {
		col = col.Unique()
	}
	if field.Is("primary") || field.Name == primaryKeyName {
		col = col.PrimaryKey()
	}
	if field.Name == primaryKeyName {
		col = col.Default("gen_random_uuid()")
	}

	// Set default
	if v := field.Get("default"); v != "" {
		col = col.Default(v)
	}

	// Return success
	return col, nil
}

// Converts a reflect type into a pg type, with optional tags
// which modify the type
func PGType(t reflect.Type) *SQType {
	var ptr bool
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		ptr = true
	}

	// Primitive types
	switch t {
	case intType:
		return &SQType{Name: "INTEGER", Ptr: ptr}
	case uintType:
		return &SQType{Name: "INTEGER", Ptr: ptr}
	case stringType:
		return &SQType{Name: "TEXT", Ptr: ptr}
	case boolType:
		return &SQType{Name: "BOOLEAN", Ptr: ptr}
	case float32Type:
		return &SQType{Name: "FLOAT4", Ptr: ptr}
	case float64Type:
		return &SQType{Name: "FLOAT8", Ptr: ptr}
	case timeType:
		return &SQType{Name: "TIMESTAMP", Ptr: ptr}
	}

	// Check for array type
	if t.Kind() == reflect.Array {
		if ty := PGType(t.Elem()); ty != nil {
			return &SQType{Name: ty.Name, Size: fmt.Sprintf("[%d]%s", t.Len(), ty.Size), Ptr: ptr}
		} else {
			return nil
		}
	}

	// Check for slice type
	if t.Kind() == reflect.Slice {
		if ty := PGType(t.Elem()); ty != nil {
			return &SQType{Name: ty.Name, Size: "[]" + ty.Size, Ptr: ptr}
		} else {
			return nil
		}
	}

	// Type not currently supported
	return nil
}
