package sqlite

// Import into namespace
//. "github.com/djthorpe/go-errors"

///////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo CFLAGS: -Iv3.40.1
#include <sqlite3.h>
*/
import "C"

///////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	if err := SQError(C.sqlite3_initialize()); err != SQLITE_OK {
		panic(err)
	}
}
