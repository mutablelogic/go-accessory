package sqlite

///////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo CFLAGS: -Iv3.40.1
#include <sqlite3.h>
*/
import "C"

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Return version
func Version() (string, int, string) {
	return C.GoString(C.sqlite3_libversion()), int(C.sqlite3_libversion_number()), C.GoString(C.sqlite3_sourceid())
}
