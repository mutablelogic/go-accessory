/*
N(...).Select() is a factory method that returns a select statement on a single table:

	N("t").Select().Distinct() => "SELECT DISTINCT * FROM t"
	N("t").Select("a,",b") => "SELECT a,b FROM t"
	N("t").Select("a,",b").Sort("a") => "SELECT a,b FROM t ORDER BY a"
	N("t").Select("a,",b").Sort().Limit(100) => "SELECT a,b FROM t LIMIT 100"
	N("t").FromSchema("public").As("x").Select(N("a").FromSchema("x").As("a"),b").SortDesc("1").Limit(100) => "SELECT x.a AS a,b FROM public.t x ORDER BY 1 DESC LIMIT 100"
*/
package querybuilder

import (
	"fmt"
	"reflect"

	// Package imports
	quote "github.com/mutablelogic/go-accessory/pkg/querybuilder/quote"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type selectExpr struct {
	name
	flags
	column []any
}

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (n name) Select(v ...any) selectExpr {
	q := selectExpr{name: n}
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

func (q selectExpr) Distinct() selectExpr {
	q.flags |= distinct
	return q
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (q selectExpr) strExprs() string {
	if len(q.column) == 0 {
		return "*"
	} else {
		return quote.JoinSep(",", q.column...)
	}
}

func (q selectExpr) String() string {
	return quote.Join("SELECT", (q.flags & distinct), q.strExprs(), "FROM", q.name)
}
