package mongodb

import (
	"fmt"
	"time"

	// Package imports
	trace "github.com/mutablelogic/go-accessory/pkg/trace"

	// Namespace Imports
	. "github.com/djthorpe/go-errors"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type ClientOpt func(*conn) error

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Set the current database
func OptDatabase(v string) ClientOpt {
	return func(conn *conn) error {
		// Apply after client is connected
		if conn.Client != nil {
			fmt.Println("setting database to", v, "=>", conn.Database(v))
			conn.db[defaultDatabase] = conn.Database(v).(*database)
		}
		return nil
	}
}

// Set the default timeout
func OptTimeout(v time.Duration) ClientOpt {
	return func(conn *conn) error {
		if v == 0 {
			v = defaultTimeout
		}
		if conn.Client == nil {
			if v <= 0 {
				return ErrBadParameter.With("timeout")
			}
			conn.timeout = v
		}
		return nil
	}
}

// Set the default timeout
func OptCollection(collection any, name string) ClientOpt {
	return func(conn *conn) error {
		if conn.Client == nil {
			return nil
		}

		// Create a new collection
		if meta := conn.registerProto(collection, name); meta == nil {
			return ErrBadParameter.Withf("Invalid collection of type %T", collection)
		}
		return nil
	}
}

// Set the trace function
func OptTrace(fn trace.Func) ClientOpt {
	return func(conn *conn) error {
		if conn.Client == nil {
			conn.tracefn = fn
		}
		return nil
	}
}
