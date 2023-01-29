package accessory

import (
	"io"
)

///////////////////////////////////////////////////////////////////////////////
// INTERFACES

// Pool represents a connection pool. You can create a connection pool with
// the following code:
//
//	pool := pool.New(ctx, uri, opts...)
//
// where uri is a mongodb:// URI and opts are pool options. You can set the
// maximum size of the pool with the following option:
//
//	opts := pool.WithMaxSize(int)
type Pool interface {
	io.Closer

	// Get a connection from the pool, or return nil
	// if a connection could not be created
	Get() Conn

	// Release a connection back to the pool
	Put(Conn)

	// Return size of pool
	Size() int
}
