package mongodb_test

import (
	"context"
	"io"
	"testing"
	"time"

	mongodb "github.com/djthorpe/go-mongodb"
	assert "github.com/stretchr/testify/assert"
)

const (
	uri = "mongodb://cm1:27017/"
)

func Test_Client_001(t *testing.T) {
	assert := assert.New(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	c, err := mongodb.New(ctx, uri)
	assert.NoError(err)
	assert.NotNil(c)
	t.Log(c)
	assert.NoError(c.Close())
}

func Test_Client_002(t *testing.T) {
	assert := assert.New(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	c, err := mongodb.New(ctx, uri, mongodb.OptDatabase("admin"))
	assert.NoError(err)
	assert.NotNil(c)
	names, err := c.Collections(ctx)
	assert.NoError(err)
	assert.NotNil(names)
	assert.NoError(c.Close())
}

func Test_Client_003(t *testing.T) {
	assert := assert.New(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	c, err := mongodb.New(ctx, uri, mongodb.OptDatabase("test"))
	assert.NoError(err)
	assert.NotNil(c)
	defer c.Close()

	// Insert a document into the database and return the key
	type Doc struct {
		Id string `bson:"_id,omitempty"`
	}
	var doc Doc
	key, err := c.Insert(ctx, &doc)
	assert.NoError(err)
	assert.NotNil(key)
	assert.Equal(key, doc.Id)
	t.Log(doc)
}

func Test_Client_004(t *testing.T) {
	assert := assert.New(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	c, err := mongodb.New(ctx, uri, mongodb.OptDatabase("test"))
	assert.NoError(err)
	assert.NotNil(c)
	defer c.Close()

	// Insert a document into the database and return the key
	type Doc struct {
		Id string `bson:"_id,omitempty"`
	}
	var doc1, doc2 Doc
	keys, err := c.InsertMany(ctx, &doc1, &doc2)
	assert.NoError(err)
	assert.NotNil(keys)
	assert.Len(keys, 2)
	assert.Equal(keys, []string{doc1.Id, doc2.Id})
}

func Test_Client_005(t *testing.T) {
	assert := assert.New(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	c, err := mongodb.New(ctx, uri, mongodb.OptDatabase("test"))
	assert.NoError(err)
	assert.NotNil(c)
	defer c.Close()

	// Insert a document into the database and return the key
	type Doc struct {
		Id string `bson:"_id,omitempty"`
	}
	var doc1, doc2 Doc
	keys, err := c.InsertMany(ctx, &doc1, &doc2)
	assert.NoError(err)
	assert.NotNil(keys)
	assert.Len(keys, 2)
	assert.Equal(keys, []string{doc1.Id, doc2.Id})
}

func Test_Client_006(t *testing.T) {
	assert := assert.New(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	c, err := mongodb.New(ctx, uri, mongodb.OptDatabase("test"))
	assert.NoError(err)
	assert.NotNil(c)
	defer c.Close()

	type Doc struct {
		Id string `bson:"_id,omitempty"`
	}
	var doc Doc

	// Insert a document into the database
	key, err := c.Insert(ctx, &doc)
	assert.NoError(err)

	// Delete a document from the database
	result, err := c.Delete(ctx, Doc{}, mongodb.F().EqualsId(key))
	assert.NoError(err)
	assert.Equal(int64(1), result)
}

func Test_Client_007(t *testing.T) {
	assert := assert.New(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	c, err := mongodb.New(ctx, uri, mongodb.OptDatabase("test"))
	assert.NoError(err)
	assert.NotNil(c)
	defer c.Close()

	type Doc struct {
		Id string `bson:"_id,omitempty"`
	}
	var doc Doc

	// Insert a document into the database
	key, err := c.Insert(ctx, &doc)
	assert.NoError(err)

	// Find the document in the database
	err = c.Find(ctx, &doc, nil, mongodb.F().EqualsId(key))
	assert.NoError(err)
	assert.Equal(doc.Id, key)
}

func Test_Client_008(t *testing.T) {
	assert := assert.New(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	c, err := mongodb.New(ctx, uri, mongodb.OptDatabase("test"))
	assert.NoError(err)
	assert.NotNil(c)
	defer c.Close()

	type Doc struct {
		Id string `bson:"_id,omitempty"`
	}
	var doc Doc

	// Insert a document into the database
	num, err := c.InsertMany(ctx, &doc, &doc)
	assert.NoError(err)
	assert.Len(num, 2)

	// Find the documents in the database
	cursor, err := c.FindMany(ctx, Doc{}, nil, nil)
	assert.NoError(err)
	assert.NotNil(cursor)
	for {
		if err := cursor.Next(ctx, &doc); err == io.EOF {
			break
		} else if err != nil {
			assert.NoError(err)
		}
		t.Log(doc)
	}
}

func Test_Client_009(t *testing.T) {
	assert := assert.New(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	c, err := mongodb.New(ctx, uri, mongodb.OptDatabase("test"))
	assert.NoError(err)
	assert.NotNil(c)
	defer c.Close()

	type Doc struct {
		Id    string `bson:"_id,omitempty"`
		Value string `bson:"value,omitempty"`
	}
	doc := Doc{Value: "test"}

	// Insert a document into the database
	key, err := c.Insert(ctx, &doc)
	assert.NoError(err)
	assert.NotEmpty(key)

	// Retrieve the document from the database
	err = c.Find(ctx, &doc, nil, mongodb.F().EqualsId(key))
	assert.NoError(err)

	// Update the document in the database
	doc.Id = ""
	doc.Value = "updated"
	matched, modified, err := c.Update(ctx, &doc, mongodb.F().EqualsId(key))
	assert.NoError(err)
	assert.Equal(int64(1), matched)
	assert.Equal(int64(1), modified)

}
