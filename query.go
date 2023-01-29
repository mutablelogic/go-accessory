package accessory

///////////////////////////////////////////////////////////////////////////////
// TYPES

type QueryFlag uint

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	NOT_NULL QueryFlag = 1 << iota
	PRIMARY_KEY
	UNIQUE_KEY
	ASC
	DESC
	ON_CONFLICT_ROLLBACK
	ON_CONFLICT_ABORT
	ON_CONFLICT_FAIL
	ON_CONFLICT_IGNORE
	ON_CONFLICT_REPLACE
	AUTO_INCREMENT
	TEMPORARY
	IF_NOT_EXISTS
	STRICT
	WITHOUT_ROWID
	NONE           QueryFlag = 0
	QUERY_MIN                = NOT_NULL
	QUERY_MAX                = WITHOUT_ROWID
	QUERY_CONFLICT           = ON_CONFLICT_ROLLBACK | ON_CONFLICT_ABORT | ON_CONFLICT_FAIL | ON_CONFLICT_IGNORE | ON_CONFLICT_REPLACE
	QUERY_SORT               = ASC | DESC
)

///////////////////////////////////////////////////////////////////////////////
// INTERFACES

// Use Q(query) to create an SQL query
// Use N(name, flags...) to create a new table name or column name
// Use S(expr...) to create a new SELECT statement with selected expressions
//    as columns
// Use E(any) to create any kind of expression

// Query represents any kind of SQL Query
type Query interface {
	// Return the SQL query that can be executed
	Query() string

	// Set flags for any query
	With(QueryFlag) Query
}

// Name represents an SQL name (table name, column name)
type Name interface {
	Query

	// Use a specific alias name
	As(string) Name

	// Use a specific schema name
	WithSchema(string) Name

	// Set a declared type
	WithType(string) Name

	// Transform into a CreateTable query with columns. Use TEMPORARY, IF_NOT_EXISTS, STRICT
	// and WITHOUT_ROWID flags to modify the table creation.
	CreateTable(...Name) CreateTable
}

type CreateTable interface {
	Query

	// Add a UNIQUE or PRIMARY_KEY clause to the table
	//WithKey(QueryFlag, ...string) CreateTable
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (v QueryFlag) Is(f QueryFlag) bool {
	return v&f != 0
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (v QueryFlag) String() string {
	if v == NONE {
		return v.FlagString()
	}
	str := ""
	for f := QUERY_MIN; f <= QUERY_MAX; f <<= 1 {
		if v&f != 0 {
			str += "|" + f.FlagString()
		}
	}
	return str[1:]
}

func (v QueryFlag) FlagString() string {
	switch v {
	case NONE:
		return ""
	case NOT_NULL:
		return "NOT NULL"
	case PRIMARY_KEY:
		return "PRIMARY KEY"
	case UNIQUE_KEY:
		return "UNIQUE"
	case ASC:
		return "ASC"
	case DESC:
		return "DESC"
	case ON_CONFLICT_ROLLBACK:
		return "ON CONFLICT ROLLBACK"
	case ON_CONFLICT_ABORT:
		return "ON CONFLICT ABORT"
	case ON_CONFLICT_FAIL:
		return "ON CONFLICT FAIL"
	case ON_CONFLICT_IGNORE:
		return "ON CONFLICT IGNORE"
	case ON_CONFLICT_REPLACE:
		return "ON CONFLICT REPLACE"
	case AUTO_INCREMENT:
		return "AUTOINCREMENT"
	default:
		return "[?? Invalid QueryFlag value]"
	}
}
