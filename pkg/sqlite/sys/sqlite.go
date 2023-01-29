package sqlite

import (
	"strings"
)

///////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo CFLAGS: -DSQLITE_THREADSAFE=2
#cgo CFLAGS: -DSQLITE_DEFAULT_WAL_SYNCHRONOUS=1
#cgo CFLAGS: -DSQLITE_ENABLE_UNLOCK_NOTIFY
#cgo CFLAGS: -DSQLITE_ENABLE_FTS3
#cgo CFLAGS: -DSQLITE_ENABLE_FTS5
#cgo CFLAGS: -DSQLITE_ENABLE_RTREE
#cgo CFLAGS: -DSQLITE_ENABLE_DBSTAT_VTAB
#cgo CFLAGS: -DSQLITE_LIKE_DOESNT_MATCH_BLOBS
#cgo CFLAGS: -DSQLITE_OMIT_DEPRECATED
#cgo CFLAGS: -DSQLITE_ENABLE_JSON1
#cgo CFLAGS: -DSQLITE_ENABLE_SESSION
#cgo CFLAGS: -DSQLITE_ENABLE_SNAPSHOT
#cgo CFLAGS: -DSQLITE_ENABLE_PREUPDATE_HOOK
#cgo CFLAGS: -DSQLITE_ENABLE_GEOPOLY
#cgo CFLAGS: -DSQLITE_USE_ALLOCA
#cgo CFLAGS: -DSQLITE_ENABLE_COLUMN_METADATA
#cgo CFLAGS: -DHAVE_USLEEP=1
#cgo CFLAGS: -DSQLITE_DQS=0

// Ref: https://github.com/crawshaw/sqlite/blob/master/sqlite.go
#cgo windows LDFLAGS: -Wl,-Bstatic -lwinpthread -Wl,-Bdynamic
#cgo linux LDFLAGS: -ldl -lm
#cgo linux CFLAGS: -std=c99
#cgo openbsd LDFLAGS: -lm
#cgo openbsd CFLAGS: -std=c99
#cgo freebsd LDFLAGS: -lm
#cgo freebsd CFLAGS: -std=c99

// Include C source files
#cgo CFLAGS: -Iv3.40.1
#include <sqlite3.c>
*/
import "C"

///////////////////////////////////////////////////////////////////////////////
// TYPES

type OpenFlags C.int
type Conn C.sqlite3

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	DefaultSchema = "main"
	DefaultMemory = ":memory:"
	DefaultFlags  = SQLITE_OPEN_CREATE | SQLITE_OPEN_READWRITE
)

const (
	SQLITE_OPEN_NONE         OpenFlags = 0
	SQLITE_OPEN_READONLY     OpenFlags = C.SQLITE_OPEN_READONLY     // The database is opened in read-only mode. If the database does not already exist, an error is returned.
	SQLITE_OPEN_READWRITE    OpenFlags = C.SQLITE_OPEN_READWRITE    // The database is opened for reading and writing if possible, or reading only if the file is write protected by the operating system. In either case the database must already exist, otherwise an error is returned.
	SQLITE_OPEN_CREATE       OpenFlags = C.SQLITE_OPEN_CREATE       // The database is created if it does not already exist
	SQLITE_OPEN_URI          OpenFlags = C.SQLITE_OPEN_URI          // The filename can be interpreted as a URI if this flag is set.
	SQLITE_OPEN_MEMORY       OpenFlags = C.SQLITE_OPEN_MEMORY       // The database will be opened as an in-memory database. The database is named by the "filename" argument for the purposes of cache-sharing, if shared cache mode is enabled, but the "filename" is otherwise ignored.
	SQLITE_OPEN_NOMUTEX      OpenFlags = C.SQLITE_OPEN_NOMUTEX      // The new database connection will use the "multi-thread" threading mode. This means that separate threads are allowed to use SQLite at the same time, as long as each thread is using a different database connection.
	SQLITE_OPEN_FULLMUTEX    OpenFlags = C.SQLITE_OPEN_FULLMUTEX    // The new database connection will use the "serialized" threading mode. This means the multiple threads can safely attempt to use the same database connection at the same time. (Mutexes will block any actual concurrency, but in this mode there is no harm in trying.)
	SQLITE_OPEN_SHAREDCACHE  OpenFlags = C.SQLITE_OPEN_SHAREDCACHE  // The database is opened shared cache enabled, overriding the default shared cache setting provided by sqlite3_enable_shared_cache().
	SQLITE_OPEN_PRIVATECACHE OpenFlags = C.SQLITE_OPEN_PRIVATECACHE // The database is opened shared cache disabled, overriding the default shared cache setting provided by sqlite3_enable_shared_cache().
	//	SQLITE_OPEN_NOFOLLOW     OpenFlags = C.SQLITE_OPEN_NOFOLLOW                         // The database filename is not allowed to be a symbolic link
	SQLITE_OPEN_MIN = SQLITE_OPEN_READONLY
	SQLITE_OPEN_MAX = SQLITE_OPEN_PRIVATECACHE
)

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (v OpenFlags) StringFlag() string {
	switch v {
	case SQLITE_OPEN_NONE:
		return "SQLITE_OPEN_NONE"
	case SQLITE_OPEN_READONLY:
		return "SQLITE_OPEN_READONLY"
	case SQLITE_OPEN_READWRITE:
		return "SQLITE_OPEN_READWRITE"
	case SQLITE_OPEN_CREATE:
		return "SQLITE_OPEN_CREATE"
	case SQLITE_OPEN_URI:
		return "SQLITE_OPEN_URI"
	case SQLITE_OPEN_MEMORY:
		return "SQLITE_OPEN_MEMORY"
	case SQLITE_OPEN_NOMUTEX:
		return "SQLITE_OPEN_NOMUTEX"
	case SQLITE_OPEN_FULLMUTEX:
		return "SQLITE_OPEN_FULLMUTEX"
	case SQLITE_OPEN_SHAREDCACHE:
		return "SQLITE_OPEN_SHAREDCACHE"
	case SQLITE_OPEN_PRIVATECACHE:
		return "SQLITE_OPEN_PRIVATECACHE"
	default:
		return "[?? Invalid OpenFlags value]"
	}
}

func (v OpenFlags) String() string {
	if v == SQLITE_OPEN_NONE {
		return v.StringFlag()
	}
	str := ""
	for f := SQLITE_OPEN_MIN; f <= SQLITE_OPEN_MAX; f <<= 1 {
		if v&f != 0 {
			str += "|" + f.StringFlag()
		}
	}
	return strings.TrimPrefix(str, "|")
}
