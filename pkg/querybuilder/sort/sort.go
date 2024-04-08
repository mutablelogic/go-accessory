package sort

import (
	quote "github.com/mutablelogic/go-accessory/pkg/querybuilder/quote"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type sort struct {
	columns       []any
	limit, offset uint64
}

type sortcol struct {
	name string
	desc bool
}

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func Sort(columns ...string) sort {
	var s sort
	return s.Asc(columns...)
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (s sortcol) String() string {
	if s.desc {
		return quote.Identifier(s.name) + " DESC"
	} else {
		return quote.Identifier(s.name)
	}
}

func (s sort) String() string {
	str := ""
	if len(s.columns) > 0 {
		str = quote.Join("ORDER BY", quote.JoinSep(",", s.columns...))
	}
	if s.limit > 0 {
		str = quote.Join(str, "LIMIT", s.limit)
	}
	if s.offset > 0 {
		str = quote.Join(str, "OFFSET", s.offset)
	}
	return str
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (s sort) Asc(columns ...string) sort {
	for _, column := range columns {
		if column != "" {
			s.columns = append(s.columns, sortcol{name: column})
		}
	}
	return s
}

func (s sort) Desc(columns ...string) sort {
	for _, column := range columns {
		if column != "" {
			s.columns = append(s.columns, sortcol{name: column, desc: true})
		}
	}
	return s
}

func (s sort) Limit(v uint64) sort {
	s.limit = v
	return s
}

func (s sort) Offset(v uint64) sort {
	s.offset = v
	return s
}
