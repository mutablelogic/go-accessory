package mongodb

import (
	"context"
	"fmt"
	"reflect"
	"time"

	// Packages
	multierror "github.com/hashicorp/go-multierror"
	bson "go.mongodb.org/mongo-driver/bson"
	driver "go.mongodb.org/mongo-driver/mongo"
	options "go.mongodb.org/mongo-driver/mongo/options"
	readpref "go.mongodb.org/mongo-driver/mongo/readpref"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
	. "github.com/mutablelogic/go-accessory"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type client struct {
	*driver.Client

	// The default timeout
	timeout time.Duration

	// Database mapping. The default database is stored
	// as an empty string
	db map[string]*database

	// Collection metadata mapping.
	col map[reflect.Type]*meta

	// Function to trace calls
	tracefn fnTrace
}

var _ Client = (*client)(nil)

type fnTrace func(context.Context, time.Duration)

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	defaultTimeout  = 10 * time.Second
	defaultDatabase = ""
	emptyCollection = ""
)

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

// Connect to MongoDB server
func Open(ctx context.Context, url string, opts ...ClientOpt) (Client, error) {
	// Create client
	this := new(client)
	this.db = make(map[string]*database, 1)
	this.col = make(map[reflect.Type]*meta, 1)
	this.timeout = defaultTimeout

	// Apply the client options BEFORE we connect
	for _, opt := range opts {
		if err := opt(this); err != nil {
			return nil, err
		}
	}

	// Trace
	now := time.Now()
	defer this.t(ctx, OpConnect, time.Since(now), url)

	// Connect
	clientOpts := []*options.ClientOptions{
		options.Client().ApplyURI(url),
		options.Client().SetConnectTimeout(this.timeout),
		options.Client().SetTimeout(this.timeout),
	}
	client, err := driver.Connect(ctx, clientOpts...)
	if err != nil {
		return nil, err
	} else {
		this.Client = client
	}

	// Apply the client options AFTER we connect
	for _, opt := range opts {
		if err := opt(this); err != nil {
			return nil, err
		}
	}

	// Return success
	return this, nil
}

// Close the client
func (client *client) Close() error {
	var result error

	// Return nil if already closed
	if client.Client == nil {
		return nil
	}

	// Disconnect with default timeout
	ctx, cancel := context.WithTimeout(context.Background(), client.Timeout())
	defer cancel()

	// Trace
	now := time.Now()
	defer client.t(ctx, OpDisconnect, time.Since(now))

	// Disconnect
	if err := client.Disconnect(ctx); err != nil {
		result = multierror.Append(result, err)
	} else {
		client.Client = nil
	}

	// Release resources
	client.db = nil
	client.col = nil

	// Return any errors
	return result
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (client *client) String() string {
	str := "<mongodb"
	str += fmt.Sprint(" timeout=", client.Timeout())
	if db := client.Database(defaultDatabase); db != nil {
		str += fmt.Sprint(" db=", db)
	}
	return str + ">"
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Ping the primary server and return any errors
func (client *client) Ping(ctx context.Context) error {
	// Return nil if already closed
	if client.Client == nil {
		return ErrOutOfOrder.With("Ping")
	}

	// Trace
	now := time.Now()
	defer client.t(ctx, OpPing, time.Since(now))

	// Perform ping
	return client.Client.Ping(c(ctx), readpref.Primary())
}

// Timeout returns the default timeout for any client operations
func (client *client) Timeout() time.Duration {
	if client.Client == nil {
		return defaultTimeout
	} else if to := client.Client.Timeout(); to != nil {
		return *to
	} else {
		return defaultTimeout
	}
}

// Database returns a database with a specific name
func (client *client) Database(v string) Database {
	if client.db == nil {
		return nil
	} else if _, exists := client.db[v]; !exists {
		client.db[v] = NewDatabase(client, v, client.collectionToName, client.updateDocumentWithKey)
	}
	return client.db[v]
}

// Databases returns all databases that exist on the server
func (client *client) Databases(ctx context.Context) ([]Database, error) {
	// Check client is open
	if client.Client == nil {
		return nil, ErrOutOfOrder.With("Databases")
	}

	// Obtain database names
	names, err := client.ListDatabaseNames(c(ctx), bson.D{})
	if err != nil {
		return nil, err
	}

	// Create database objects
	result := make([]Database, 0, len(names))
	for _, name := range names {
		result = append(result, client.Database(name))
	}

	// Return success
	return result, nil
}

// Exists returns true if a database with given name exists
func (client *client) Exists(ctx context.Context, v string) bool {
	if client.Client == nil {
		return false
	}

	// Return true if database exists
	names, err := client.ListDatabaseNames(c(ctx), bson.M{"name": v})
	if err != nil {
		return false
	}

	// Return true if only one database exists
	return len(names) == 1
}

// Do executes a function within a transaction. If the function returns
// any error, the transaction is rolled back. Otherwise, the transaction
// is applied to the database.
func (client *client) Do(ctx context.Context, fn func(context.Context) error) error {
	session, err := client.Client.StartSession(&options.SessionOptions{})
	if err != nil {
		return err
	}
	defer session.EndSession(c(ctx))

	// Trace
	now := time.Now()
	defer client.t(ctx, OpTransaction, time.Since(now))

	// Perform operations within a transaction
	if err := session.StartTransaction(&options.TransactionOptions{}); err != nil {
		return err
	}

	// Commit or rollback
	var result error
	if err := fn(ctx); err != nil {
		result = multierror.Append(result, err)
		if err := session.AbortTransaction(c(ctx)); err != nil {
			result = multierror.Append(result, err)
		}
	} else {
		if err := session.CommitTransaction(c(ctx)); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Return any errors
	return result
}

// Return an empty filter specification
func (client *client) F() Filter {
	return NewFilter()
}

// Return an empty sort specification
func (client *client) S() Sort {
	return NewSort()
}

// Return the name of the default database, or empty string if none
func (client *client) Name() string {
	if db := client.Database(defaultDatabase); db != nil {
		return db.Name()
	}
	return ""
}

///////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// c always returns a context
func c(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	} else {
		return ctx
	}
}

// t is called to trace an operation
func (client *client) t(ctx context.Context, op Operation, dur time.Duration, args ...string) {
	if client.tracefn != nil {
		// TODO: Append operation and arguments
		client.tracefn(c(ctx), dur)
	}
}

func (client *client) protoToCollection(proto any) (*meta, reflect.Type) {
	t := derefType(reflect.TypeOf(proto))
	if t.Kind() != reflect.Struct {
		return nil, t
	}
	if meta, exists := client.col[t]; exists {
		return meta, t
	} else {
		return nil, t
	}
}

func (client *client) collectionToName(protos ...any) string {
	// No protos = no way!
	if len(protos) == 0 {
		return emptyCollection
	}

	// Get name from collection or type
	var name string
	collection, t := client.protoToCollection(protos[0])
	if collection != nil {
		name = collection.Name
	} else {
		name = t.Name()
	}

	// Return emptyCollection if remaining protos are different
	if len(protos) > 1 {
		if otherName := client.collectionToName(protos[1:]...); otherName == emptyCollection || otherName != name {
			return emptyCollection
		}
	}

	// Return success
	return name
}

func (client *client) updateDocumentWithKey(doc, key any) (string, error) {
	// Get the collection from the document
	collection, _ := client.protoToCollection(doc)
	if collection == nil || collection.Key == nil {
		return "", ErrNotModified
	}

	// Set the key
	return collection.SetKey(doc, key)
}
