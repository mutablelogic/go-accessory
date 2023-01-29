package pool

import (
	// Namespace imports
	. "github.com/djthorpe/go-errors"
	. "github.com/mutablelogic/go-accessory"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type poolconn struct {
	Conn
}

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

// Close implementation returns an error when Close is called directly
func (*poolconn) Close() error {
	return ErrOutOfOrder.With("Cannot call close on a pooled connection")
}
