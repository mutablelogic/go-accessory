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
// PUBLIC METHODS

// Return number of keywords
func KeywordCount() int {
	return int(C.sqlite3_keyword_count())
}

// Return keyword
func KeywordName(index int) string {
	var cStr *C.char
	var cLen C.int

	if err := SQError(C.sqlite3_keyword_name(C.int(index), &cStr, &cLen)); err != SQLITE_OK {
		return ""
	} else {
		return C.GoStringN(cStr, cLen)
	}
}

// Lookup keyword
func KeywordCheck(v string) bool {
	var cStr *C.char
	var cLen C.int

	// Populate CString
	cStr = C.CString(v)
	cLen = C.int(len(v))
	defer C.free(unsafe.Pointer(cStr))

	// Return check
	return intToBool(int(C.sqlite3_keyword_check(cStr, cLen)))
}
