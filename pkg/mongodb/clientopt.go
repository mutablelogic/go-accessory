package mongodb

import (
	"reflect"
	"time"

	// Namespace Imports
	. "github.com/djthorpe/go-errors"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type ClientOpt func(*client) error

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Set the current database
func OptDatabase(v string) ClientOpt {
	return func(client *client) error {
		// Apply after client is connected
		if client.Client != nil {
			client.db[defaultDatabase] = client.Database(v).(*database)
		}
		return nil
	}
}

// Set the default timeout
func OptTimeout(v time.Duration) ClientOpt {
	return func(client *client) error {
		if v == 0 {
			v = defaultTimeout
		}
		if client.Client == nil {
			if v <= 0 {
				return ErrBadParameter.With("timeout")
			}
			client.timeout = v
		}
		return nil
	}
}

// Set the default timeout
func OptCollection(collection any, name string) ClientOpt {
	return func(client *client) error {
		if client.Client == nil {
			return nil
		}

		// Create a new collection
		if meta := NewMeta(reflect.TypeOf(collection), name); meta == nil {
			return ErrBadParameter.Withf("Invalid collectionof type %T", collection)
		} else {
			client.col[meta.Type] = meta
		}
		return nil
	}
}

// Set the trace function
func OptTrace(fn traceFunc) ClientOpt {
	return func(client *client) error {
		if client.Client == nil {
			client.tracefn = fn
		}
		return nil
	}
}
