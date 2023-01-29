package sqlite

import (
	"fmt"
	"sync"
	"unsafe"
)

///////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo CFLAGS: -Iv3.40.1
#include <sqlite3.h>
#include <stdlib.h>

extern void go_config_logger(void* userInfo, int code, char* msg);
static inline int _sqlite3_config_logging(int enable) {
  if(enable) {
    return sqlite3_config(SQLITE_CONFIG_LOG, go_config_logger, NULL);
  } else {
    return sqlite3_config(SQLITE_CONFIG_LOG, NULL, NULL);
  }
}
*/
import "C"

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

var logFn sync.Once

///////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func initLogging() {
	logFn.Do(func() {
		C._sqlite3_config_logging(1)
	})
}

//export go_config_logger
func go_config_logger(userInfo unsafe.Pointer, code C.int, message *C.char) {
	fmt.Println(SQError(code), C.GoString(message))
}
