package pool

import (
	"context"
	"fmt"
	"net/url"
	"sync"
	"sync/atomic"

	// Package imports
	multierror "github.com/hashicorp/go-multierror"
	mongodb "github.com/mutablelogic/go-accessory/pkg/mongodb"
	trace "github.com/mutablelogic/go-accessory/pkg/trace"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
	. "github.com/mutablelogic/go-accessory"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type pool struct {
	p     sync.Pool
	max   int64
	size  atomic.Int64
	drain atomic.Bool

	// Connection parameters
	uri     *url.URL
	mongodb []mongodb.ClientOpt

	// Trace function
	trace trace.Func
}

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	schemeMongo1 = "mongodb"
	schemeMongo2 = "mongodb+srv"
	schemeSqlite = "sqlite"
)

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

// Create a pool with the given URL. The URL should be of scheme "mongodb"
// "file", or "sqlite".
func New(ctx context.Context, uri *url.URL, opts ...Option) Pool {
	pool := new(pool)

	// Check parameters
	if uri == nil {
		trace.Err(trace.WithUrl(ctx, trace.OpConnect, uri), pool.trace, ErrBadParameter.With("uri"))
		return nil
	} else {
		pool.uri = uri
	}

	// Client options - before client is created
	for _, opt := range opts {
		if err := opt(pool); err != nil {
			trace.Err(trace.WithUrl(ctx, trace.OpConnect, uri), pool.trace, err)
			return nil
		}
	}

	// Set the connection factory function
	switch uri.Scheme {
	case schemeMongo1, schemeMongo2:
		pool.p.New = func() any {
			// Check for draining
			if pool.drain.Load() {
				return nil
			}
			// Create MongoDB connection
			if conn, err := pool.NewMongoDB(context.Background()); err != nil {
				trace.Err(trace.WithUrl(ctx, trace.OpConnect, uri), pool.trace, err)
				return nil
			} else {
				return conn
			}
		}
		// Add the first connection to the pool
		if conn, err := pool.NewMongoDB(ctx); err != nil {
			trace.Err(trace.WithUrl(ctx, trace.OpConnect, uri), pool.trace, err)
			return nil
		} else {
			pool.size.Add(1)
			pool.Put(conn)
		}
	default:
		trace.Err(trace.WithUrl(ctx, trace.OpConnect, uri), pool.trace, ErrBadParameter.With(uri))
		return nil
	}

	// Client options - after client is created
	for _, opt := range opts {
		if err := opt(pool); err != nil {
			trace.Err(trace.WithUrl(ctx, trace.OpConnect, uri), pool.trace, err)
			return nil
		}
	}

	// Return success
	return pool
}

// Drain the pool and close all connections
func (pool *pool) Close() error {
	var result error

	// Signal we are draining
	pool.drain.Store(true)

	// Drain until no more connections
	for {
		if conn := pool.p.Get(); conn == nil {
			break
		} else if err := conn.(Conn).Close(); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Signal we are no longer draining
	pool.drain.Store(false)

	// Return any errors
	return result
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (pool *pool) String() string {
	str := "<pool"
	if pool.uri != nil {
		str += fmt.Sprintf(" uri=%q", pool.uri.String())
	}
	if size := pool.Size(); size > 0 {
		str += fmt.Sprint(" size=", pool.Size())
	}
	if pool.max > 0 {
		str += fmt.Sprint(" max_size=", pool.max)
	}
	return str + ">"
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Get a connection from the connection pool
func (pool *pool) Get() Conn {
	if pool.max > 0 && pool.size.Load() >= pool.max {
		trace.Err(context.Background(), pool.trace, ErrOutOfOrder.With("maximum number of connections reached"))
		return nil
	}
	conn := pool.p.Get()
	if conn == nil {
		return nil
	}
	pool.size.Add(1)
	return conn.(Conn)
}

// Put a connection back into the connection pool
func (pool *pool) Put(v Conn) {
	if v != nil {
		pool.p.Put(v)
		pool.size.Add(-1)
	}
}

func (pool *pool) Size() int {
	return int(pool.size.Load())
}

///////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// Create a new MongoDB connection with the required options
func (pool *pool) NewMongoDB(ctx context.Context) (Conn, error) {
	return mongodb.Open(ctx, pool.uri, pool.mongodb...)
}
