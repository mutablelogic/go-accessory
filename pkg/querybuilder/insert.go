/*
N(...).Insert() is a factory method that returns an insert:

	N("t").Insert() => "INSERT INTO t DEFAULT VALUES"
	N("t").Insert("a", "b") => "INSERT INTO t (a,b) VALUES ($1, $2)"
	N("t").Insert("a", "b").Returning() => "INSERT INTO t (a,b) VALUES ($1, $2) RETURNING *"
	N("t").Insert("a", "b").Returning("a","b") => "INSERT INTO t (a,b) VALUES ($1, $2) RETURNING a,b"
*/
package querybuilder

import (
	"fmt"
	"reflect"
	"strings"

	// Packages
	quote "github.com/mutablelogic/go-accessory/pkg/querybuilder/quote"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type insertTable struct {
	flags
	name
	column    []name
	values    []any
	returning []name
}

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (n name) Insert(v ...any) insertTable {
	q := insertTable{name: n}
	for _, v := range v {
		switch v := v.(type) {
		case string:
			q.column = append(q.column, N(v))
		case name:
			q.column = append(q.column, v)
		default:
			panic("Invalid type: " + fmt.Sprint(reflect.TypeOf(v)))
		}
	}
	return q
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (q insertTable) Returning(v ...any) insertTable {
	q.returning = make([]name, 0, len(v))
	for _, v := range v {
		switch v := v.(type) {
		case string:
			q.returning = append(q.returning, N(v))
		case name:
			q.returning = append(q.returning, v)
		default:
			panic("Invalid type: " + fmt.Sprint(reflect.TypeOf(v)))
		}
	}
	return q
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func names(v []name) string {
	result := make([]string, 0, len(v))
	for _, v := range v {
		result = append(result, v.Name())
	}
	return quote.Identifiers(result...)
}

func (q insertTable) placeholders() string {
	result := make([]string, 0, len(q.column))
	for i := range q.column {
		result = append(result, fmt.Sprintf("$%d", i+1))
	}
	return strings.Join(result, ",")
}

func (q insertTable) returns() string {
	if q.returning == nil {
		return ""
	}
	if len(q.returning) == 0 {
		return "RETURNING *"
	}
	return "RETURNING " + names(q.returning)
}

func (q insertTable) String() string {
	if len(q.column) == 0 {
		return join("INSERT INTO", q.name.String(), "DEFAULT VALUES")
	} else if len(q.values) == 0 {
		return join("INSERT INTO", q.name.String(), "("+names(q.column)+")", "VALUES", "("+q.placeholders()+")", q.returns())
	} else {
		return "TODO"
	}
}
