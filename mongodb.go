package mongodb

import (
	"context"
	"io"
	"time"
)

///////////////////////////////////////////////////////////////////////////////
// INTERFACES

// Client represents a connection to a MongoDB server. Open a connection to
// the client with
//
//	mongodb.Open(context.Context, string, ...ClientOpt) (Client, error)
//
// which returns a client object. The client options can be used to set
// the default database, timeout, and the mapping between structs and collection
// names:
//
//	clientopt := mongodb.WithDatabase(string)
//	clientopt := mongodb.WithTimeout(time.Duration)
//	clientopt := mongodb.WithCollection(any, string)
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

	// Return true if a database with given name exists
	Exists(context.Context, string) bool

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

	// Return all existing collections in the database
	Collections(context.Context) ([]Collection, error)

	// Insert a single document to the database and return key for the document
	// represented as a string. The document key is updated if the document is
	// writable.
	Insert(context.Context, any) (string, error)

	// Insert many documents of the same type to the database and return keys for the documents.
	// Documents are updated with their key if they are writable.
	InsertMany(context.Context, ...any) ([]string, error)
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
	// Match a document _id (hex)
	ObjectId(string) error
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
