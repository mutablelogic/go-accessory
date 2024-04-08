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
	notnull
	unique
	primarykey
	distinct
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
	if f.Is(notnull) {
		str += "NOT NULL "
	}
	if f.Is(unique) {
		str += "UNIQUE "
	}
	if f.Is(primarykey) {
		str += "PRIMARY KEY "
	}
	if str == "" {
		return ""
	} else {
		return str[:len(str)-1]
	}
}
