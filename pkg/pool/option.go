package pool

import (
	"net/url"
	"time"

	// Package imports
	mongodb "github.com/mutablelogic/go-accessory/pkg/mongodb"
	trace "github.com/mutablelogic/go-accessory/pkg/trace"

	// Namespace Imports
	. "github.com/djthorpe/go-errors"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type Option func(*pool) error

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Set the maximum size of the pool
func OptMaxSize(v int64) Option {
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
func OptDatabase(v string) Option {
	return func(pool *pool) error {
		if pool.uri != nil && (pool.uri.Scheme == schemeMongo1 || pool.uri.Scheme == schemeMongo2) {
			pool.mongodb = append(pool.mongodb, mongodb.OptDatabase(v))
		}
		return nil
	}
}

// Set the default timeout
func OptTimeout(v time.Duration) Option {
	return func(pool *pool) error {
		if pool.uri != nil && (pool.uri.Scheme == schemeMongo1 || pool.uri.Scheme == schemeMongo2) {
			pool.mongodb = append(pool.mongodb, mongodb.OptTimeout(v))
		}
		return nil
	}
}

// Set the collection metadata
func OptCollection(collection any, name string) Option {
	return func(pool *pool) error {
		if pool.uri != nil && (pool.uri.Scheme == schemeMongo1 || pool.uri.Scheme == schemeMongo2) {
			pool.mongodb = append(pool.mongodb, mongodb.OptCollection(collection, name))
		}
		return nil
	}
}

// Set the trace function
func OptTrace(fn trace.Func) Option {
	return func(pool *pool) error {
		pool.trace = fn
		if pool.uri != nil && (pool.uri.Scheme == schemeMongo1 || pool.uri.Scheme == schemeMongo2) {
			pool.mongodb = append(pool.mongodb, mongodb.OptTrace(fn))
		}
		return nil
	}
}

// Set the attach function
func OptAttach(url *url.URL, schema string) Option {
	return func(pool *pool) error {
		return ErrNotImplemented.With("OptAttach")
	}
}
