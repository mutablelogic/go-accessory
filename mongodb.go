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

	// Return the names of existing collections
	Collections(context.Context) ([]string, error)

	// Insert a single document to the database and return key for the document
	// represented as a string. The document key is updated if the document is
	// writable
	Insert(context.Context, any) (string, error)

	/*
		// Insert many documents of the same type to the database and return keys for the documents
		InsertMany(context.Context, ...any) ([]string, error)

		// Delete zero or one documents and returns the number of deleted documents (which should be
		// zero or one. The document argument is used to determine the collection name, and the
		// filter argument is used to determine a document to delete. If there is more than
		// one filter, they are ANDed together
		Delete(context.Context, any, ...Filter) (int64, error)

		// DeleteMany deletes zero or more documents from collection and returns the number of deleted documents.
		DeleteMany(context.Context, any, ...Filter) (int64, error)

		// Find selects a single document based on filter and sort parameters.
		// It returns ErrNotFound if no document is found
		Find(context.Context, any, Sort, ...Filter) error

		// FindMany returns an iterable cursor based on filter and sort parameters.
		// It returns ErrNotFound if no document is found
		FindMany(context.Context, any, Sort, ...Filter) (Cursor, error)

		// Update a single document with given values in a collection, and returns number
		// of documents matched and modified, neither of which should be more than one.
		Update(context.Context, any, ...Filter) (int64, int64, error)
	*/
}

// Cursor represents an iterable cursor to a result set
type Cursor interface {
	io.Closer

	// Find next document in the result set and unmarshal into doc. Will
	// return io.EOF when no more documents are available.
	Next(context.Context, any) error
}

// Filter represents a filter expression for a query
type Filter interface {
	// Match a document _id (hex)
	ObjectId(string) error
}

// Sort represents a sort specification for a query
type Sort interface {
	// Add ascending sort order
	Asc(...string) error

	// Add descending sort order
	Desc(...string) error

	// Limit the number of documents returned
	Limit(int64) error
}
