package mongodb

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	bson "go.mongodb.org/mongo-driver/bson"
	driver "go.mongodb.org/mongo-driver/mongo"
	options "go.mongodb.org/mongo-driver/mongo/options"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
)

// Packages

///////////////////////////////////////////////////////////////////////////////
// TYPES

type Database struct {
	*driver.Database
}

///////////////////////////////////////////////////////////////////////////////
// STRIGIFY

func (database *Database) String() string {
	str := "<database"
	str += fmt.Sprintf(" name=%q", database.Name())
	return str + ">"
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Collections returns the list of collections in the database
func (database *Database) Collections(ctx context.Context) ([]string, error) {
	return database.ListCollectionNames(ctx, bson.D{})
}

// Insert a single document to the database and return key for the document
// If writable, the document InsertID field is updated
func (database *Database) Insert(ctx context.Context, v any) (string, error) {
	// Obtain the collection name
	name, err := collectionNameFromStruct(reflect.TypeOf(v))
	if err != nil {
		return "", err
	}

	// Perform the insert
	result, err := database.Collection(name).InsertOne(ctx, v, &options.InsertOneOptions{})
	if err != nil {
		return "", err
	} else if err := updateDocId(reflect.ValueOf(v), id(result.InsertedID)); err != nil {
		return "", err
	}

	// Return success
	return id(result.InsertedID), nil
}

// Insert one or more documents to the database and return number of documents matched and modified
// all the document types must be the same
func (database *Database) InsertMany(ctx context.Context, v ...any) ([]string, error) {
	if len(v) == 0 {
		return nil, ErrBadParameter.With("no documents to insert")
	}

	// Check all documents are of the same type
	typeOf := reflect.TypeOf(v[0])
	for _, v := range v[1:] {
		if reflect.TypeOf(v) != typeOf {
			return nil, ErrBadParameter.With("documents must be of the same type")
		}
	}

	// Obtain the collection name
	name, err := collectionNameFromStruct(typeOf)
	if err != nil {
		return nil, err
	}

	// Perform the insert
	result, err := database.Collection(name).InsertMany(ctx, v, &options.InsertManyOptions{})
	if err != nil {
		return nil, err
	}

	// Update the document IDs
	keys := make([]string, len(result.InsertedIDs))
	for i, insertedId := range result.InsertedIDs {
		if err := updateDocId(reflect.ValueOf(v[i]), id(insertedId)); err != nil {
			return nil, err
		} else {
			keys[i] = id(insertedId)
		}
	}

	// Return success
	return keys, nil
}

// Delete a single document from collection and returns the number of deleted documents.
// If more than one filter expression is provided, they are ANDed together
func (database *Database) Delete(ctx context.Context, collection any, filter ...*Filter) (int64, error) {
	// Obtain the collection name
	name, err := collectionNameFromStruct(reflect.TypeOf(collection))
	if err != nil {
		return -1, err
	}

	// Return error it no filter is provided
	if len(filter) == 0 {
		return -1, ErrBadParameter.With("no filter argument provided")
	}

	// Perform the delete
	result, err := database.Collection(name).DeleteOne(ctx, and(filter...), &options.DeleteOptions{})
	if err != nil {
		return -1, err
	} else {
		return result.DeletedCount, nil
	}
}

// DeleteMany deletes zero or more documents from collection and returns the number of deleted documents.
// If more than one filter expression is provided, they are ANDed together
func (database *Database) DeleteMany(ctx context.Context, collection any, filter ...*Filter) (int64, error) {
	// Obtain the collection name
	name, err := collectionNameFromStruct(reflect.TypeOf(collection))
	if err != nil {
		return -1, err
	}

	// Return error it no filter is provided
	if len(filter) == 0 {
		return -1, ErrBadParameter.With("no filter argument provided")
	}

	// Perform the delete
	result, err := database.Collection(name).DeleteMany(ctx, and(filter...), &options.DeleteOptions{})
	if err != nil {
		return -1, err
	} else {
		return result.DeletedCount, nil
	}
}

// Find selects a single document in a collection based on filter and sort parameters.
// If more than one filter expression is provided, they are ANDed together
func (database *Database) Find(ctx context.Context, doc any, sort *Sort, filter ...*Filter) error {
	// Obtain the collection name
	name, err := collectionNameFromStruct(reflect.TypeOf(doc))
	if err != nil {
		return err
	}

	// Do the find
	result := database.Collection(name).FindOne(ctx, and(filter...), &options.FindOneOptions{
		Sort: sort.doc(),
	})

	// Check for errors
	if err := result.Err(); err != nil {
		if errors.Is(result.Err(), driver.ErrNoDocuments) {
			return ErrNotFound.Withf("document not found")
		} else {
			return err
		}
	}

	// Decode the document
	if doc != nil {
		return result.Decode(doc)
	} else {
		return nil
	}
}

// FindMany returns in a collection based on filter and sort parameters, and returns
// an iteratable cursor to the result set. If more than one filter expression is provided,
// they are ANDed together
func (database *Database) FindMany(ctx context.Context, collection any, sort *Sort, filter ...*Filter) (*Cursor, error) {
	// Obtain the collection name
	name, err := collectionNameFromStruct(reflect.TypeOf(collection))
	if err != nil {
		return nil, err
	}

	// Do the find
	cursor, err := database.Collection(name).Find(ctx, and(filter...), &options.FindOptions{
		Sort: sort.doc(),
	})
	if err != nil {
		return nil, err
	} else {
		return newCursor(cursor), nil
	}
}

// Update a single document in a collection, and returns number of documents matched and modified
func (database *Database) Update(ctx context.Context, update any, filter ...*Filter) (int64, int64, error) {
	// Obtain the collection name
	name, err := collectionNameFromStruct(reflect.TypeOf(update))
	if err != nil {
		return -1, -1, err
	}

	if result, err := database.Collection(name).UpdateOne(ctx, and(filter...), bson.D{{"$set", update}}, &options.UpdateOptions{}); err != nil {
		return -1, -1, err
	} else {
		return result.MatchedCount, result.ModifiedCount, nil
	}
}

///////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func collectionNameFromStruct(t reflect.Type) (string, error) {
	// Dereference to a struct
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return "", ErrBadParameter.Withf("expected struct, got %T", t)
	}

	// TODO

	// Return the name of the struct
	return t.Name(), nil
}

func updateDocId(v reflect.Value, id string) error {
	// Dereference to a struct
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return ErrBadParameter.Withf("updateDocId: expected struct, got %v", v.Type())
	}
	fields := reflect.VisibleFields(v.Type())
	for _, f := range fields {
		tag := f.Tag.Get("bson")
		if tag == "_id" || strings.HasPrefix(tag, "_id,") {
			v := v.FieldByIndex(f.Index)
			if !v.CanSet() {
				continue
			}
			if v.Kind() != reflect.String {
				return ErrBadParameter.Withf("updateDocId: expected string, got %v", v.Type())
			}
			v.SetString(id)
		}
	}
	return nil
}
