package pg

import (
	"context"
	"errors"
	"time"

	// Package imports
	meta "github.com/mutablelogic/go-accessory/pkg/meta"
	trace "github.com/mutablelogic/go-accessory/pkg/trace"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
	. "github.com/mutablelogic/go-accessory/pkg/querybuilder"
)

///////////////////////////////////////////////////////////////////////////////
// INTERFACES

type Write interface {
	// WriteInsert will insert a row into the database. The metadata type and
	// the data must be compatible.
	WriteInsert(context.Context, *meta.Collection, any) error

	// WriteInsertWithSchema will insert a row into the database. The metadata type and
	// the data must be compatible.
	WriteInsertWithSchema(context.Context, *meta.Collection, string, any) error

	// WriteUpdate will update a row in the database, based on the primary key values
	// being set in the data. The metadata type and the data must be compatible.
	WriteUpdate(context.Context, *meta.Collection, any) error

	// WriteUpdateWithSchema will update a row in the database, based on the primary key values
	// being set in the data. The metadata type and the data must be compatible.
	WriteUpdateWithSchema(context.Context, *meta.Collection, string, any) error

	// Delete will delete a row in the database, based on the primary key values
	WriteDelete(context.Context, *meta.Collection, any) error

	// DeleteWithSchema will delete a row in the database, based on the primary key values
	WriteDeleteWithSchema(context.Context, *meta.Collection, string, any) error
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (c *conn) WriteInsert(ctx context.Context, meta *meta.Collection, data any) error {
	return c.WriteInsertWithSchema(ctx, meta, "", data)
}

func (c *conn) WriteInsertWithSchema(ctx context.Context, meta *meta.Collection, schema string, data any) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	// Catch error messages
	var result error

	// Rows to insert - skip when omitempty is set and the data contains a zero-value
	// Any _id field must be zero-valued or else an error occurs
	params := make([]any, 0, len(meta.Fields))
	values := make([]any, 0, len(meta.Fields))
	for _, field := range meta.Fields {
		if v, z, err := meta.Value(field, data); err != nil {
			result = errors.Join(result, err)
		} else if !z { // TODO: add test for omitempty
			params = append(params, N(field.Name))
			values = append(values, v)
		}
	}

	// Return any errors
	if result != nil {
		return result
	}

	// Create statement
	// TODO: Returning should only be for primary keys
	st := N(meta.Name).WithSchema(schema).Insert(params...).Returning()

	// Execute statement
	defer trace.Do(trace.WithExec(ctx, st), c.tracefn, time.Now())
	// TODO: Query and fill primary key with returning data, if data is a pointer

	// Return success
	return nil
}

func (c *conn) WriteUpdate(ctx context.Context, meta *meta.Collection, data any) error {
	return c.WriteUpdateWithSchema(ctx, meta, "", data)
}

func (c *conn) WriteUpdateWithSchema(ctx context.Context, meta *meta.Collection, schema string, data any) error {
	return ErrNotImplemented
}

func (c *conn) WriteDelete(ctx context.Context, meta *meta.Collection, data any) error {
	return c.WriteDeleteWithSchema(ctx, meta, "", data)
}

func (c *conn) WriteDeleteWithSchema(ctx context.Context, meta *meta.Collection, schema string, data any) error {
	return ErrNotImplemented
}
