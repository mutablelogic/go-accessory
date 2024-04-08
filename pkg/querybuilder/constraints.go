/*
Key(...) is a factory method that returns a primary key constraint:

	Key() => PRIMARY KEY
	Key("a","b") => PRIMARY KEY (a,b)

Key(...).Foreign(name, ...) is a factory method that returns a foreign key constraint:

	Key().Foreign("other_table") => REFERENCES other_table
	Key("b","c").Foreign("other_table","c1","c2") => FOREIGN KEY (b, c) REFERENCES other_table (c1, c2)
	Key("a").Foreign("other_table") => FOREIGN KEY a REFERENCES other_table
	Key().Foreign("other_table","c1") => REFERENCES other_table (c1)
	Key().Foreign("other_table").OnDeleteRestrict() => REFERENCES other_table (c1) ON DELETE RESTRICT
	Key().Foreign("other_table").OnDeleteCascade() => REFERENCES other_table (c1) ON DELETE CASCADE
	Key().Foreign("other_table").OnDeleteNoAction() => REFERENCES other_table (c1) ON DELETE NO ACTION

Key(...).Unique() is a factory method that returns a unique key constraint:

	Key().Unique() => UNIQUE
	Key("a","b").Unique() => UNIQUE (a,b)
*/
package querybuilder

import (
	// Package imports
	quote "github.com/mutablelogic/go-accessory/pkg/querybuilder/quote"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type keyConstraint struct {
	name
	flags
	column  []any
	foreign struct {
		name   string
		column []any
	}
}

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

// Create a key which defaults to a primary key
func Key(v ...string) keyConstraint {
	q := keyConstraint{}
	q.flags = primarykey
	for _, v := range v {
		q.column = append(q.column, N(v))
	}
	return q
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (q keyConstraint) Name(v string) keyConstraint {
	q.name = N(v)
	return q
}

func (q keyConstraint) Unique() keyConstraint {
	q.flags = uniquekey
	return q
}

func (q keyConstraint) Foreign(name string, column ...string) keyConstraint {
	q.flags = foreignkey
	q.foreign.name = name
	for _, v := range column {
		q.foreign.column = append(q.foreign.column, N(v))
	}
	return q
}

func (q keyConstraint) OnDeleteRestrict() keyConstraint {
	q.flags |= restrict
	return q
}

func (q keyConstraint) OnDeleteCascade() keyConstraint {
	q.flags |= cascade
	return q
}

func (q keyConstraint) OnDeleteNoAction() keyConstraint {
	q.flags |= noAction
	return q
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (q keyConstraint) leftStr() string {
	if len(q.column) > 0 {
		return quote.Join("(" + quote.JoinSep(",", q.column...) + ")")
	} else {
		return ""
	}
}

func (q keyConstraint) foreignStr() string {
	if !q.flags.Is(foreignkey) {
		return ""
	} else if len(q.foreign.column) > 0 {
		return quote.Join("REFERENCES", q.foreign.name, "("+quote.JoinSep(",", q.foreign.column...)+")")
	} else {
		return quote.Join("REFERENCES", q.foreign.name)
	}
}

func (q keyConstraint) onDeleteStr() string {
	if q.flags.Is(foreignkey) && q.flags.Is(restrict|cascade|noAction) {
		return quote.Join("ON DELETE", q.flags&(restrict|cascade|noAction))
	} else {
		return ""
	}
}

func (q keyConstraint) String() string {
	str := ""
	if !q.flags.Is(foreignkey) || len(q.column) > 0 {
		str = quote.Join(q.flags)
	}
	return quote.Join(str, q.leftStr(), q.foreignStr(), q.onDeleteStr())
}
