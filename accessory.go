package accessory

import (
	"context"
	"io"
	"time"
)

///////////////////////////////////////////////////////////////////////////////
// INTERFACES

// Client represents a connection to a database server. Open a connection to
// the client with
//
//	mongodb.Open(context.Context, string, ...ClientOpt) (Client, error)
//
// which returns a client object. The client options can be used to set
// the default database:
//
//	clientopt := mongodb.WithDatabase(string)
//
// You can set the operation timeout using the following option:
//
//	clientopt := mongodb.WithTimeout(time.Duration)
//
// You can map a go struct intstance to a collection name:
//
//	clientopt := mongodb.WithCollection(any, string)
//
// and you can set up a trace function to record operation timings:
//
//	clientopt := mongodb.WithTrace(func(context.Context, time.Duration))
type Client interface {
	io.Closer

	// You can call all database operations on the client instance, which will
	// use the default database or return an error if no default database
	// is set
	Database

	// Return the default timeout for the client
	Timeout() time.Duration

	// Ping the client, return an error if not reachable
	Ping(context.Context) error

	// Return a database object for a specific database
	Database(string) Database

	// Return all existing databases on the server
	Databases(context.Context) ([]Database, error)

	// Perform operations within a transaction. Rollback or apply
	// changes to the database depending on error return.
	Do(context.Context, func(context.Context) error) error

	// Return a filter specification
	F() Filter

	// Return a sort specification
	S() Sort
}

// Database represents a specific database on the server on which operations
// can be performed.
type Database interface {
	// Return the name of the database
	Name() string

	// Return a collection object for a specific struct
	Collection(any) Collection

	// Insert documents of the same type to the database. The document key is updated
	// if the document is writable.
	Insert(context.Context, ...any) error
}

type Collection interface {
	// Return the name of the collection
	Name() string

	// Delete zero or one documents and returns the number of deleted documents (which should be
	// zero or one. The filter argument is used to determine a document to delete. If there is more than
	// one filter, they are ANDed together
	Delete(context.Context, ...Filter) (int64, error)

	// DeleteMany deletes zero or more documents and returns the number of deleted documents.
	DeleteMany(context.Context, ...Filter) (int64, error)

	// Find selects a single document based on filter and sort parameters.
	// It returns ErrNotFound if no document is found
	Find(context.Context, Sort, ...Filter) (any, error)

	// FindMany returns an iterable cursor based on filter and sort parameters.
	// It returns ErrNotFound if no document is found
	FindMany(context.Context, Sort, ...Filter) (Cursor, error)

	// Update zero or one document with given values and return the number
	// of documents matched and modified, neither of which should be more than one.
	Update(context.Context, any, ...Filter) (int64, int64, error)

	// Update zero or more document with given values and return the number
	// of documents matched and modified, neither of which should be more than one.
	UpdateMany(context.Context, any, ...Filter) (int64, int64, error)
}

// Cursor represents an iterable cursor to a result set
type Cursor interface {
	io.Closer

	// Find next document in the result set and return the document. Will
	// return (nil, io.EOF) when no more documents are available.
	Next(context.Context) (any, error)
}

// Filter represents a filter expression for a query
type Filter interface {
	// Match a document primary key. For MongoDB, this can be an ObjectID represented in HEX, or
	// other string.
	Key(string) error
}

// Sort represents a sort specification for a query
type Sort interface {
	// Add ascending sort order fields
	Asc(...string) error

	// Add descending sort order fields
	Desc(...string) error

	// Limit the number of documents returned
	Limit(int64) error
}
