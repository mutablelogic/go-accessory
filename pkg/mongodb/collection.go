package mongodb

import (
	"fmt"
	"reflect"
	"strings"

	// Packages
	primitive "go.mongodb.org/mongo-driver/bson/primitive"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

// Collection is the mapping between a reflect.Type and a collection name
type collection struct {
	// Go type (which is always a struct)
	Type reflect.Type

	// Collection name
	Name string

	// Field index which is used as the primary key
	Key []int
}

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	// Tag to use for identifying fields
	structTag = "bson"

	// Field name which is used as the primary key
	structKey = "_id"
)

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewCollection(t reflect.Type, name string) *collection {
	collection := new(collection)
	collection.Name = name
	t = derefType(t)
	if t.Kind() != reflect.Struct {
		return nil
	} else {
		collection.Type = t
	}

	// Find the field which is used as the primary key
	fields := reflect.VisibleFields(t)
	for _, field := range fields {
		name, _ := structTagValue(field)
		if name == "" {
			continue
		} else if name == structKey {
			collection.Key = field.Index
		}
	}

	return collection
}

// Set the key for a document. Return ErrNotModified if the key
// cannot be set in the document.
func (collection *collection) SetKey(doc, key any) (string, error) {
	v := derefValue(reflect.ValueOf(doc))
	if v.Kind() != reflect.Struct || v.Type() != collection.Type {
		return "", ErrBadParameter.Withf("SetKey: invalid document of type %T, expecting %s", doc, collection.Type)
	}
	if collection.Key == nil || !v.CanSet() {
		return "", ErrNotModified.Withf("SetKey: cannot set key in document of type %T", doc)
	}
	// Obtain the field to set
	f := v.FieldByIndex(collection.Key)
	if !f.CanSet() {
		return "", ErrNotModified.Withf("SetKey: cannot set key in document of type %T", doc)
	}
	// If the field type matches the key type, then set directly
	if f.Type() == reflect.TypeOf(key) {
		f.Set(reflect.ValueOf(key))
		return keyToString(key), nil
	}
	// If the field is a string, then convert the key to a string
	if f.Type() == reflect.TypeOf("") {
		keystr := keyToString(key)
		f.SetString(keystr)
		return keystr, nil
	}
	// We don't support other types of keys
	return "", ErrInternalAppError.Withf("SetKey: unsupported key type %T", f.Type())
}

func keyToString(key any) string {
	switch key := key.(type) {
	case string:
		return key
	case primitive.ObjectID:
		return key.Hex()
	default:
		return fmt.Sprint(key)
	}
}

///////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// structTag returns the name for a field and options, or an empty string if
// the field should be ignored.
func structTagValue(f reflect.StructField) (string, map[string]string) {
	// Check for ignored field
	value := strings.TrimSpace(f.Tag.Get(structTag))
	if value == "-" {
		return "", nil
	}

	name := f.Name
	flags := make(map[string]string)
	for i, tag := range strings.Split(value, ",") {
		tag = strings.TrimSpace(tag)
		switch i {
		case 0:
			if tag != "" {
				name = tag
			}
		default:
			kv := strings.SplitN(tag, ":", 2)
			if len(kv) == 2 {
				flags[kv[0]] = strings.TrimSpace(kv[1])
			} else {
				flags[kv[0]] = ""
			}
		}
	}

	// Return success
	return name, flags
}

func derefType(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

func derefValue(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v
}
