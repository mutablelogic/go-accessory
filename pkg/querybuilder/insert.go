/*
N(...).Insert() is a factory method that returns an insert:

	N("t").Insert() => "INSERT INTO t DEFAULT VALUES"
	N("t").Insert("a", "b") => "INSERT INTO t (a,b) VALUES (?, ?)"
*/
package querybuilder

import (
	"fmt"
	"reflect"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type insertTable struct {
	flags
	name
	column []any
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
// STRINGIFY

func (q insertTable) String() string {
	if len(q.column) == 0 {
		return join("INSERT INTO", q.name.String(), "DEFAULT VALUES")
	} else {
		return join("INSERT INTO", q.name.String(), "("+joinsep(",", q.column...)+")")
	}
}
