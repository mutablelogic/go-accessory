package meta

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

// Attr is a struct field attribute mapped to a collection field
type Attr struct {
	Type     reflect.Type
	Name     string
	Index    []int
	Tags     []string
	Children []*Attr
}

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewAttr(index []int, field reflect.StructField, tag string) *Attr {
	tags := tagValue(field, tag)
	if len(tags) == 0 {
		return nil
	}
	return &Attr{
		Name:  tags[0],
		Type:  field.Type,
		Tags:  tags[1:],
		Index: append(append([]int{}, index...), field.Index...),
	}
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (self *Attr) String() string {
	str := "<attr"
	str += fmt.Sprintf(" type=%q", self.Type)
	str += fmt.Sprintf(" name=%q", self.Name)
	if self.Index != nil {
		str += fmt.Sprintf(" index=%v", self.Index)
	}
	if len(self.Tags) > 0 {
		str += fmt.Sprintf(" tags=%q", self.Tags)
	}
	if len(self.Children) > 0 {
		str += fmt.Sprintf(" children=%v", self.Children)
	}
	return str + ">"
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
func attrForStruct(index []int, t reflect.Type, tag string) []*Attr {
	// t should be struct
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil
	}

	// Iterate over fields
	var attrs = []*Attr{}
	for _, field := range reflect.VisibleFields(t) {
		attr := NewAttr(index, field, tag)
		if attr == nil {
			continue
		}
		attrs = append(attrs, attr)
	}

	// Return attributes
	return attrs
}
