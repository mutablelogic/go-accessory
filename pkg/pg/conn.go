// PostgreSQL driver for go-accessory
package pg

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"sync"
	"time"

	// Packages

	pgx "github.com/jackc/pgx/v5"
	trace "github.com/mutablelogic/go-accessory/pkg/trace"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
	. "github.com/mutablelogic/go-accessory"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type conn struct {
	sync.Mutex
	*pgx.Conn
	opts

	// Default database name
	def string

	// Database mapping
	databases map[string]*database
}

var _ Conn = (*conn)(nil)

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	defaultTimeout = 10 * time.Second
	defaultSchema  = "public"
)

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

// Connect to MongoDB server
func Open(ctx context.Context, url *url.URL, opts ...Opt) (Conn, error) {
	// Create client
	conn := new(conn)
	conn.databases = make(map[string]*database, 1)

	// Set the defaults
	conn.opts.timeout = defaultTimeout
	conn.opts.schema = defaultSchema

	// Apply the client options BEFORE we connect
	for _, opt := range opts {
		if err := opt(&conn.opts); err != nil {
			return nil, err
		}
	}

	// Add connect_timeout to URL
	urlSet(url, "connect_timeout", fmt.Sprint(uint(conn.Timeout().Seconds())))

	// If there is an application name, set it
	if conn.opts.applicationName != "" {
		urlSet(url, "application_name", conn.opts.applicationName)
	}

	// If there is a user, set it (and optionally, a password)
	if conn.opts.user != "" {
		urlSet(url, "user", conn.opts.user)
		if conn.opts.password != "" {
			urlSet(url, "password", conn.opts.password)
		}
	}

	// Connection Trace
	defer trace.Do(trace.WithUrl(ctx, trace.OpConnect, url), conn.tracefn, time.Now())

	// Connect to the server
	if c, err := pgx.Connect(ctx, url.String()); err != nil {
		return nil, err
	} else {
		conn.Conn = c
	}

	// Get current database
	if err := conn.Conn.QueryRow(ctx, "SELECT current_database()").Scan(&conn.def); err != nil {
		return nil, err
	}
	if conn.def == "" {
		if err := conn.Conn.QueryRow(ctx, "SELECT current_user").Scan(&conn.def); err != nil {
			return nil, err
		}
	}

	// Return success
	return conn, nil
}

// Close the client
func (conn *conn) Close() error {
	conn.Mutex.Lock()
	defer conn.Mutex.Unlock()

	var result error

	// Close the connection
	if conn.Conn != nil {
		result = errors.Join(result, conn.Conn.Close(context.Background()))
	}

	// Return any errors
	return result
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (conn *conn) String() string {
	str := "<pg.conn"
	if conn.Conn != nil {
		str += fmt.Sprintf(" url=%q", conn.Conn.Config().ConnString())
		str += fmt.Sprintf(" db=%q", conn.Name())
	}
	str += fmt.Sprint(" timeout=", conn.Timeout())
	return str + ">"
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Return the timeout for the client
func (conn *conn) Timeout() time.Duration {
	return conn.timeout
}

// Ping the client, return an error if not reachable
func (conn *conn) Ping(ctx context.Context) error {
	conn.Mutex.Lock()
	defer conn.Mutex.Unlock()

	// Ping Trace
	defer trace.Do(trace.WithOp(ctx, trace.OpPing), conn.tracefn, time.Now())
	return conn.Conn.Ping(ctx)
}

// Return a database object for a specific database
func (conn *conn) Database(name string) Database {
	// Obtain cached database or create a new one
	if d, exists := conn.databases[name]; exists {
		return d
	} else if d := conn.new_database(name, conn.opts.schema); d == nil {
		return nil
	} else {
		conn.databases[name] = d
		return d
	}
}

// Return all existing databases on the server
func (conn *conn) Databases(ctx context.Context) ([]Database, error) {
	if conn.IsClosed() {
		return nil, ErrOutOfOrder.With("connection closed")
	}

	var result []Database
	databases, err := conn.databaseNames(ctx)
	if err != nil {
		return nil, err
	}
	for _, name := range databases {
		if db := conn.Database(name); db != nil {
			result = append(result, db)
		}
	}
	return result, nil
}

// Perform operations within a transaction. Rollback or apply
// changes to the database depending on error return.
func (conn *conn) Do(ctx context.Context, fn func(context.Context) error) error {
	// Add a transaction counter to the context
	ctx = trace.WithTx(ctx)

	// Perforn a 'BEGIN' transaction
	if err := conn.BeginTx(ctx); err != nil {
		return err
	}

	// Commit or rollback
	var result error
	if err := fn(ctx); err != nil {
		// Trace
		defer trace.Do(trace.WithOp(ctx, trace.OpRollback), conn.tracefn, time.Now())

		// Rollback
		result = errors.Join(result, err)
		if err := conn.RollbackTx(ctx); err != nil {
			result = errors.Join(result, err)
		}
	} else {
		// Trace
		defer trace.Do(trace.WithOp(ctx, trace.OpCommit), conn.tracefn, time.Now())

		// Commit
		if err := conn.CommitTx(ctx); err != nil {
			result = errors.Join(result, err)
		}
	}

	// Return any errors
	return result
}

// Return the name of the default database
func (conn *conn) Name() string {
	return conn.def
}

// Return a collection object for a specific struct
func (this *conn) Collection(proto any) Collection {
	if db := this.Database(this.def); db != nil {
		return db.Collection(proto)
	} else {
		return nil
	}
}

// Insert documents of the same type to the database within a transaction.
// The document keys are updated if the document is writable.
func (this *conn) Insert(ctx context.Context, v ...any) error {
	if db := this.Database(this.def); db == nil {
		return ErrNotFound.With("database", this.def)
	} else {
		return db.Insert(ctx, v...)
	}
}

///////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func urlSet(u *url.URL, key, value string) *url.URL {
	q := u.Query()
	q.Set(key, value)
	u.RawQuery = q.Encode()
	return u
}

func (c *conn) databaseNames(ctx context.Context) ([]string, error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	rows, err := c.Query(ctx, "SELECT datname FROM pg_catalog.pg_database WHERE NOT datistemplate")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Collect database names
	var result []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		} else {
			result = append(result, name)
		}
	}
	return result, nil
}
