package sqlite

import (
	"fmt"
	"unsafe"

	// Modules
	multierror "github.com/hashicorp/go-multierror"

	// Import into namespace
	. "github.com/djthorpe/go-errors"
)

///////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo CFLAGS: -Iv3.40.1
#include <sqlite3.h>
#include <stdlib.h>
*/
import "C"

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

// Open URL
func OpenUrl(url string, flags OpenFlags, vfs string) (*Conn, error) {
	return OpenPath(url, flags|SQLITE_OPEN_URI, vfs)
}

// Open Path
func OpenPath(path string, flags OpenFlags, vfs string) (*Conn, error) {
	var cVfs, cName *C.char
	var c *C.sqlite3

	// TODO: Look into logging later
	//initFn.Do(func() {
	//	C._sqlite3_config_logging(1)
	//})

	// Check for thread safety
	if C.sqlite3_threadsafe() == 0 {
		return nil, ErrInternalAppError.With("sqlite library was not compiled for thread-safe operation")
	}

	// Set memory database if empty string
	if path == "" || path == DefaultMemory {
		path = DefaultMemory
		flags |= SQLITE_OPEN_MEMORY
	}

	// Set flags, add read/write flag if create flag is set
	if flags == 0 {
		flags = DefaultFlags
	}
	if flags|SQLITE_OPEN_CREATE > 0 {
		flags |= SQLITE_OPEN_READWRITE
	}

	// Remove custom flags, which are not supported by sqlite3_open_v2
	// but are used by higher level packages to add caching, etc.
	flags &= (SQLITE_OPEN_MAX << 1) - 1

	// Populate CStrings
	if vfs != "" {
		cVfs = C.CString(vfs)
		defer C.free(unsafe.Pointer(cVfs))
	}
	cName = C.CString(path)
	defer C.free(unsafe.Pointer(cName))

	// Call sqlite3_open_v2
	if err := SQError(C.sqlite3_open_v2(cName, &c, C.int(flags), cVfs)); err != SQLITE_OK {
		if c != nil {
			C.sqlite3_close_v2(c)
		}
		return nil, err.With(C.GoString(C.sqlite3_errmsg((*C.sqlite3)(c))))
	}

	// Set extended error codes
	if err := SQError(C.sqlite3_extended_result_codes(c, 1)); err != SQLITE_OK {
		C.sqlite3_close_v2(c)
		return nil, err.With(C.GoString(C.sqlite3_errmsg((*C.sqlite3)(c))))
	}

	return (*Conn)(c), nil
}

// Close Connection
func (c *Conn) Close() error {
	var result error

	// Close any active statements
	/*var s *Statement
	for {
		s = c.NextStatement(s)
		if s == nil {
			break
		}
		fmt.Println("finalizing", uintptr(unsafe.Pointer(s)))
		if err := s.Finalize(); err != nil {
			result = multierror.Append(result, err)
		}
	}*/

	// Close database connection
	if err := SQError(C.sqlite3_close_v2((*C.sqlite3)(c))); err != SQLITE_OK {
		result = multierror.Append(result, err)
	}

	// Return any errors
	return result
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (c *Conn) String() string {
	str := "<conn"
	if filename := c.Filename(""); filename != "" {
		str += fmt.Sprintf(" filename=%q", filename)
	}
	if readonly := c.Readonly(""); readonly {
		str += " readonly"
	}
	if autocommit := c.Autocommit(); autocommit {
		str += " autocommit"
	}
	if rowid := c.LastInsertId(); rowid != 0 {
		str += fmt.Sprint(" last_insert_id=", rowid)
	}
	if changes := c.Changes(); changes != 0 {
		str += fmt.Sprint(" rows_affected=", changes)
	}
	return str + ">"
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Get Filename
func (c *Conn) Filename(schema string) string {
	var cSchema *C.char

	// Set schema to default if empty string
	if schema == "" {
		schema = DefaultSchema
	}

	// Populate CStrings
	cSchema = C.CString(schema)
	defer C.free(unsafe.Pointer(cSchema))

	// Call and return
	cFilename := C.sqlite3_db_filename((*C.sqlite3)(c), cSchema)
	if cFilename == nil {
		return ""
	} else {
		return C.GoString(cFilename)
	}
}

// Get Read-only state. Also returns false if database not found
func (c *Conn) Readonly(schema string) bool {
	var cSchema *C.char

	// Set schema to default if empty string
	if schema == "" {
		schema = DefaultSchema
	}

	// Populate CStrings
	cSchema = C.CString(schema)
	defer C.free(unsafe.Pointer(cSchema))

	// Call and return
	r := int(C.sqlite3_db_readonly((*C.sqlite3)(c), cSchema))
	if r == -1 {
		return false
	} else {
		return intToBool(r)
	}
}

// Set extended result codes
func (c *Conn) SetExtendedResultCodes(v bool) error {
	if err := SQError(C.sqlite3_extended_result_codes((*C.sqlite3)(c), C.int(boolToInt(v)))); err != SQLITE_OK {
		return err
	} else {
		return nil
	}
}

// Cache Flush
func (c *Conn) CacheFlush() error {
	if err := SQError(C.sqlite3_db_cacheflush((*C.sqlite3)(c))); err != SQLITE_OK {
		return err
	} else {
		return nil
	}
}

// Release Memory
func (c *Conn) ReleaseMemory() error {
	if err := SQError(C.sqlite3_db_release_memory((*C.sqlite3)(c))); err != SQLITE_OK {
		return err
	} else {
		return nil
	}
}

// Return autocommit state
func (c *Conn) Autocommit() bool {
	return intToBool(int(C.sqlite3_get_autocommit((*C.sqlite3)(c))))
}

// Get last insert id
func (c *Conn) LastInsertId() int64 {
	return int64(C.sqlite3_last_insert_rowid((*C.sqlite3)(c)))
}

// Set last insert id
func (c *Conn) SetLastInsertId(v int64) {
	C.sqlite3_set_last_insert_rowid((*C.sqlite3)(c), C.sqlite3_int64(v))
}

// Get number of changes (rows affected)
func (c *Conn) Changes() int {
	return int(C.sqlite3_changes((*C.sqlite3)(c)))
}

// Interrupt all queries for connection
func (c *Conn) Interrupt() {
	C.sqlite3_interrupt((*C.sqlite3)(c))
}
