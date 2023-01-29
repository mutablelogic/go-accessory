package mongodb

import (
	"context"
	"io"
	"reflect"

	// Modules
	multierror "github.com/hashicorp/go-multierror"
	driver "go.mongodb.org/mongo-driver/mongo"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type cursor struct {
	ctx context.Context
	c   *driver.Cursor
	t   reflect.Type
	eof bool
}

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewCursor(c *driver.Cursor, t reflect.Type) *cursor {
	this := new(cursor)
	this.c = c
	this.t = t
	return this
}

func (cursor *cursor) Close() error {
	var result error
	if cursor.c == nil {
		return nil
	}
	if err := cursor.c.Close(context.Background()); err != nil {
		result = multierror.Append(result, err)
	}
	cursor.c = nil
	return result
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (cursor *cursor) Next(ctx context.Context) (any, error) {
	if cursor.eof {
		return nil, io.EOF
	}
	if next := cursor.c.Next(ctx); next {
		doc := reflect.New(cursor.t).Interface()
		if err := cursor.c.Decode(doc); err != nil {
			return nil, err
		} else {
			return doc, nil
		}
	}
	cursor.eof = true
	if err := cursor.c.Close(ctx); err != nil {
		return nil, err
	}
	cursor.c = nil
	return nil, io.EOF
}
