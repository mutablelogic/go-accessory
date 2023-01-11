package mongodb

import (
	"context"
	"fmt"
	"net/url"
	"reflect"
	"time"

	// Packages
	multierror "github.com/hashicorp/go-multierror"
	trace "github.com/mutablelogic/go-accessory/pkg/trace"
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

	// The URL used to connect
	url *url.URL

	// The default timeout
	timeout time.Duration

	// Database mapping. The default database is stored
	// as an empty string
	db map[string]*database

	// Collection metadata mapping.
	meta map[reflect.Type]*meta

	// Function to trace calls
	tracefn trace.Func
}

var _ Client = (*client)(nil)

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
func Open(ctx context.Context, url *url.URL, opts ...ClientOpt) (Client, error) {
	// Create client
	this := new(client)
	this.db = make(map[string]*database, 1)
	this.meta = make(map[reflect.Type]*meta, 1)
	this.timeout = defaultTimeout
	this.url = url

	// Apply the client options BEFORE we connect
	for _, opt := range opts {
		if err := opt(this); err != nil {
			return nil, err
		}
	}

	// Ensure context is not nil
	ctx = c(ctx)

	// Trace
	defer trace.Do(trace.WithUrl(ctx, trace.OpConnect, url), this.tracefn, time.Now())

	// Connect
	clientOpts := []*options.ClientOptions{
		options.Client().ApplyURI(url.String()),
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
	defer trace.Do(trace.WithUrl(ctx, trace.OpDisconnect, client.url), client.tracefn, time.Now())

	// Disconnect
	if err := client.Disconnect(ctx); err != nil {
		result = multierror.Append(result, err)
	} else {
		client.Client = nil
	}

	// Release resources
	client.db = nil
	client.meta = nil

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
	defer trace.Do(trace.WithUrl(ctx, trace.OpPing, client.url), client.tracefn, time.Now())

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
		client.db[v] = NewDatabase(client, v, client.protosToMeta, client.tracefn)
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

	// Add a transaction counter to the context
	ctx = trace.WithTx(c(ctx))

	// Perform operations within a transaction
	if err := session.StartTransaction(&options.TransactionOptions{}); err != nil {
		return err
	}

	// Commit or rollback
	var result error
	if err := fn(ctx); err != nil {
		// Trace
		defer trace.Do(trace.WithOp(ctx, trace.OpRollback), client.tracefn, time.Now())

		// Rollback
		result = multierror.Append(result, err)
		if err := session.AbortTransaction(c(ctx)); err != nil {
			result = multierror.Append(result, err)
		}
	} else {
		// Trace
		defer trace.Do(trace.WithOp(ctx, trace.OpCommit), client.tracefn, time.Now())

		// Commit
		if err := session.CommitTransaction(c(ctx)); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Return any errors
	return result
}

// Return a collection in the default database
func (client *client) Collection(proto any) Collection {
	return client.Database(defaultDatabase).Collection(proto)
}

// Return the name of the default database, or empty string if none
func (client *client) Name() string {
	return client.Database(defaultDatabase).Name()
}

// Return an empty filter specification
func (client *client) F() Filter {
	return NewFilter()
}

// Return an empty sort specification
func (client *client) S() Sort {
	return NewSort()
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

// register a mapping from a prototype to a collection name
func (client *client) registerProto(proto any, name string) *meta {
	t := derefType(reflect.TypeOf(proto))
	if t.Kind() != reflect.Struct {
		return nil
	}
	if meta, exists := client.meta[t]; exists && meta.Name == name {
		return meta
	} else if meta := NewMeta(t, name); meta != nil {
		client.meta[t] = meta
		return meta
	} else {
		return nil
	}
}

// return metadata from prototype
func (client *client) protoToMeta(proto any) *meta {
	t := derefType(reflect.TypeOf(proto))
	return client.meta[t]
}

// Return metadata from more than one prototype which
// are all of the same type, or else return nil
func (client *client) protosToMeta(protos ...any) *meta {
	// No protos = no way!
	if len(protos) == 0 {
		return nil
	}
	// Check for nil
	if protos[0] == nil {
		return nil
	}

	// Get name from collection or type
	meta := client.protoToMeta(protos[0])
	if meta == nil {
		meta = NewMeta(reflect.TypeOf(protos[0]), "")
		if meta == nil {
			return nil
		} else {
			client.meta[meta.Type] = meta
		}
	}

	// Return emptyCollection if remaining protos are different
	if len(protos) > 1 {
		if otherMeta := client.protosToMeta(protos[1:]...); otherMeta == nil || otherMeta != meta {
			return nil
		}
	}

	// Return success
	return meta
}
