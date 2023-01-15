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

type conn struct {
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

var _ Conn = (*conn)(nil)

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
func Open(ctx context.Context, url *url.URL, opts ...ClientOpt) (Conn, error) {
	// Create client
	this := new(conn)
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
	conn, err := driver.Connect(ctx, clientOpts...)
	if err != nil {
		return nil, err
	} else {
		this.Client = conn
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
func (conn *conn) Close() error {
	var result error

	// Return nil if already closed
	if conn.Client == nil {
		return nil
	}

	// Disconnect with default timeout
	ctx, cancel := context.WithTimeout(context.Background(), conn.Timeout())
	defer cancel()

	// Trace
	defer trace.Do(trace.WithUrl(ctx, trace.OpDisconnect, conn.url), conn.tracefn, time.Now())

	// Disconnect
	if err := conn.Disconnect(ctx); err != nil {
		result = multierror.Append(result, err)
	} else {
		conn.Client = nil
	}

	// Release resources
	conn.db = nil
	conn.meta = nil

	// Return any errors
	return result
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (conn *conn) String() string {
	str := "<mongodb.conn"
	str += fmt.Sprint(" timeout=", conn.Timeout())
	if db := conn.Database(defaultDatabase); db != nil {
		str += fmt.Sprint(" db=", db)
	}
	return str + ">"
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Ping the primary server and return any errors
func (conn *conn) Ping(ctx context.Context) error {
	// Return nil if already closed
	if conn.Client == nil {
		return ErrOutOfOrder.With("Ping")
	}

	// Trace
	defer trace.Do(trace.WithUrl(ctx, trace.OpPing, conn.url), conn.tracefn, time.Now())

	// Perform ping
	return conn.Client.Ping(c(ctx), readpref.Primary())
}

// Timeout returns the default timeout for any client operations
func (conn *conn) Timeout() time.Duration {
	if conn.Client == nil {
		return defaultTimeout
	} else if to := conn.Client.Timeout(); to != nil {
		return *to
	} else {
		return defaultTimeout
	}
}

// Database returns a database with a specific name
func (conn *conn) Database(v string) Database {
	if conn.db == nil {
		return nil
	} else if _, exists := conn.db[v]; !exists {
		if v != defaultDatabase {
			conn.db[v] = NewDatabase(conn, v, conn.protosToMeta, conn.tracefn)
		}
	}
	return conn.db[v]
}

// Databases returns all databases that exist on the server
func (conn *conn) Databases(ctx context.Context) ([]Database, error) {
	// Check client is open
	if conn.Client == nil {
		return nil, ErrOutOfOrder.With("Databases")
	}

	// Obtain database names
	names, err := conn.ListDatabaseNames(c(ctx), bson.D{})
	if err != nil {
		return nil, err
	}

	// Create database objects
	result := make([]Database, 0, len(names))
	for _, name := range names {
		result = append(result, conn.Database(name))
	}

	// Return success
	return result, nil
}

// Exists returns true if a database with given name exists
func (conn *conn) Exists(ctx context.Context, v string) bool {
	if conn.Client == nil {
		return false
	}

	// Return true if database exists
	names, err := conn.ListDatabaseNames(c(ctx), bson.M{"name": v})
	if err != nil {
		return false
	}

	// Return true if only one database exists
	return len(names) == 1
}

// Do executes a function within a transaction. If the function returns
// any error, the transaction is rolled back. Otherwise, the transaction
// is applied to the database.
func (conn *conn) Do(ctx context.Context, fn func(context.Context) error) error {
	session, err := conn.Client.StartSession(&options.SessionOptions{})
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
		defer trace.Do(trace.WithOp(ctx, trace.OpRollback), conn.tracefn, time.Now())

		// Rollback
		result = multierror.Append(result, err)
		if err := session.AbortTransaction(c(ctx)); err != nil {
			result = multierror.Append(result, err)
		}
	} else {
		// Trace
		defer trace.Do(trace.WithOp(ctx, trace.OpCommit), conn.tracefn, time.Now())

		// Commit
		if err := session.CommitTransaction(c(ctx)); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Return any errors
	return result
}

// Return a collection in the default database
func (conn *conn) Collection(proto any) Collection {
	return conn.Database(defaultDatabase).Collection(proto)
}

// Return the name of the default database, or empty string if none
func (conn *conn) Name() string {
	return conn.Database(defaultDatabase).Name()
}

// Return an empty filter specification
func (conn *conn) F() Filter {
	return NewFilter()
}

// Return an empty sort specification
func (conn *conn) S() Sort {
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
func (conn *conn) registerProto(proto any, name string) *meta {
	t := derefType(reflect.TypeOf(proto))
	if t.Kind() != reflect.Struct {
		return nil
	}
	if meta, exists := conn.meta[t]; exists && meta.Name == name {
		return meta
	} else if meta := NewMeta(t, name); meta != nil {
		conn.meta[t] = meta
		return meta
	} else {
		return nil
	}
}

// return metadata from prototype
func (conn *conn) protoToMeta(proto any) *meta {
	t := derefType(reflect.TypeOf(proto))
	return conn.meta[t]
}

// Return metadata from more than one prototype which
// are all of the same type, or else return nil
func (conn *conn) protosToMeta(protos ...any) *meta {
	// No protos = no way!
	if len(protos) == 0 {
		return nil
	}
	// Check for nil
	if protos[0] == nil {
		return nil
	}

	// Get name from collection or type
	meta := conn.protoToMeta(protos[0])
	if meta == nil {
		meta = NewMeta(reflect.TypeOf(protos[0]), "")
		if meta == nil {
			return nil
		} else {
			conn.meta[meta.Type] = meta
		}
	}

	// Return emptyCollection if remaining protos are different
	if len(protos) > 1 {
		if otherMeta := conn.protosToMeta(protos[1:]...); otherMeta == nil || otherMeta != meta {
			return nil
		}
	}

	// Return success
	return meta
}
