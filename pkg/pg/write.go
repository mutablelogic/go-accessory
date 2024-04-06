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

type Write interface {
	// WriteInsert will insert a row into the database. The metadata type and
	// the data must be compatible.
	WriteInsert(context.Context, *meta.Collection, ...any) error

	// WriteInsertWithSchema will insert a row into the database. The metadata type and
	// the data must be compatible.
	WriteInsertWithSchema(context.Context, *meta.Collection, string, ...any) error
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (c *conn) WriteInsert(ctx context.Context, meta *meta.Collection, data ...any) error {
	return c.WriteInsertWithSchema(ctx, meta, "", data...)
}

func (c *conn) WriteInsertWithSchema(ctx context.Context, meta *meta.Collection, schema string, data ...any) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	// Create statement
	st := N(meta.Name).WithSchema(schema).Insert()

	// Execute statement
	defer trace.Do(trace.WithExec(ctx, st), c.tracefn, time.Now())
	if _, err := c.Conn.Exec(ctx, st.String()); err != nil {
		return fmt.Errorf("Collection %q: %w", meta.Name, err)
	}

	// Return success
	return nil
}
