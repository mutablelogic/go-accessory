package pool

import (
	"time"

	// Package imports
	mongodb "github.com/mutablelogic/go-accessory/pkg/mongodb"
	trace "github.com/mutablelogic/go-accessory/pkg/trace"

	// Namespace Imports
	. "github.com/djthorpe/go-errors"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type ClientOpt func(*pool) error

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Set the maximum size of the pool
func OptMaxSize(v int64) ClientOpt {
	return func(pool *pool) error {
		if v < 0 {
			return ErrBadParameter.With("OptMaxSize")
		} else {
			pool.max = v
		}
		return nil
	}
}

// Set the database name
func OptDatabase(v string) ClientOpt {
	return func(pool *pool) error {
		if pool.uri.Scheme == schemeMongoDB {
			pool.mongodb = append(pool.mongodb, mongodb.OptDatabase(v))
		}
		return nil
	}
}

// Set the default timeout
func OptTimeout(v time.Duration) ClientOpt {
	return func(pool *pool) error {
		if pool.uri.Scheme == schemeMongoDB {
			pool.mongodb = append(pool.mongodb, mongodb.OptTimeout(v))
		}
		return nil
	}
}

// Set the collection metadata
func OptCollection(collection any, name string) ClientOpt {
	return func(pool *pool) error {
		if pool.uri.Scheme == schemeMongoDB {
			pool.mongodb = append(pool.mongodb, mongodb.OptCollection(collection, name))
		}
		return nil
	}
}

// Set the trace function
func OptTrace(fn trace.Func) ClientOpt {
	return func(pool *pool) error {
		pool.trace = fn
		if pool.uri.Scheme == schemeMongoDB {
			pool.mongodb = append(pool.mongodb, mongodb.OptTrace(fn))
		}
		return nil
	}
}
