package meta

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

// Field is a struct field attribute mapped to a collection field
type Field struct {
	Type  reflect.Type
	Name  string
	Index []int
	Tags  []string
}

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewField(index []int, field reflect.StructField, tag string) *Field {
	tags := tagValue(field, tag)
	if len(tags) == 0 {
		return nil
	}
	return &Field{
		Name:  tags[0],
		Type:  field.Type,
		Tags:  tags[1:],
		Index: append(append([]int{}, index...), field.Index...),
	}
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (self *Field) String() string {
	str := "<field"
	str += fmt.Sprintf(" type=%q", self.Type)
	str += fmt.Sprintf(" name=%q", self.Name)
	if self.Index != nil {
		str += fmt.Sprintf(" index=%v", self.Index)
	}
	if len(self.Tags) > 0 {
		str += fmt.Sprintf(" tags=%q", self.Tags)
	}
	return str + ">"
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Is returns true if the attribute has the specified tag
func (attr *Field) Is(tag string) bool {
	for _, t := range attr.Tags {
		if t == tag {
			return true
		} else if strings.HasPrefix(t, tag+":") {
			return true
		}
	}
	return false
}

// Get returns the value of a specific tag
func (attr *Field) Get(tagname string) string {
	for _, tag := range attr.Tags {
		if strings.HasPrefix(tag, tagname+":") {
			return tag[len(tagname)+1:]
		}
	}
	return ""
}

///////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// tagValue returns the value of the field based on tag or field name
// and returns nil if the field should be ignored (not assignable)
func tagValue(field reflect.StructField, tagName string) []string {
	// Check for private field
	if field.Name != "" && unicode.IsLower(rune(field.Name[0])) {
		return nil
	}
	// Check for anonymous field
	if field.Anonymous {
		return nil
	}
	tags := strings.Split(field.Tag.Get(tagName), ",")
	if tags[0] == "-" {
		return nil
	} else if tags[0] == "" {
		tags[0] = field.Name
	}
	return tags
}

// attrForStruct recursively enumerates the fields of a struct and returns
// an array of attributes
func attrForStruct(index []int, t reflect.Type, tag string) []*Field {
	// t should be struct
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil
	}

	// Iterate over fields
	var attrs = []*Field{}
	for _, field := range reflect.VisibleFields(t) {
		attr := NewField(index, field, tag)
		if attr == nil {
			continue
		}
		attrs = append(attrs, attr)
	}

	// Return attributes
	return attrs
}
