package mongodb

import (
	"context"

	// Namespace imports
	. "github.com/mutablelogic/go-accessory"
)

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Return the name of the default database, or empty string if none
func (client *client) Name() string {
	if db := client.Database(defaultDatabase); db != nil {
		return db.Name()
	}
	return ""
}

// Return a collection in the default database
func (client *client) Collection(proto any) Collection {
	return client.Database(defaultDatabase).Collection(proto)
}

func (client *client) Insert(ctx context.Context, doc ...any) error {
	return client.Database(defaultDatabase).Insert(ctx, doc)
}
