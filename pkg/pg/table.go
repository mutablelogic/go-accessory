package pg

import (
	"context"
	"fmt"
	"time"

	// Package imports
	meta "github.com/mutablelogic/go-accessory/pkg/meta"
	trace "github.com/mutablelogic/go-accessory/pkg/trace"

	// Namespace imports
	. "github.com/mutablelogic/go-accessory/pkg/querybuilder"
)

///////////////////////////////////////////////////////////////////////////////
// INTERFACES

type Table interface {
	// CreateTable creates a table in the current database. If the second
	// argument is true, then the table is only created if it doesn't already
	// exist. The third argument contains the
	CreateTable(context.Context, *meta.Collection, bool) error

	// CreateTableWithSchema creates a table in the current database with the
	// named schema. If the third argument is true, then the table is only
	// created if it doesn't already exist.
	CreateTableWithSchema(context.Context, *meta.Collection, string, bool) error

	// CreateTempTable creates a temporary table in the current database.
	CreateTempTable(context.Context, *meta.Collection) error
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (c *conn) CreateTable(ctx context.Context, meta *meta.Collection, ifnotexists bool) error {
	return c.CreateTableWithSchema(ctx, meta, "", ifnotexists)
}

func (c *conn) CreateTableWithSchema(ctx context.Context, meta *meta.Collection, schema string, ifnotexists bool) error {
	return c.createTableWithSchemaEx(ctx, meta, schema, ifnotexists, false)
}

func (c *conn) CreateTempTable(ctx context.Context, meta *meta.Collection) error {
	return c.createTableWithSchemaEx(ctx, meta, "", false, true)
}

///////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (c *conn) createTableWithSchemaEx(ctx context.Context, meta *meta.Collection, schema string, ifnotexists bool, temp bool) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	cols, err := PGColumns(meta)
	if err != nil {
		return fmt.Errorf("Collection %q: %w", meta.Name, err)
	}

	// Create statement
	st := N(meta.Name).WithSchema(schema).CreateTable(cols...)
	if ifnotexists {
		st = st.IfNotExists()
	}
	if temp {
		st = st.Temporary()
	}

	// Execute statement
	defer trace.Do(trace.WithExec(ctx, st), c.tracefn, time.Now())
	if _, err := c.Conn.Exec(ctx, st.String()); err != nil {
		return fmt.Errorf("Collection %q: %w", meta.Name, err)
	}

	// Return success
	return nil
}
