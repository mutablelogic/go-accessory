package mongodb

import (
	"context"
	"io"

	// Modules
	driver "go.mongodb.org/mongo-driver/mongo"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type Cursor struct {
	*driver.Cursor
	eof bool
}

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func newCursor(cursor *driver.Cursor) *Cursor {
	this := new(Cursor)
	this.Cursor = cursor
	return this
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (cursor *Cursor) Next(ctx context.Context, doc any) error {
	if cursor.eof {
		return io.EOF
	}
	if next := cursor.Cursor.Next(ctx); next {
		if err := cursor.Decode(doc); err != nil {
			return err
		}
		return nil
	}
	cursor.eof = true
	if err := cursor.Close(ctx); err != nil {
		return err
	}
	return io.EOF
}
