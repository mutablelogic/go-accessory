package mongodb

import (
	"context"
	"fmt"
	"time"

	// Packages
	multierror "github.com/hashicorp/go-multierror"
	bson "go.mongodb.org/mongo-driver/bson"
	primitive "go.mongodb.org/mongo-driver/bson/primitive"
	driver "go.mongodb.org/mongo-driver/mongo"
	options "go.mongodb.org/mongo-driver/mongo/options"
	readpref "go.mongodb.org/mongo-driver/mongo/readpref"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type Client struct {
	*driver.Client

	// Database mapping. The default database is stored
	// as an empty string
	db map[string]*Database
}

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	defaultTimeout  = 10 * time.Second
	defaultDatabase = ""
)

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

// Connect to MongoDB server
func New(ctx context.Context, url string, opts ...ClientOpt) (*Client, error) {
	client, err := driver.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		return nil, err
	}

	// Ping the primary
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	// Create client
	this := new(Client)
	this.Client = client
	this.db = make(map[string]*Database, 1)

	// Apply the client options
	for _, opt := range opts {
		if err := opt(this); err != nil {
			return nil, err
		}
	}

	// Return success
	return this, nil
}

func (client *Client) Close() error {
	var result error
	ctx, cancel := context.WithTimeout(context.Background(), client.Timeout())
	defer cancel()
	if err := client.Disconnect(ctx); err != nil {
		result = multierror.Append(result, err)
	}

	// Return any errors
	return result
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (client *Client) String() string {
	str := "<mongodb"
	str += fmt.Sprint(" timeout=", client.Timeout())
	if db := client.Database(); db != nil {
		str += fmt.Sprint(" db=", db)
	}
	return str + ">"
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Timeout returns the timeout for the client
func (client *Client) Timeout() time.Duration {
	if to := client.Client.Timeout(); to != nil {
		return *to
	} else {
		return defaultTimeout
	}
}

// Return a database, or default database if no argument is provided
func (client *Client) Database(name ...string) *Database {
	if len(name) == 0 {
		return client.db[defaultDatabase]
	}
	if len(name) == 1 {
		if db, exists := client.db[name[0]]; exists {
			return db
		}
	}
	return nil
}

// Databases returns the list of databases on the server
func (client *Client) Databases(ctx context.Context) ([]*Database, error) {
	names, err := client.ListDatabaseNames(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	def := client.Database()
	result := make([]*Database, 0, len(names))
	for _, name := range names {
		if def != nil && name == def.Name() {
			result = append(result, def)
		} else {
			result = append(result, client.newDatabase(name))
		}
	}
	// Return success
	return result, nil
}

// Exists returns true if a specific database with given name exists
func (client *Client) Exists(ctx context.Context, name string) bool {
	names, err := client.ListDatabaseNames(ctx, bson.M{"name": name})
	if err != nil {
		return false
	}
	return len(names) == 1
}

// Collections returns the list of collections in the default database
func (client *Client) Collections(ctx context.Context) ([]string, error) {
	db := client.Database()
	if db == nil {
		return nil, ErrBadParameter.With("database not selected")
	}
	if !client.Exists(ctx, db.Name()) {
		return nil, ErrNotFound.With(db.Name())
	}
	return db.Collections(ctx)
}

// Insert a single document to the database and return key for the document
// If writable, the document InsertID field is updated
func (client *Client) Insert(ctx context.Context, v any) (string, error) {
	db := client.Database()
	if db == nil {
		return "", ErrBadParameter.With("database not selected")
	} else {
		return db.Insert(ctx, v)
	}
}

// InsertMany inserts one or more documents of the same type to the database and
// return keys for the document
func (client *Client) InsertMany(ctx context.Context, v ...any) ([]string, error) {
	db := client.Database()
	if db == nil {
		return nil, ErrBadParameter.With("database not selected")
	} else {
		return db.InsertMany(ctx, v...)
	}
}

// Delete deletes a single document from the collection and returns the number of
// documents deleted, which should be zero or one. At least one filter is needed
// in order to select the document to delete
func (client *Client) Delete(ctx context.Context, collection any, filter ...*Filter) (int64, error) {
	db := client.Database()
	if db == nil {
		return -1, ErrBadParameter.With("database not selected")
	} else {
		return db.Delete(ctx, collection, filter...)
	}
}

// DeleteMany deletes a zero or more documents from the collection and returns the number of
// documents deleted. At least one filter is needed in order to select the documents to delete
func (client *Client) DeleteMany(ctx context.Context, collection any, filter ...*Filter) (int64, error) {
	db := client.Database()
	if db == nil {
		return -1, ErrBadParameter.With("database not selected")
	} else {
		return db.DeleteMany(ctx, collection, filter...)
	}
}

// Find locates a document in the collection after filtering and sorting. Returns ErrNotFound
// if no document was found
func (client *Client) Find(ctx context.Context, doc any, sort *Sort, filter ...*Filter) error {
	db := client.Database()
	if db == nil {
		return ErrBadParameter.With("database not selected")
	} else {
		return db.Find(ctx, doc, sort, filter...)
	}
}

///////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (client *Client) newDatabase(name string) *Database {
	return &Database{Database: client.Client.Database(name)}
}

func id(key any) string {
	switch key := key.(type) {
	case string:
		return key
	case primitive.ObjectID:
		return key.Hex()
	default:
		panic("unreachable")
	}
}
