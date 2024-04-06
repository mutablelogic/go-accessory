package pg

import (
	"context"

	// Namespace imports
	. "github.com/mutablelogic/go-accessory/pkg/querybuilder"
)

///////////////////////////////////////////////////////////////////////////////
// INTERFACES

type Schema interface {
	// CreateSchema creates a schema in the current database. If the second
	// argument is true, then the schema is only created if it doesn't already
	// exist.
	CreateSchema(string, bool) error
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (c *conn) CreateSchema(ctx context.Context, name string, ifnotexists bool) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	st := N(name).CreateSchema()
	if ifnotexists {
		st = st.IfNotExists()
	}
	defer 
	_, err := c.Conn.Exec(ctx, st.String())
	return err
}
