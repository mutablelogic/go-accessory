package querybuilder

//////////////////////////////////////////////////////////////////////////////
// TYPES

type flags uint64

//////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	temporary flags = 1 << iota
	unlogged
	ifNotExists
	ifExists
	notnull
	distinct
	uniquekey
	primarykey
	foreignkey
	cascade
	restrict
	noAction
)

//////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (f flags) Is(v flags) bool {
	return f&v != 0
}

//////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (f flags) String() string {
	str := ""
	if f.Is(distinct) {
		str += "DISTINCT "
	}
	if f.Is(temporary) {
		str += "TEMPORARY "
	}
	if f.Is(unlogged) {
		str += "UNLOGGED "
	}
	if f.Is(ifNotExists) {
		str += "IF NOT EXISTS "
	}
	if f.Is(ifExists) {
		str += "IF EXISTS "
	}
	if f.Is(notnull) {
		str += "NOT NULL "
	}
	if f.Is(uniquekey) {
		str += "UNIQUE "
	}
	if f.Is(primarykey) {
		str += "PRIMARY KEY "
	}
	if f.Is(foreignkey) {
		str += "FOREIGN KEY "
	}
	if f.Is(cascade) {
		str += "CASCADE "
	}
	if f.Is(restrict) {
		str += "RESTRICT "
	}
	if f.Is(noAction) {
		str += "NO ACTION "
	}
	if str == "" {
		return ""
	} else {
		return str[:len(str)-1]
	}
}
